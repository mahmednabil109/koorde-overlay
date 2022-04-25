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
	conn     *grpc.ClientConn
}

func (p *Peer) InitConnection() error {
	if p.kc != nil {
		return nil
	}

	conn, err := grpc.Dial(p.NetAddr.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can't dial: %v", err)
		return err
	}
	p.conn = conn

	p.kc = pd.NewKoordeClient(p.conn)
	log.Printf("connection Done With %s", p.NetAddr.String())
	return nil
}

func (p *Peer) CloseConnection() {
	p.kc = nil
	p.conn.Close()
}
