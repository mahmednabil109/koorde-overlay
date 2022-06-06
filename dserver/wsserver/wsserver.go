package wsserver

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/mahmednabil109/koorde-overlay/dserver/drpc"
)

type ACTION_FUNC func(string, string)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	WriteBufferSize: 1024,
}

type response struct {
	Node   string `json:"id"`
	Action int    `json:"action"`
	Block  string `json:"block"`
}

type notification struct {
	Id        string   `json:"id"`
	Successor string   `json:"successor"`
	D         string   `json:"d"`
	Blocks    []string `json:"blocks"`
}

type Server struct {
	UIConn  *websocket.Conn
	Mux     sync.Mutex
	Actions map[int]ACTION_FUNC
	Data    <-chan *drpc.DPointers
	Blocks  <-chan []string

	updateStop chan bool
}

func (s *Server) InitAndServe(actions map[int]ACTION_FUNC, data <-chan *drpc.DPointers, blocks <-chan []string, port int) {

	// init the data chan
	s.Data = data
	s.Blocks = blocks

	// init actions
	s.Actions = actions

	// handle upgrade
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Panic(err)
		}
		s.Mux.Lock()
		s.UIConn = conn
		s.Mux.Unlock()

		log.Print("Init ui viewer")
		// initiate the listener and the viewer
		go s.listen()
	})

	// the purpose is to make only one controller & one view
	// when new ui is connected stop others
	go s.update()

	log.Printf("websocket server start listening %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil))
}

// send updates to the ui from the nodes
func (s *Server) update() {
	for {
		select {
		case d := <-s.Data:
			if s.UIConn == nil {
				continue
			}

			s.Mux.Lock()
			err := s.UIConn.WriteJSON(notification{
				Id:        d.Id,
				Successor: d.Successor,
				D:         d.D,
			})
			s.Mux.Unlock()

			if err != nil {
				log.Print("unable to send update")
				s.UIConn.Close()
				continue
			}
		case b := <-s.Blocks:
			if s.UIConn == nil {
				continue
			}

			s.Mux.Lock()
			err := s.UIConn.WriteJSON(notification{
				Blocks: b,
			})
			s.Mux.Unlock()

			if err != nil {
				log.Print("unable to send update")
				s.UIConn.Close()
				continue
			}
		case <-s.updateStop:
			return
		}
	}
}

// listen to any action from the ui
func (s *Server) listen() {
	for {
		var reply response
		err := s.UIConn.ReadJSON(&reply)
		if err != nil {
			log.Print(err)
			break
		}
		// logging
		log.Print(reply)

		// if the action is defined by the dserver
		// execute it
		if action, ok := s.Actions[reply.Action]; ok {
			log.Print(action)
			action(reply.Node, reply.Block)
		} else {
			log.Print("action not found")
		}

	}
}
