package main

import (
	"context"
	"fmt"
	"log"
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

func main() {
	// test constants
	const (
		NODE_NUM = 20
	)

	nodes := make([]u.Nnode, 0)

	// make concurrent joins :)
	// var wg sync.WaitGroup
	// var nodes_mux sync.Mutex

	for i := 0; i < NODE_NUM; i++ {
		// wg.Add(1)
		// go func(i int) {
		// defer wg.Done()

		cmd := exec.Command("../../bin/koorde-overlay",
			"-port", fmt.Sprintf("%d", 8080+i),
			"-first", fmt.Sprint(u.IF(i == 0).Int(1, 0)),
			"-bootstrap", "127.0.0.1:8080")
		err := cmd.Start()
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)

		// nodes_mux.Lock()
		// defer nodes_mux.Unlock()
		nodes = append(nodes, u.Nnode{
			Cmd:  cmd,
			Port: 8080 + i,
		})
		nodes[len(nodes)-1].Init()
		// }(i)
	}
	// wg.Wait()
	// time.Sleep(time.Second)

	for i := 0; i < NODE_NUM; i++ {
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
	var succ_wg sync.WaitGroup

	for i := range nodes {
		succ_wg.Add(1)

		go func(i int) {
			defer succ_wg.Done()

			next_i := i + 1
			if i == len(nodes)-1 {
				next_i = 0
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			reply, err := nodes[i].KC.GetSuccessorRPC(ctx, &pd.Empty{})
			if err != nil {
				panic(err)
			}

			log.Printf("%s -> %s & %s %v",
				nodes[i].ID,
				nodes[next_i].ID,
				reply.SrcId,
				reply.SrcId == nodes[next_i].ID)
		}(i)
	}
	succ_wg.Wait()

	// check the D pointers
	// concurrent
	var d_wg sync.WaitGroup

	for i := range nodes {
		d_wg.Add(1)

		go func(i int) {
			defer d_wg.Done()
			N_id := node.ID(utils.ParseID(nodes[i].ID))
			D_id, _ := N_id.LeftShift()

			for j := range nodes {
				var nnn u.Nnode
				if j == len(nodes)-1 {
					nnn = nodes[0]
				} else {
					nnn = nodes[j+1]
				}

				C_id := node.ID(utils.ParseID(nodes[j].ID))
				NNN_id := node.ID(utils.ParseID(nnn.ID))

				if D_id.InLXRange(C_id, NNN_id) {
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()

					reply, err := nodes[i].KC.DGetPointers(ctx, &pd.Empty{})
					if err != nil {
						panic(err)
					}

					log.Printf("%s ** %s & %s %v",
						nodes[i].ID,
						nodes[j].ID,
						reply.D,
						reply.D == nodes[j].ID)
				}
			}
		}(i)
	}
	d_wg.Wait()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	for i := range nodes {
		nodes[i].Cmd.Wait()
	}

}
