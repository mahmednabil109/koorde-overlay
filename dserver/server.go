package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/mahmednabil109/koorde-overlay/dserver/drpc"
	"google.golang.org/grpc"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		WriteBufferSize: 1024,
	}
	UIPORT = flag.Int("ui-port", 8282, "ui port")
	PORT   = flag.Int("port", 16585, "port")
)

type DServer struct {
	drpc.UnimplementedDServerServer

	Data chan *drpc.DPointers
}

func (s *DServer) Init() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *PORT))
	if err != nil {
		panic(err)
	}
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

func main() {

	var uiConn *websocket.Conn
	var mux sync.Mutex

	dserver := DServer{
		Data: make(chan *drpc.DPointers),
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}
		mux.Lock()
		defer mux.Unlock()
		uiConn = conn

	})

	go func(s DServer) {
		for d := range s.Data {
			if uiConn == nil {
				continue
			}

			mux.Lock()
			err := uiConn.WriteJSON(struct {
				Id        string `json:"id"`
				Successor string `json:"successor"`
				D         string `json:"d"`
			}{
				Id:        d.Id,
				Successor: d.Successor,
				D:         d.D,
			})
			if err != nil {
				log.Print("unable to send update")
				uiConn.Close()
				return
			}
			mux.Unlock()
		}
	}(dserver)

	go dserver.Init()
	log.Printf("websocket server start listening %d", *UIPORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", *UIPORT), nil))

	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGINT)
	// <-sigs

	// uiConn.Close()
	// TODO clean shutdown to the http server
}
