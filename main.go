package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mahmednabil109/koorde-overlay/mock"
	"github.com/mahmednabil109/koorde-overlay/node"
	"github.com/mahmednabil109/koorde-overlay/utils"
)

var (
	port      = flag.Int("port", 16585, "port for grpc of the localnode")
	first     = flag.Int("first", 1, "flag to mark the node as the first in the network")
	bootstrap = flag.String("bootstrap", "127.0.0.1:16585", "ip for bootstraping node")
	d         = flag.Int("debug", 0, "debug flag")
	dserver   = flag.String("dserver", "127.0.0.1:16585", "address of the dserver")
	dport     = flag.Int("dport", 16586, "dport")
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	var n node.Localnode
	n.ConsensusAPI = &mock.Consensus{}

	InitContext := context.WithValue(context.Background(), "port", *port)
	InitContext = context.WithValue(InitContext, "d", *d)
	InitContext = context.WithValue(InitContext, "dserver-addr", *dserver)
	InitContext = context.WithValue(InitContext, "dport", *dport)

	log.Printf("%+v", InitContext)

	err := n.Init(InitContext)
	if err != nil {
		log.Printf("faild to init the localnode: %v", err)
	}
	f, err := os.Create(fmt.Sprintf("%s.log", n.NetAddr.String()))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)

	if *first == 0 {
		log.Println("start to join")
		n.Join(utils.ParseIP(*bootstrap), *port)
		log.Printf("nodeID %v", n)
		log.Printf("Successor %v", n.Successor)
		log.Printf("D %v", n.D)
		log.Printf("Predecessor %v", n.Predecessor)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs

	// n.NodeShutdown <- true
	log.Printf("Programe ended")
}
