package client

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (

	// time allowed to write message to peer
	writeWait = 10 * time.Second

	// time allowed to read message from peer
	pongWait = 60 * time.Second

	// ping period must be less from pongPeriod
	pingPeriod = (pongWait * 9) / 10

	// size allowed for message size
	maxMessageSize = 512 // kb
)

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Conn struct {

	// connection
	ws *websocket.Conn

	// send msg (in byte type)
	send chan []byte
}

// func read messages from the ws connects to the hub
func (s subscriptions) readPump(h *Hub) {
	c := s.conn
	defer func() {
		h.unregister <- s
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(appData string) error {
		err := c.ws.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}

		return nil
	})

	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("---Error-> IsUnexpectedCloseError--->", err)
			}
			break
		}

		log.Println("msg: ", string(msg))
		m := message{msg, s.room}
		h.broadcast <- m
	}
}

func (c *Conn) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (s *subscriptions) writePump() {
	c := s.conn

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {

		select {
		case msg, ok := <-c.send:
			if !ok {
				if err := c.write(websocket.CloseMessage, []byte{}); err != nil {
					fmt.Println("Error: ", err.Error())
					return
				}
			}

			log.Println("msg: ", string(msg))

			if err := c.write(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func ServeWS(w http.ResponseWriter, r *http.Request, roomId string, h *Hub) {
	fmt.Print(roomId)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	c := &Conn{send: make(chan []byte, 256), ws: ws}
	s := subscriptions{c, roomId}
	h.register <- s
	go s.writePump()
	go s.readPump(h)
}
