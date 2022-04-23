package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mahmednabil109/koorde-overlay/node"
	"github.com/mahmednabil109/koorde-overlay/utils"
)

var (
	port      = flag.Int("port", 16585, "port for grpc of the localnode")
	first     = flag.Int("first", 1, "flag to mark the node as the first in the network")
	bootstrap = flag.String("bootstrap", "127.0.0.1:16585", "ip for bootstraping node")
)

func main() {
	flag.Parse()

	var n node.Localnode
	err := n.Init(*port)
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
		log.Printf("nodeID %s", n.NodeAddr.String())
		log.Printf("Successor %s", n.Successor.NodeAddr.String())
		log.Printf("D %s", n.D.NodeAddr.String())
		log.Printf("Predecessor %s", n.Predecessor.NodeAddr.String())
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs

	log.Printf("Programe ended")
}
