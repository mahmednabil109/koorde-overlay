package node

import (
	"log"
	"math/rand"
	"sort"
	"testing"
	"time"
)

type nodeStore struct {
	nodes []Localnode
}

func (ns *nodeStore) Len() int {
	return len(ns.nodes)
}

func (ns *nodeStore) Swap(i, j int) {
	ns.nodes[i], ns.nodes[j] = ns.nodes[j], ns.nodes[i]
}

func (ns *nodeStore) Less(i, j int) bool {
	return !ns.nodes[i].NodeAddr.InLXRange(ns.nodes[j].NodeAddr, MAX_ID)
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// TODO #1 denug multi node in same memory space
func Test_lookup(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	NODE_LEN := 50
	ports := []int{}

	for i := 0; i < 50; i++ {
		ports = append(ports, 8080+i)
	}
	nodes, err := construct_simple_network(ports)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != NODE_LEN {
		t.FailNow()
	}

	for i := range nodes {
		if nodes[i].NodeAddr.Equal(nodes[i].Successor.NodeAddr) {
			t.Fatalf("node and its successor %d", i)
		}
	}

	for _, n := range nodes {
		log.Printf("%s => %s", n.NodeAddr, n.Successor.NodeAddr)
	}

	// for i := 0; i < 25; i++ {
	j, k := randInt(1, NODE_LEN), randInt(1, NODE_LEN)
	res, err := nodes[j].Lookup(nodes[k].NodeAddr)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("successor id %s", nodes[j].Successor.NodeAddr)
	log.Printf("lookup for %s, from %s", nodes[k].NodeAddr, nodes[j].NodeAddr)
	log.Printf("lookup (%d, %d) successed and peer: %+v", j, k, res.NodeAddr)
	// }

}

func construct_simple_network(network_ports []int) (nodes []Localnode, err error) {
	nodes = make([]Localnode, 0)

	for _, port := range network_ports {
		n := Localnode{}
		err = n.Init(port)
		if err != nil {
			return
		}
		nodes = append(nodes, n)
	}
	sort.Sort(&nodeStore{nodes})

	// inits the Successor pointer
	for i := range nodes {
		if i == len(nodes)-1 {
			nodes[i].Successor = nodes[0].Peer
		} else {
			nodes[i].Successor = nodes[i+1].Peer
		}
		nodes[i].Successor.InitConnection()
	}

	// inits the D pointer
	for _, n := range nodes {
		D_id, _ := n.NodeAddr.LeftShift()
		for i, nn := range nodes {
			var next_nn Localnode
			if i == len(nodes)-1 {
				next_nn = nodes[0]
			} else {
				next_nn = nodes[i+1]
			}
			if D_id.InLXRange(nn.NodeAddr, next_nn.NodeAddr) {
				n.D = nn.Peer
				break
			}
		}
		n.D.InitConnection()
	}

	return
}
