package hub

import (
	"encoding/json"
	"log"
	"net/http"
	"sijil-core/internals/database"
	"time"

	"github.com/gorilla/websocket"
)

const  (

	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 512

)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	// Absoulutely crucial for security, in prod env
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

// Wrapper appround ws conn and the line
type Client struct {
	Hub *Hub
	Conn *websocket.Conn
	Send chan []byte 
	ProjectID int 
	userID int 
}

// -- This fxn reads from the websocket --
// pumps message from the ws conn, to the hub
func (c *Client) readPump() {

	 defer func ()  {
			c.Conn.Close()
			c.Hub.unregister <- c 
	 }()
	 c.Conn.SetReadLimit(maxMessageSize)
	 c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	 c.Conn.SetPongHandler(func(string) error {c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil})

	 for {

		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

	 }

}

// -- This function writes to the websocket --
// pumps messages from the hub to the ws conn
func (c *Client) writePump() {

	ticker := time.NewTicker(pingPeriod)
	// When this function ends (e.g connection breake, clean up)
	defer func ()  {
			ticker.Stop()
			c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <- c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <- ticker.C: 
			// Send a ping message to the client to keep conn alive
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil{
				return
			}
		}
	}

}

// Hub maintains the set of active clients and broadcast messages
type Hub struct {
	
	rooms map[int]map[*Client]bool 
	broadcast chan database.LogEntry
	register chan *Client
	unregister chan *Client
}

func NewHub() *Hub {

	return &Hub{
		rooms: make(map[int]map[*Client]bool),
		broadcast: make(chan database.LogEntry),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}

}

func (h *Hub) Run() {

	for  {

		select {

		case client := <- h.register: 
			if h.rooms[client.ProjectID] == nil {
				h.rooms[client.ProjectID] = make(map[*Client]bool)
			}

			h.rooms[client.ProjectID][client] = true 
		
		case client := <- h.unregister: 
			if room, ok := h.rooms[client.ProjectID]; ok {
				if _, ok := room[client]; ok {
					delete(room, client)
					close(client.Send)
				}

				if len(room) == 0 {
					delete(h.rooms, client.ProjectID)
				}
			}
		case logEntry := <- h.broadcast: 

			room := h.rooms[logEntry.ProjectID]

			if room == nil {
				continue
			}
			message, _ := json.Marshal(logEntry)

			for client := range  room {
				select {
				case client.Send <- message:
					// Message sent successfully
				default: 
					close(client.Send)
					delete(room, client)
				}
			}
		}
	}
}

func (h *Hub) BroadcastLog(logEntry database.LogEntry) {

	h.broadcast <- logEntry
}

// This is the "Front door" for our server handler
// handles ws request from the peer
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, projectID, userID int) {

	// upgrade http to ws 
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create the new client
	client := &Client{
		 Hub: hub, 
		 Conn: conn,
		 Send: make(chan []byte, 256),
		ProjectID: projectID,
		userID: userID,
		}

	// Register client
	client.Hub.register <- client

	// Thse goroutines run in the background for life of the conn
	go client.writePump()
	go client.readPump()
};