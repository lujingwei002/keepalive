package ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/lujingwei002/keepalive"
	"github.com/lujingwei002/keepalive/config"
)

type server struct {
	app    keepalive.Application
	config *config.Server
}

func New() keepalive.Server {
	return &server{}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *server) handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Client connected")
	app := s.app
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error while reading message:", err)
			return
		}
		fmt.Printf("Received message: %s\n", messageType)
		switch messageType {
		case websocket.BinaryMessage:
			app.HandleClientMessage(p)
			break
		}
		// Echo the message back to the client
		// err = conn.WriteMessage(messageType, p)
		// if err != nil {
		// fmt.Println("Error while writing message:", err)
		// return
		// }
	}
}

func (s *server) Init(app keepalive.Application, config *config.Server) error {
	if config.Ws == nil {
		return fmt.Errorf("server ws config not set")
	}
	c := config.Ws
	http.HandleFunc(c.Path, s.handleConnection)
	s.config = config
	s.app = app
	return nil

}

func (s *server) Start(app keepalive.Application) error {
	c := s.config.Ws
	fmt.Printf("WebSocket server started at %s\n", c.Bind)
	app.Go(func() error {
		err := http.ListenAndServe(c.Bind, nil)
		if err != nil {
			fmt.Println("Error starting server:", err)
		}
		return err
	})
	return nil
}
