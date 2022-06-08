package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"time"

	u "github.com/mahmednabil109/koorde-overlay/mulitnodetest"
	pd "github.com/mahmednabil109/koorde-overlay/rpc"
)

var (
	NODE_NUM = flag.Int("node-num", 1, "number of the nodes in the network")
)

func main() {
	flag.Parse()

	nodes := make([]u.Nnode, 0)
	rand.Seed(time.Now().UnixNano())

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
			cmd := exec.Command("../../../bin/koorde-overlay",
				"-port", fmt.Sprintf("%d", 8080+i),
				"-first", fmt.Sprint(u.IF(i == 0).Int(1, 0)),
				"-bootstrap", bootstrap_ip,
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

			// simple delay
			time.Sleep(100 * time.Millisecond)
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
	log.Print("Wait until network to settel")

	time.Sleep(10*time.Second + time.Second*time.Duration(len(nodes)))

	// not concurent
	for i := 0; i < *NODE_NUM; i++ {
		pre := time.Now()

		num := string(fmt.Sprintf("%d", i))
		_, err := nodes[i].KC.InitBroadCastRPC(context.Background(), &pd.BlockPacket{Info: num})
		if err != nil {
			panic(err)
		}

		log.Printf(
			"broadcast %d from %s:%d in %v",
			i,
			nodes[i].ID,
			nodes[i].Port,
			time.Since(pre),
		)
	}

	// checking the correctness
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

			set := map[string]struct{}{}

			for _, b := range blocks.Block {
				_, ok := set[b.Info]
				if ok {
					log.Printf("dublicate !! %+v", b)
				}
				set[b.Info] = struct{}{}
			}

			log.Printf(
				"blocks num in %d is %d %d %v",
				nodes[i].Port,
				len(blocks.Block),
				len(set),
				len(set) == *NODE_NUM,
			)
		}(i)
	}
	g_wg.Wait()

}
