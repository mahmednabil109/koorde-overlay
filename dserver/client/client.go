package client

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/mahmednabil109/koorde-overlay/dserver/drpc"
	"github.com/mahmednabil109/koorde-overlay/mock"
	"github.com/mahmednabil109/koorde-overlay/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
)

type Node interface {
	BroadCast(string) error
	GetLocalBlocks() []mock.Block
	GetID() string
}

type Client struct {
	Node Node

	// client part
	dc   drpc.DServerClient
	conn *grpc.ClientConn
	// server part
	Port int
	s    *grpc.Server
	drpc.UnimplementedDServerServer
}

func (c *Client) Init(node Node, addr string, port int) error {
	c.Node = node
	c.Port = port

	if c.conn != nil {
		return nil
	}

	// init the client comm
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(RetryPolicy),
	)
	if err != nil {
		log.Print(err)
		return err
	}

	c.conn = conn
	c.dc = drpc.NewDServerClient(c.conn)
	log.Printf("conn done with %s", addr)

	// send connection data to the dserver
	ctx := context.Background()
	_, err = c.dc.Connect(ctx, &drpc.DPeer{
		Addr: fmt.Sprintf("%s:%d", utils.GetMyIP().String(), c.Port),
		Id:   c.Node.GetID(),
	})
	log.Print(err)

	// init the server comm
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Port))
	if err != nil {
		log.Fatalf("faild to listen to %v", err)
		return err
	}
	c.s = grpc.NewServer()
	drpc.RegisterDServerServer(c.s, c)
	go func() {
		log.Printf("grpc start listening 16585")
		if err := c.s.Serve(lis); err != nil {
			log.Fatalf("faild to serve %v", err)
		}
	}()

	return nil
}

func (c *Client) Update(id, successor, d string) {
	_, err := c.dc.UpdatePointers(context.Background(), &drpc.DPointers{Id: id, Successor: successor, D: d})
	if err != nil {
		log.Print(err)
	}
}

// grpc methods
func (c *Client) GetLocalBlock(ctx context.Context, e *drpc.DEmpty) (*drpc.DBlocks, error) {
	var blocks []*drpc.DBlock

	// parse all blocks
	for _, b := range c.Node.GetLocalBlocks() {
		blocks = append(blocks, &drpc.DBlock{Bid: b.Info})
	}

	return &drpc.DBlocks{Blocks: blocks}, nil
}

func (c *Client) BroadcastBlock(ctx context.Context, e *drpc.DBlock) (*drpc.DEmpty, error) {
	err := c.Node.BroadCast(e.Bid)
	return &drpc.DEmpty{}, err
}
