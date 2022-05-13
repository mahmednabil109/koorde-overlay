package mulitnodetest

import (
	"fmt"
	"log"
	"os/exec"

	pd "github.com/mahmednabil109/koorde-overlay/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	retryPolicy = `{
		"methodConfig": [{
		  "name": [{"service": "rpc.Koorde"}],
		  "waitForReady": true,
		  "retryPolicy": {
			  "MaxAttempts": 4,
			  "InitialBackoff": ".01s",
			  "MaxBackoff": ".01s",
			  "BackoffMultiplier": 1.0,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ]
		  }
		}]}`
)

type NodeStore struct {
	N []Nnode
}

func (ns *NodeStore) Len() int {
	return len(ns.N)
}

func (ns *NodeStore) Swap(i, j int) {
	ns.N[i], ns.N[j] = ns.N[j], ns.N[i]
}

func (ns *NodeStore) Less(i, j int) bool {
	return ns.N[i].ID < ns.N[j].ID
}

type Nnode struct {
	Port int
	Cmd  *exec.Cmd
	ID   string
	KC   pd.KoordeClient
}

func (n *Nnode) Init() error {

	conn, err := grpc.Dial(
		fmt.Sprintf("127.0.0.1:%d", n.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(retryPolicy))

	if err != nil {
		log.Fatalf("can't dial: %v", err)
		return err
	}

	n.KC = pd.NewKoordeClient(conn)
	log.Printf("connection Done With %s", fmt.Sprintf("127.0.0.1:%d", n.Port))
	return nil
}

type IF bool

func (c IF) Int(a, b int) int {
	if c {
		return a
	}
	return b
}
