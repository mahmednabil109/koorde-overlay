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
	D         *Peer
	Successor *Peer
	s         *grpc.Server
	pd.UnimplementedKoordeServer
}

/* RPC impelementation */

func (ln *Localnode) BootStarpRPC(ctx context.Context, bootstrapPacket *pd.BootStrapPacket) (*pd.BootStrapReply, error) {
	return &pd.BootStrapReply{}, nil
}

func (ln *Localnode) LookupRPC(ctx context.Context, lookupPacket *pd.LookupPacket) (*pd.PeerPacket, error) {
	return &pd.PeerPacket{}, nil
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

/* Localnode API */

// Init initializes the first node in the network
// It inits the Successor, D pointers with default values (node itslef)
func (ln *Localnode) Init(port int) error {
	ln.NetAddr = &net.TCPAddr{IP: []byte{127, 0, 0, 1}, Port: port}
	ln.NodeAddr = utils.SHA1OF(ln.NetAddr.String())
	ln.Start = ln.NodeAddr
	ln.Interval = []ID{ln.NodeAddr, ln.NodeAddr}
	ln.Successor = &ln.Peer
	ln.D = &ln.Peer
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

func (ln *Localnode) Lookup(id, idShift, iid ID) (*Peer, error) {
	if id.InLXRange(ln.NodeAddr, ln.Successor.NodeAddr) {
		return ln.Successor, nil
	}
	if iid.InLXRange(ln.NodeAddr, ln.Successor.NodeAddr) {
		if ln.D.kc == nil {
			ln.D.InitConnection()
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_idShift, _ := idShift.LeftShift()
		lookupPacket := &pd.LookupPacket{
			SrcId:   ln.NodeAddr.String(),
			SrcIp:   ln.NetAddr.String(),
			Id:      id.String(),
			IdShift: _idShift.String(),
			Iid:     iid.TopShift(idShift).String()}
		reply, err := ln.D.kc.LookupRPC(ctx, lookupPacket)

		if err != nil {
			log.Fatalf("lookup faild: %v", err)
			return nil, err
		}
		return parse_peer_packet(reply), nil
	}

	if ln.Successor.kc == nil {
		ln.Successor.InitConnection()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	lookupPacket := &pd.LookupPacket{
		SrcId:   ln.NetAddr.String(),
		SrcIp:   ln.NetAddr.String(),
		Id:      id.String(),
		IdShift: idShift.String(),
		Iid:     iid.String()}
	reply, err := ln.Successor.kc.LookupRPC(ctx, lookupPacket)

	if err != nil {
		log.Fatalf("lookup faild: %v", err)
		return nil, err
	}
	return parse_peer_packet(reply), nil
}

func (ln *Localnode) UpdateNeighbors() {}
func (ln *Localnode) UpdateOthers()    {}
func (ln *Localnode) BroadCast()       {}

/* Helper Methods */

// init_grpc_server creates a tcp socket an registers
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
func parse_peer_packet(reply *pd.PeerPacket) *Peer {
	return &Peer{
		NodeAddr: utils.ParseID(reply.SrcId),
		NetAddr:  utils.ParseIP(reply.SrcIp),
		Start:    utils.ParseID(reply.Start),
		Interval: []ID{utils.ParseID(reply.Interval[0]), utils.ParseID(reply.Interval[1])},
	}
}
