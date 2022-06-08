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
	LOOP     = flag.Int("loop", 10, "number of iterations")
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
			cmd := exec.Command("../../../bin/koorde-overlay-network",
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

	log.Print("Start the testing")

	delay_sum := 0.0

	total_pre := time.Now()
	logs := make([]string, *LOOP)
	var lookup_wg sync.WaitGroup

	// not concurent
	for i := 0; i < *LOOP; i++ {
		lookup_wg.Add(1)

		go func(i int) {
			defer lookup_wg.Done()

			j, k := rand.Intn(*NODE_NUM), rand.Intn(*NODE_NUM)

			pre := time.Now()
			reply, err := nodes[j].KC.DLKup(context.Background(), &pd.PeerPacket{SrcId: nodes[k].ID})

			if err != nil {
				panic(err)
			}

			delay_sum += time.Since(pre).Seconds()

			// log.Printf(
			// 	"lookup %d from %s:%d -> %s:%d %v path %d in %v",
			// 	i,
			// 	nodes[j].ID,
			// 	nodes[j].Port,
			// 	nodes[k].ID,
			// 	nodes[k].Port,
			// 	nodes[k].ID == reply.SrcId,
			// 	reply.PathLen,
			// 	time.Since(pre),
			// )
			logs[i] = fmt.Sprintf(
				"lookup %d from %s:%d -> %s:%d %v path %d in %v",
				i,
				nodes[j].ID,
				nodes[j].Port,
				nodes[k].ID,
				nodes[k].Port,
				nodes[k].ID == reply.SrcId,
				reply.PathLen,
				time.Since(pre),
			)
		}(i)

	}
	lookup_wg.Wait()

	for _, k := range logs {
		log.Print(k)
	}

	log.Printf("total time %v", time.Since(total_pre))
	log.Printf("avg RT %v", delay_sum/float64(*LOOP))
	log.Printf("avg TPS %v", float64(*LOOP)/delay_sum)

}
