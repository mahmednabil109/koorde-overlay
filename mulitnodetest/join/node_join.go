package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sort"
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
		NODE_NUM = 100
	)

	nodes := make([]u.Nnode, 0)

	for i := 0; i < NODE_NUM; i++ {

		cmd := exec.Command("../../bin/koorde-overlay",
			"-port", fmt.Sprintf("%d", 8080+i),
			"-first", fmt.Sprint(u.IF(i == 0).Int(1, 0)),
			"-bootstrap", "127.0.0.1:8080")
		err := cmd.Start()
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)

		nodes = append(nodes, u.Nnode{
			Cmd:  cmd,
			Port: 8080 + i,
		})
		nodes[i].Init()
	}

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
	for i := range nodes {
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
	}

	// check the D pointers
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
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	for i := range nodes {
		nodes[i].Cmd.Wait()
	}

}
