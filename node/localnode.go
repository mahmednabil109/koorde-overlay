package node

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pd "github.com/mahmednabil109/koorde-overlay/rpc"
	"github.com/mahmednabil109/koorde-overlay/utils"
	"google.golang.org/grpc"
)

type Localnode struct {
	Peer
	D         Peer
	Successor Peer
	s         *grpc.Server
	pd.UnimplementedKoordeServer
}

/* RPC impelementation */

func (ln *Localnode) BootStarpRPC(ctx context.Context, bootstrapPacket *pd.BootStrapPacket) (*pd.BootStrapReply, error) {
	return &pd.BootStrapReply{}, nil
}

func (ln *Localnode) LookupRPC(bctx context.Context, lookupPacket *pd.LookupPacket) (*pd.PeerPacket, error) {

	k := ID(utils.ParseID(lookupPacket.K))
	kShift := ID(utils.ParseID(lookupPacket.KShift))
	i := ID(utils.ParseID(lookupPacket.I))

	log.Printf("@=%+v,%+v", ln.NodeAddr, ln.Successor.NodeAddr)
	log.Printf("inside %s in (%s %s] !!", k, ln.NodeAddr, ln.Successor.NodeAddr)
	if k.InLXRange(ln.NodeAddr, ln.Successor.NodeAddr) {
		log.Printf("Successor || %s", ln.NetAddr)
		return form_peer_packet(&ln.Successor), nil
	}
	if i.InLXRange(ln.NodeAddr, ln.Successor.NodeAddr) {
		log.Printf("Forward -> %s", ln.D.NetAddr)
		if ln.D.kc == nil {
			// TODO handle failer and pointer replacemnet
			ln.D.InitConnection()
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		KShift, _ := kShift.LeftShift()
		lookupPacket := &pd.LookupPacket{
			SrcId:  ln.NodeAddr.String(),
			SrcIp:  ln.NetAddr.String(),
			K:      k.String(),
			KShift: KShift.String(),
			I:      i.TopShift(kShift).String()}
		reply, err := ln.D.kc.LookupRPC(ctx, lookupPacket)

		if err != nil {
			log.Fatalf("lookup faild: %v", err)
			return nil, err
		}
		return reply, nil
	}
	log.Printf("Correction -> %s", ln.Successor.NetAddr)
	if ln.Successor.kc == nil {
		// TODO handle failer and pointer replacemnet
		ln.Successor.InitConnection()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reply, err := ln.Successor.kc.LookupRPC(ctx, lookupPacket)

	if err != nil {
		log.Fatalf("lookup faild: %v", err)
		return nil, err
	}
	return reply, nil
}

func (ln *Localnode) SuccessorRPC(ctx context.Context, e *pd.Empty) (*pd.PeerPacket, error) {
	return &pd.PeerPacket{}, nil
}

func (ln *Localnode) UpdateNeighborRPC(ctx context.Context, e *pd.Empty) (*pd.Empty, error) {
	return &pd.Empty{}, nil
}

func (ln *Localnode) BroadCastRPC(ctx context.Context, b *pd.BlockPacket) (*pd.Empty, error) {
	return &pd.Empty{}, nil
}

// DEBUG RPC

func (n *Localnode) DSetSuccessor(ctx context.Context, p *pd.PeerPacket) (*pd.Empty, error) {
	n.Successor = parse_peer_packet(p)
	n.Successor.InitConnection()
	return &pd.Empty{}, nil
}

func (n *Localnode) DSetD(ctx context.Context, p *pd.PeerPacket) (*pd.Empty, error) {
	n.D = parse_peer_packet(p)
	n.D.InitConnection()
	return &pd.Empty{}, nil
}

func (n *Localnode) DGetID(ctx context.Context, e *pd.Empty) (*pd.PeerPacket, error) {
	return &pd.PeerPacket{SrcId: n.NodeAddr.String()}, nil
}

func (n *Localnode) DGetPointers(ctx context.Context, e *pd.Empty) (*pd.Pointers, error) {
	return &pd.Pointers{Succ: n.Successor.NodeAddr.String(), D: n.D.NodeAddr.String()}, nil
}

func (n *Localnode) DLKup(ctx context.Context, p *pd.PeerPacket) (*pd.PeerPacket, error) {
	reply, err := n.Lookup(utils.ParseID(p.SrcId))
	return form_peer_packet(&reply), err
}

/* Localnode API */

// Init initializes the first node in the network
// it inits the Successor, D pointers with default values (node itslef)
func (ln *Localnode) Init(port int) error {
	ln.NetAddr = &net.TCPAddr{IP: []byte{127, 0, 0, 1}, Port: port}
	ln.NodeAddr = utils.SHA1OF(ln.NetAddr.String())
	ln.Start = ln.NodeAddr
	ln.Interval = []ID{ln.NodeAddr, ln.NodeAddr}
	ln.Successor = ln.Peer
	ln.D = ln.Peer
	err := init_grpc_server(ln, port)
	return err
}

// Join initializes the node by executing Chord Join Algorithm
// it inits the Successor, D pointers
func (ln *Localnode) Join(nodeAddr *net.TCPAddr, port int) error {
	ln.Init(port)
	if nodeAddr == nil {
		return nil
	}

	peer := Peer{NetAddr: nodeAddr}
	peer.InitConnection()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bootstrapPacket := &pd.BootStrapPacket{
		SrcId: ln.NodeAddr.String(),
		SrcIp: ln.NetAddr.String()}
	bootstrapReply, err := peer.kc.BootStarpRPC(ctx, bootstrapPacket)
	if err != nil {
		log.Fatalf("cannot bootstrap: %v", err)
		return err
	}
	// init the successor and de brujin pointer
	ln.Successor = parse_peer_packet(bootstrapReply.Successor)
	ln.D = parse_peer_packet(bootstrapReply.D)

	// init the connection of node pointers
	// TODO handle failer and pointer replacemnet
	ln.Successor.InitConnection()
	ln.D.InitConnection()

	return nil
}

func (ln *Localnode) Lookup(k ID) (Peer, error) {
	kShift, i := select_imaginary_node(k, ln.NodeAddr, ln.Successor.NodeAddr)

	lookupPacket := &pd.LookupPacket{
		SrcId:  ln.NodeAddr.String(),
		SrcIp:  ln.NetAddr.String(),
		K:      k.String(),
		KShift: kShift.String(),
		I:      i.TopShift(kShift).String()}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reply := make(chan *pd.PeerPacket)
	err := make(chan error)

	go func(ln *Localnode) {
		r, e := ln.LookupRPC(ctx, lookupPacket)
		reply <- r
		err <- e
	}(ln)

	r := <-reply
	e := <-err
	return parse_peer_packet(r), e
}

func (ln *Localnode) UpdateNeighbors() {}
func (ln *Localnode) UpdateOthers()    {}
func (ln *Localnode) BroadCast()       {}

/* Helper Methods */

// Select the best imaginary node to start the lookup from
// that is in the range (m, m.Successor] in the ring
func select_imaginary_node(k, m, successor ID) (ID, ID) {

	for i := 2*len(m) - 1; i >= 0; i-- {
		_id := m.MaskLowerWith(k, i)

		if ID(_id).InLXRange(m, successor) {
			for j := 0; j < i; j++ {
				k, _ = k.LeftShift()
			}
			return k, _id
		}
	}

	// no Match
	return k, m.AddOne()
}

// init_grpc_server creates a tcp socket and registers
// a new grpc server for Localnode.s
func init_grpc_server(ln *Localnode, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("faild to listen to %v", err)
		return err
	}

	ln.s = grpc.NewServer()
	pd.RegisterKoordeServer(ln.s, ln)
	go func() {
		log.Printf("grpc start listening %v", ln.NetAddr)
		if err := ln.s.Serve(lis); err != nil {
			log.Fatalf("faild to serve %v", err)
		}
	}()
	return nil
}

// parse_lookup_reply parses the pd.PeerPacket into a Peer struct
func parse_peer_packet(reply *pd.PeerPacket) Peer {
	if reply == nil {
		return Peer{}
	}
	return Peer{
		NodeAddr: utils.ParseID(reply.SrcId),
		NetAddr:  utils.ParseIP(reply.SrcIp),
		Start:    utils.ParseID(reply.Start),
		Interval: []ID{utils.ParseID(reply.Interval[0]), utils.ParseID(reply.Interval[1])},
	}
}

func form_peer_packet(peer *Peer) *pd.PeerPacket {
	if peer == nil {
		return nil
	}
	return &pd.PeerPacket{
		SrcId:    peer.NodeAddr.String(),
		SrcIp:    peer.NetAddr.String(),
		Start:    peer.Start.String(),
		Interval: []string{peer.Interval[0].String(), peer.Interval[1].String()},
	}
}
