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
	"sort"
	"sync"
	"syscall"
	"time"

	u "github.com/mahmednabil109/koorde-overlay/mulitnodetest"
	"github.com/mahmednabil109/koorde-overlay/node"
	pd "github.com/mahmednabil109/koorde-overlay/rpc"
	"github.com/mahmednabil109/koorde-overlay/utils"
)

var (
	NODE_NUM = flag.Int("node-num", 3, "number of nodes in the network")
	LOOKUPS  = flag.Int("lookups", 100, "number of lookups to test with")
)

func main() {
	flag.Parse()

	nodes := make([]u.Nnode, 0)

	// make concurrent joins :)
	var wg sync.WaitGroup
	var nodes_mux sync.Mutex

	for i := 0; i < *NODE_NUM; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			cmd := exec.Command("../../bin/koorde-overlay",
				"-port", fmt.Sprintf("%d", 8080+i),
				"-first", fmt.Sprint(u.IF(i == 0).Int(1, 0)),
				"-bootstrap", "127.0.0.1:8080")
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
	// time.Sleep(time.Second)

	for i := 0; i < *NODE_NUM; i++ {
		log.Print(node.ID(utils.SHA1OF(fmt.Sprintf("127.0.0.1:%d", 8080+i))).String())
	}

	for i := range nodes {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		peer, err := nodes[i].KC.DGetID(ctx, &pd.Empty{})
		if err != nil {
			panic(err)
		}
		nodes[i].ID = peer.SrcId
	}

	sort.Sort(&u.NodeStore{N: nodes})

	// check the successors
	// concurrent

	// go func() {
	// 	for {
	// var succ_wg sync.WaitGroup

	// 		for i := range nodes {
	// 			succ_wg.Add(1)

	// 			go func(i int) {
	// 				defer succ_wg.Done()

	// 				next_i := i + 1
	// 				if i == len(nodes)-1 {
	// 					next_i = 0
	// 				}
	// 				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// 				defer cancel()
	// 				reply, err := nodes[i].KC.GetSuccessorRPC(ctx, &pd.Empty{})
	// 				if err != nil {
	// 					panic(err)
	// 				}

	// 				log.Printf("%s -> %s & %s %v",
	// 					nodes[i].ID,
	// 					nodes[next_i].ID,
	// 					reply.SrcId,
	// 					reply.SrcId == nodes[next_i].ID)
	// 			}(i)
	// 		}
	// 		succ_wg.Wait()

	// 		// check the D pointers
	// 		// concurrent
	// 		var d_wg sync.WaitGroup

	// 		for i := range nodes {
	// 			d_wg.Add(1)

	// 			go func(i int) {
	// 				defer d_wg.Done()
	// 				N_id := node.ID(utils.ParseID(nodes[i].ID))
	// 				D_id, _ := N_id.LeftShift()

	// 				for j := range nodes {
	// 					var nnn u.Nnode
	// 					if j == len(nodes)-1 {
	// 						nnn = nodes[0]
	// 					} else {
	// 						nnn = nodes[j+1]
	// 					}

	// 					C_id := node.ID(utils.ParseID(nodes[j].ID))
	// 					NNN_id := node.ID(utils.ParseID(nnn.ID))

	// 					if D_id.InLXRange(C_id, NNN_id) {
	// 						ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// 						defer cancel()

	// 						reply, err := nodes[i].KC.DGetPointers(ctx, &pd.Empty{})
	// 						if err != nil {
	// 							panic(err)
	// 						}

	// 						log.Printf("%s ** %s & %s %v",
	// 							nodes[i].ID,
	// 							nodes[j].ID,
	// 							reply.D,
	// 							reply.D == nodes[j].ID)
	// 					}
	// 				}
	// 			}(i)
	// 		}
	// 		d_wg.Wait()

	// 		time.Sleep(2 * time.Second)
	// 	}
	// }()

	time.Sleep(30 * time.Second)
	log.Println("Start the lookups")

	pre := time.Now()
	var lookup_wg sync.WaitGroup

	for k := 0; k < *LOOKUPS; k++ {
		lookup_wg.Add(1)

		rand.Seed(time.Now().UnixNano())
		i, j := rand.Intn(*NODE_NUM), rand.Intn(*NODE_NUM)

		go func(i, j, k int) {
			defer lookup_wg.Done()

			pre := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			log.Printf("lookup %d from %s :%d -> %s :%d", k, nodes[i].ID, nodes[i].Port, nodes[j].ID, nodes[j].Port)
			reply, err := nodes[i].KC.DLKup(ctx, &pd.PeerPacket{SrcId: nodes[j].ID})
			if err != nil {
				panic(err)
			}

			log.Printf("lookup %d result: %v %v hi in %v", k, reply.SrcId, reply.SrcId == nodes[j].ID, time.Since(pre))
		}(i, j, k)
	}

	lookup_wg.Wait()

	log.Printf("%d lookups takes: %s", *LOOKUPS, time.Since(pre))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	for i := range nodes {
		nodes[i].Cmd.Wait()
	}

}
