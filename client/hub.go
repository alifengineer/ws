package client

type message struct {
	data []byte
	room string
}

type subscriptions struct {
	conn *Conn
	room string
}

type Hub struct {

	// hub rooms
	rooms map[string]map[*Conn]bool

	// Inboudn messages from the connections
	broadcast chan message

	// register requests from the connections
	register chan subscriptions

	// unregister requests from th connections
	unregister chan subscriptions
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan message),
		register:   make(chan subscriptions),
		unregister: make(chan subscriptions),
		rooms:      make(map[string]map[*Conn]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case r := <-h.register:
			connections := h.rooms[r.room]
			if connections == nil {
				connections = make(map[*Conn]bool)
				h.rooms[r.room] = connections
			}
			h.rooms[r.room][r.conn] = true
		case ur := <-h.unregister:
			connections := h.rooms[ur.room]
			if connections != nil {
				if _, ok := connections[ur.conn]; ok {
					delete(connections, ur.conn)
					close(ur.conn.send)
					if len(connections) == 0 {
						delete(h.rooms, ur.room)
					}
				}
			}

		case m := <-h.broadcast:
			connections := h.rooms[m.room]
			for c := range connections {
				select {
				case c.send <- m.data:
					//
				default:
					close(c.send)
					delete(connections, c)

					if len(connections) == 0 {
						delete(h.rooms, m.room)
					}
				}
			}
		}
	}
}
