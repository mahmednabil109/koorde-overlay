package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mahmednabil109/koorde-overlay/node"
	"github.com/mahmednabil109/koorde-overlay/utils"
)

var (
	port      = flag.Int("port", 16585, "port for grpc of the localnode")
	first     = flag.Bool("first", true, "flag to mark the node as the first in the network")
	bootstrap = flag.String("bootstrap", "127.0.0.1:16585", "ip for bootstraping node")
)

func main() {
	flag.Parse()

	var n node.Localnode
	if *first {
		err := n.Init(*port)
		if err != nil {
			log.Printf("faild to init the localnode: %v", err)
		}
	} else {
		n.Join(utils.ParseIP(*bootstrap), *port)
	}

	log.Printf("%+v \n", n)

	k := utils.ParseID("09b0ce42948043810a1f2cc7e7079aec7582f290")
	res, _ := n.Lookup(k, k, n.NodeAddr)
	log.Printf("%+v", res)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Printf("Programe ended")
}
