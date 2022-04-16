package node

import (
	"log"
	"net"

	pd "github.com/mahmednabil109/koorde-overlay/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Peer struct {
	NetAddr  *net.TCPAddr
	NodeAddr ID
	Start    ID
	Interval []ID
	kc       pd.KoordeClient
}

func (p *Peer) InitConnection() error {
	conn, err := grpc.Dial(p.NetAddr.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can't dial: %v", err)
		return err
	}

	p.kc = pd.NewKoordeClient(conn)
	return nil
}
