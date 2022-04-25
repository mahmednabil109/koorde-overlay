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
	NODE_NUM    = flag.Int("node-num", 100, "number of nodes in the system")
	LOOKUPS_NUM = flag.Int("lookups", 10000, "number of lookups to test with")
)

func main() {
	nodes := make([]u.Nnode, 0)

	// init the nodes paralle
	var wg sync.WaitGroup
	var nodes_mux sync.Mutex

	for i := 0; i < *NODE_NUM; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			cmd := exec.Command("../../bin/koorde-overlay", "-port", fmt.Sprintf("%d", 8080+i))
			f, err := os.Create(fmt.Sprintf("%d.err.log", 8080+i))
			if err != nil {
				panic(err)
			}
			cmd.Stderr = f
			err = cmd.Start()
			time.Sleep(1 * time.Second)
			if err != nil {
				panic(err)
			}

			nodes_mux.Lock()
			defer nodes_mux.Unlock()

			nodes = append(nodes, u.Nnode{
				Cmd:  cmd,
				Port: 8080 + i,
			})

			// initConnection
			err = nodes[len(nodes)-1].Init()
			if err != nil {
				panic(err)
			}

		}(i)
	}
	wg.Wait()

	for i := range nodes {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		peer, err := nodes[i].KC.DGetID(ctx, &pd.Empty{})
		if err != nil {
			panic(err)
		}
		nodes[i].ID = peer.SrcId
	}
	sort.Sort(&u.NodeStore{N: nodes})

	for i := range nodes {
		log.Printf("%s", nodes[i].ID)
	}

	// SET Successor POINTERS
	for i := range nodes {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		nni := i + 1
		if i == len(nodes)-1 {
			nni = 0
		}

		pp := &pd.PeerPacket{
			SrcId:    nodes[nni].ID,
			SrcIp:    fmt.Sprintf("127.0.0.1:%d", nodes[nni].Port),
			Interval: []string{nodes[nni].ID, nodes[nni].ID},
			Start:    nodes[nni].ID,
		}
		_, err := nodes[i].KC.DSetSuccessor(ctx, pp)
		if err != nil {
			panic(err)
		}
	}

	// SET D POINTERS
	for i := range nodes {
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

				pp := &pd.PeerPacket{
					SrcId:    nodes[j].ID,
					SrcIp:    fmt.Sprintf("127.0.0.1:%d", nodes[j].Port),
					Interval: []string{nodes[j].ID, nodes[j].ID},
					Start:    nodes[j].ID,
				}
				_, err := nodes[i].KC.DSetD(ctx, pp)
				if err != nil {
					panic(err)
				}
			}

		}

	}

	for i := range nodes {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := nodes[i].KC.DGetPointers(ctx, &pd.Empty{})
		if err != nil {
			panic(err)
		}
		log.Printf("%s => %s ** %s", nodes[i].ID, res.Succ, res.D)
	}

	// make concurrent lookups
	var lookup_wg sync.WaitGroup

	pre := time.Now()
	for k := 0; k < *LOOKUPS_NUM; k++ {
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

			log.Printf("lookup %d result: %v hi in %v", k, reply.SrcId, time.Since(pre))
		}(i, j, k)

	}
	lookup_wg.Wait()

	log.Printf("%d lookups takes: %s", *LOOKUPS_NUM, time.Since(pre))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	for i := range nodes {
		nodes[i].Cmd.Wait()
	}
}
