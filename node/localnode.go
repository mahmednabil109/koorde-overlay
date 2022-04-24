package node

import (
	"context"
	"fmt"
	"log"
	"net"
	"sort"
	"time"

	pd "github.com/mahmednabil109/koorde-overlay/rpc"
	"github.com/mahmednabil109/koorde-overlay/utils"
	"google.golang.org/grpc"
)

type Localnode struct {
	pd.UnimplementedKoordeServer
	Peer
	D           Peer
	DParents    []Peer
	Successor   Peer
	Predecessor Peer
	s           *grpc.Server
}

/* RPC impelementation */

func (ln *Localnode) BootStarpRPC(bctx context.Context, bootstrapPacket *pd.BootStrapPacket) (*pd.BootStrapReply, error) {
	src_id := ID(utils.ParseID(bootstrapPacket.SrcId))

	successor, err := ln.Lookup(src_id)
	if err != nil {
		return nil, err
	}

	d_id, _ := src_id.LeftShift()
	// lookup returns the successor
	d, err := ln.Lookup(d_id)
	if err != nil {
		return nil, err
	}
	// getting the predecessor
	d.InitConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d_pre, err := d.kc.GetPredecessorRPC(ctx, &pd.Empty{})
	if err != nil {
		return nil, err
	}

	return &pd.BootStrapReply{
		Successor: form_peer_packet(&successor),
		D:         d_pre,
	}, nil
}

func (ln *Localnode) LookupRPC(bctx context.Context, lookupPacket *pd.LookupPacket) (*pd.PeerPacket, error) {

	k := ID(utils.ParseID(lookupPacket.K))
	kShift := ID(utils.ParseID(lookupPacket.KShift))
	i := ID(utils.ParseID(lookupPacket.I))

	log.Printf("@=%+v,%+v", ln.NodeAddr, ln.Successor.NodeAddr)
	log.Printf("first %s in (%s %s] %v !!", k, ln.NodeAddr, ln.Successor.NodeAddr, k.InLXRange(ln.NodeAddr, ln.Successor.NodeAddr))

	if k.Equal(ln.NodeAddr) {
		log.Printf("Me || %s", ln.NetAddr)
		return form_peer_packet(&ln.Peer), nil
	}

	if k.InLXRange(ln.NodeAddr, ln.Successor.NodeAddr) {
		log.Printf("Successor || %s", ln.NetAddr)
		return form_peer_packet(&ln.Successor), nil
	}

	log.Printf("second %s in (%s %s] %v !!", i, ln.NodeAddr, ln.Successor.NodeAddr, i.InLXRange(ln.NodeAddr, ln.Successor.NodeAddr))

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

func (ln *Localnode) UpdatePredecessorRPC(bctx context.Context, p *pd.PeerPacket) (*pd.PeerPacket, error) {
	old_predecessor := form_peer_packet(&ln.Predecessor)
	ln.Predecessor = parse_peer_packet(p)

	return old_predecessor, nil
}

func (ln *Localnode) UpdateSuccessorRPC(bctx context.Context, p *pd.PeerListPacket) (*pd.PeerListPacket, error) {
	ln.Successor = parse_peer_packet(p.Peers[0])
	Succ_id := ID(utils.ParseID(p.Peers[1].SrcId))

	// slice of the pointers to hand out
	pointers := make([]*pd.PeerPacket, 0)
	for _, n := range ln.DParents {
		D_id, _ := n.NodeAddr.LeftShift()
		if D_id.InLXRange(ln.Successor.NodeAddr, Succ_id) {
			pointers = append(pointers, form_peer_packet(&n))
		}
	}
	log.Printf("hand over %v %d", pointers, len(ln.DParents))
	return &pd.PeerListPacket{Peers: pointers}, nil
}

func (ln *Localnode) UpdateDPointerRPC(bctx context.Context, p *pd.PeerPacket) (*pd.Empty, error) {
	ln.D = parse_peer_packet(p)

	return &pd.Empty{}, nil
}

func (ln *Localnode) AddDParentRPC(bctx context.Context, p *pd.PeerPacket) (*pd.Empty, error) {
	peer := parse_peer_packet(p)
	ln.DParents = append(ln.DParents, peer)

	// sort the Dparents Pointers by the id
	sort.Slice(ln.DParents, func(i, j int) bool {
		return !ln.DParents[i].NodeAddr.InLXRange(ln.DParents[j].NodeAddr, MAX_ID)
	})
	log.Print(ln.DParents)
	return &pd.Empty{}, nil
}

func (ln *Localnode) UpdateNeighborRPC(ctx context.Context, e *pd.Empty) (*pd.Empty, error) {
	return &pd.Empty{}, nil
}

func (ln *Localnode) BroadCastRPC(ctx context.Context, b *pd.BlockPacket) (*pd.Empty, error) {
	return &pd.Empty{}, nil
}

func (ln *Localnode) GetSuccessorRPC(ctx context.Context, e *pd.Empty) (*pd.PeerPacket, error) {
	// TODO make sure that the pointer is valid
	return form_peer_packet(&ln.Successor), nil
}

func (ln *Localnode) GetPredecessorRPC(ctx context.Context, e *pd.Empty) (*pd.PeerPacket, error) {
	// TODO make sure that the pointer is valid
	return form_peer_packet(&ln.Predecessor), nil
}

/* DEBUG RPC */

func (n *Localnode) DJoin(ctx context.Context, p *pd.PeerPacket) (*pd.Empty, error) {
	n.Join(utils.ParseIP(p.SrcIp), 8081)

	return &pd.Empty{}, nil
}

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
	ln.Predecessor = ln.Peer
	ln.DParents = []Peer{ln.Peer}
	ln.D = ln.Peer
	err := init_grpc_server(ln, port)
	return err
}

// Join initializes the node by executing Chord Join Algorithm
// it inits the Successor, D pointers
func (ln *Localnode) Join(nodeAddr *net.TCPAddr, port int) error {
	log.Printf("Join %s", nodeAddr.String())

	if ln.s == nil {
		ln.Init(port)
	}

	if nodeAddr == nil {
		return nil
	}

	peer := Peer{NetAddr: nodeAddr}
	peer.InitConnection()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	bootstrapPacket := &pd.BootStrapPacket{
		SrcId: ln.NodeAddr.String(),
		SrcIp: ln.NetAddr.String(),
	}
	bootstrapReply, err := peer.kc.BootStarpRPC(ctx, bootstrapPacket)
	if err != nil {
		log.Fatalf("cannot bootstrap: %v", err)
		return err
	}
	// init the successor and de brujin pointer
	ln.Successor = parse_peer_packet(bootstrapReply.Successor)
	ln.D = parse_peer_packet(bootstrapReply.D)

	// this must be before the successor initiation
	// to handle the case when i am my own d pointer
	// init connection with Dpointer and add itself as a Dparent
	currentNodePacket := form_peer_packet(&ln.Peer)
	err = ln.D.InitConnection()
	if err != nil {
		log.Fatal("fail to init connection with predecssor")
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ln.D.kc.AddDParentRPC(ctx, currentNodePacket)

	// init the connection of node pointers
	// TODO handle failer and pointer replacemnet
	// init connection with predecessor and update its pointers
	err = ln.Successor.InitConnection()
	if err != nil {
		log.Fatal("fail to init connection with successor")
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	predecessorPacket, err := ln.Successor.kc.UpdatePredecessorRPC(ctx, currentNodePacket)
	if err != nil {
		log.Fatal(err)
		return err
	}

	ln.Predecessor = parse_peer_packet(predecessorPacket)
	err = ln.Predecessor.InitConnection()
	if err != nil {
		log.Fatal("fail to init connection with predecessor")
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	packet := &pd.PeerListPacket{
		Peers: []*pd.PeerPacket{
			currentNodePacket,
			form_peer_packet(&ln.Successor),
		},
	}
	dParentsPacket, err := ln.Predecessor.kc.UpdateSuccessorRPC(ctx, packet)
	if err != nil {
		log.Fatal(err)
		return err
	}

	dParents := make([]Peer, 0)
	for i, p := range dParentsPacket.Peers {
		dParents = append(dParents, parse_peer_packet(p))
		// update DPointer in each dParent
		go func(p *Peer) {
			err := p.InitConnection()
			if err != nil {
				log.Fatal("error when init connection with Dparent")
				// best effort
				// return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			p.kc.UpdateDPointerRPC(ctx, currentNodePacket)

		}(&dParents[i])
	}
	ln.DParents = dParents

	return nil
}

func (ln *Localnode) Lookup(k ID) (Peer, error) {
	kShift, i := select_imaginary_node(k, ln.NodeAddr, ln.Successor.NodeAddr)
	log.Printf("init %s %s %s", k.String(), kShift.String(), i.String())

	lookupPacket := &pd.LookupPacket{
		SrcId:  ln.NodeAddr.String(),
		SrcIp:  ln.NetAddr.String(),
		K:      k.String(),
		KShift: kShift.String(),
		I:      i.String()}
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
		_id := m.MaskLowerWith(k, i).AddOne(i)

		if ID(_id).InLXRange(m, successor) {
			for j := 0; j < i; j++ {
				k, _ = k.LeftShift()
			}
			return k, _id
		}
	}
	// no Match
	return k, m.AddOne(0)
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
