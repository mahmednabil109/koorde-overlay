package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/mahmednabil109/koorde-overlay/dserver/drpc"
	"github.com/mahmednabil109/koorde-overlay/dserver/wsserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ACTION_GETBLOCKS = iota
	ACTION_INITBRAOD
)

var (
	RetryPolicy = `{
		"methodConfig": [{
		  "name": [{"service": "drpc.DServer"}],
		  "waitForReady": true,
		  "retryPolicy": {
			  "MaxAttempts": 4,
			  "InitialBackoff": ".01s",
			  "MaxBackoff": ".01s",
			  "BackoffMultiplier": 1.0,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ]
		  }
		}]}`
	UIPORT = flag.Int("ui-port", 8282, "ui port")
	PORT   = flag.Int("port", 16585, "port")
)

type DServer struct {
	drpc.UnimplementedDServerServer

	conns  map[string]drpc.DServerClient
	Data   chan *drpc.DPointers
	Blocks chan []string
}

func (s *DServer) Init() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *PORT))
	if err != nil {
		panic(err)
	}

	s.conns = make(map[string]drpc.DServerClient)
	server := grpc.NewServer()
	drpc.RegisterDServerServer(server, s)
	log.Printf("grpc start listening %v", *PORT)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}

func (s *DServer) UpdatePointers(ctx context.Context, pointers *drpc.DPointers) (*drpc.DEmpty, error) {
	go func() { s.Data <- pointers }()
	return &drpc.DEmpty{}, nil
}

func (s *DServer) Connect(ctx context.Context, p *drpc.DPeer) (*drpc.DEmpty, error) {
	s.conns[p.Id] = init_grpc_client(p.Addr)
	log.Print(s.conns)

	return &drpc.DEmpty{}, nil
}

func init_grpc_client(ip string) drpc.DServerClient {
	log.Printf("connect to %s", ip)

	conn, err := grpc.Dial(
		ip,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(RetryPolicy),
	)
	if err != nil {
		panic(err)
	}
	return drpc.NewDServerClient(conn)
}

func handleGetBlocks(s DServer, id, _ string) {
	log.Print("get blocks called")
	client, ok := s.conns[id]
	if !ok {
		log.Printf("%s Not found", id)
		return
	}

	reply, err := client.GetLocalBlock(context.Background(), &drpc.DEmpty{})
	log.Print(reply)
	if err != nil {
		log.Print(err)
	}

	go func(r *drpc.DBlocks) {
		blocks := []string{}
		for _, b := range r.Blocks {
			blocks = append(blocks, b.Bid)
		}
		s.Blocks <- blocks
	}(reply)
}

func handleInitBroadcast(s DServer, id, info string) {
	client, ok := s.conns[id]
	if !ok {
		log.Printf("%s Not found", id)
		return
	}

	_, err := client.BroadcastBlock(context.Background(), &drpc.DBlock{Bid: info})
	if err != nil {
		log.Print(err)
	}
}

func main() {

	var wserver wsserver.Server

	dataChan := make(chan *drpc.DPointers)
	blocksChan := make(chan []string)

	dserver := DServer{
		Data:   dataChan,
		Blocks: blocksChan,
	}

	go dserver.Init()
	wserver.InitAndServe(
		map[int]wsserver.ACTION_FUNC{
			ACTION_GETBLOCKS: func(id, info string) { handleGetBlocks(dserver, id, info) },
			ACTION_INITBRAOD: func(id, info string) { handleInitBroadcast(dserver, id, info) },
		},
		dataChan,
		blocksChan,
		*UIPORT,
	)

	// TODO clean shutdown to the http server
}
