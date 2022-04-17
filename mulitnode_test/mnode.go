package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/mahmednabil109/koorde-overlay/node"
	pd "github.com/mahmednabil109/koorde-overlay/rpc"
	"github.com/mahmednabil109/koorde-overlay/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type nodeStore struct {
	n []nnode
}

func (ns *nodeStore) Len() int {
	return len(ns.n)
}

func (ns *nodeStore) Swap(i, j int) {
	ns.n[i], ns.n[j] = ns.n[j], ns.n[i]
}

func (ns *nodeStore) Less(i, j int) bool {
	return ns.n[i].ID < ns.n[j].ID
}

type nnode struct {
	Port int
	Cmd  *exec.Cmd
	ID   string
	KC   pd.KoordeClient
}

func (n *nnode) Init() error {

	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", n.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can't dial: %v", err)
		return err
	}

	n.KC = pd.NewKoordeClient(conn)
	log.Printf("connection Done With %s", fmt.Sprintf("127.0.0.1:%d", n.Port))
	return nil
}

func main() {
	nodes := make([]nnode, 0)

	for i := 0; i < 100; i++ {
		cmd := exec.Command("../bin/koorde-overlay", "-port", fmt.Sprintf("%d", 8080+i))
		err := cmd.Start()
		time.Sleep(1 * time.Second)
		if err != nil {
			panic(err)
		}

		nodes = append(nodes, nnode{
			Cmd:  cmd,
			Port: 8080 + i,
		})

		// initConnection
		nodes[i].Init()
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
	sort.Sort(&nodeStore{nodes})

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
			var nnn nnode
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

	pre := time.Now()
	for k := 1; k < 1000; k++ {

		rand.Seed(time.Now().UnixNano())
		i, j := rand.Intn(19), rand.Intn(19)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		log.Printf("lookup from %s :%d -> %s :%d", nodes[i].ID, nodes[i].Port, nodes[j].ID, nodes[j].Port)
		reply, err := nodes[i].KC.DLKup(ctx, &pd.PeerPacket{SrcId: nodes[j].ID})
		if err != nil {
			panic(err)
		}

		log.Printf("lookup result: %+v", reply)
	}
	log.Printf("1000 lookups takes: %s", time.Since(pre))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	for i := range nodes {
		nodes[i].Cmd.Wait()
	}
}
