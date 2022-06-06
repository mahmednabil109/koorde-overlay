package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	u "github.com/mahmednabil109/koorde-overlay/mulitnodetest"
	pd "github.com/mahmednabil109/koorde-overlay/rpc"
)

var (
	NODE_NUM = flag.Int("node-num", 3, "number of nodes in the network")
	DEBUG    = flag.Int("d", 0, "debug flag")
	DSERVER  = flag.String("dserver", "127.0.0.1:16585", "debug server addr")
	LOOKUPS  = flag.Int("lookups", 100, "number of lookups to test with")
)

type IF bool

func (c IF) Int(a, b int) int {
	if c {
		return a
	}
	return b
}

func main() {
	flag.Parse()

	nodes := make([]u.Nnode, 0)
	rand.Seed(time.Now().UnixNano())

	// make concurrent joins :)
	var wg sync.WaitGroup
	var nodes_mux sync.Mutex

	for i := 0; i < *NODE_NUM; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			bootstrap_ip := "127.0.0.1:8080"
			// if i > 1 {
			// 	bootstrap_ip = fmt.Sprintf("127.0.0.1:%d", 8080+rand.Intn(i))
			// }
			cmd := exec.Command("../../bin/koorde-overlay",
				"-port", fmt.Sprintf("%d", 8080+i),
				"-first", fmt.Sprint(u.IF(i == 0).Int(1, 0)),
				"-bootstrap", bootstrap_ip,
				"-debug", fmt.Sprint(*DEBUG),
				"-dserver", *DSERVER,
			)
			f, err := os.Create(fmt.Sprintf("%d.err.log", 8080+i))
			if err != nil {
				panic(err)
			}
			defer f.Close()

			cmd.Stderr = f
			err = cmd.Start()
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Second)

			nodes_mux.Lock()
			defer nodes_mux.Unlock()
			nodes = append(nodes, u.Nnode{
				Cmd:  cmd,
				Port: 8080 + i,
			})
			nodes[len(nodes)-1].Init()
		}(i)
	}
	wg.Wait()

	for i := range nodes {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		peer, err := nodes[i].KC.DGetID(ctx, &pd.Empty{})
		if err != nil {
			panic(err)
		}
		nodes[i].ID = peer.SrcId
		log.Printf("%d has %s", nodes[i].Port, peer.SrcId)
	}
	log.Print("Wait until network to sattel")

	time.Sleep(10*time.Second + time.Second*time.Duration(len(nodes)))
	// time.Sleep(10 * time.Second)

	log.Print("Start Lookup Test")
	var l_wg sync.WaitGroup

	for i := range nodes {
		l_wg.Add(1)

		go func(i int) {
			defer l_wg.Done()

			c_node := nodes[i]
			for j := range nodes {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				reply, err := c_node.KC.DLKup(ctx, &pd.PeerPacket{SrcId: nodes[j].ID})
				if err != nil {
					panic(err)
				}
				log.Printf(
					"lookup form %d to %d -> %v",
					c_node.Port,
					nodes[j].Port,
					nodes[j].ID == reply.SrcId,
				)
			}

		}(i)
	}
	l_wg.Wait()

	log.Print("Getting Pointers")

	// j number of retries to check stability
	for j := 0; j < 2; j++ {

		var p_wg sync.WaitGroup
		for i := range nodes {
			p_wg.Add(1)

			go func(i int) {
				defer p_wg.Done()

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				reply, err := nodes[i].KC.DGetPointers(ctx, &pd.Empty{})
				if err != nil {
					panic(err)
				}

				log.Printf("%d -> %s ,** %s", nodes[i].Port, reply.Succ, reply.D)
			}(i)
		}
		p_wg.Wait()
	}

	log.Print("Start the Broadcasting")

	var b_wg sync.WaitGroup

	for i := range nodes {
		b_wg.Add(1)

		go func(i int) {
			defer b_wg.Done()
			// i := 0
			// for j := range nodes {
			// 	if nodes[j].Port == 8085 {
			// 		i = j
			// 		break
			// 	}
			// }

			start := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			num := string(fmt.Sprint(nodes[i].Port)[2:])
			log.Printf("Init Braodcast from %d %s hi in %v", nodes[i].Port, num, time.Since(start))

			_, err := nodes[i].KC.InitBroadCastRPC(ctx, &pd.BlockPacket{Info: num})
			if err != nil {
				panic(err)
			}
			time.Sleep(1 * time.Second)
		}(i)
	}
	b_wg.Wait()
	time.Sleep(time.Second)

	log.Printf("Done BraodCasting")

	var g_wg sync.WaitGroup
	for i := range nodes {
		g_wg.Add(1)

		go func(i int) {
			defer g_wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			blocks, err := nodes[i].KC.DGetBlocks(ctx, &pd.Empty{})
			if err != nil {
				panic(err)
			}

			log.Printf("blocks in %d is %+v", nodes[i].Port, blocks)
		}(i)
	}
	g_wg.Wait()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	for i := range nodes {
		nodes[i].Cmd.Wait()
	}

}
