package socket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Conn     *websocket.Conn
	Username string
	RoomID   string
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Room struct {
	Clients map[*Client]bool
	mu      sync.Mutex
}

var Rooms = make(map[string]*Room)
var roomsMu sync.Mutex

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{Conn: conn}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error for %s: %v", client.Username, err)
			LeaveRoom(client)
			break
		}

		var m Message
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}

		switch m.Type {
		case "join":
			payload, ok := m.Payload.(map[string]interface{})
			if !ok {
				log.Println("Invalid join payload")
				continue
			}
			roomID, _ := payload["room_id"].(string)
			username, _ := payload["username"].(string)
			if roomID == "" || username == "" {
				log.Println("Missing roomID or username")
				continue
			}

			// Special case: If HOST joined, create the room if it doesn't exist
			if username == "HOST" {
				roomsMu.Lock()
				if _, ok := Rooms[roomID]; !ok {
					Rooms[roomID] = &Room{
						Clients: make(map[*Client]bool),
					}
					log.Printf("Room %s created by HOST", roomID)
				}
				roomsMu.Unlock()
			}

			client.RoomID = roomID
			client.Username = username
			JoinRoom(client)
		case "start_game", "next_question", "show_results", "finish_game", "cancel_game":
			payload, ok := m.Payload.(map[string]interface{})
			if !ok {
				continue
			}
			roomID, _ := payload["room_id"].(string)
			BroadcastToRoom(roomID, m)
		case "submit_answer":
			payload, ok := m.Payload.(map[string]interface{})
			if ok {
				payload["username"] = client.Username
				m.Payload = payload
				BroadcastToRoom(client.RoomID, m)
			}
		}
	}
}

func LeaveRoom(client *Client) {
	if client.RoomID == "" {
		client.Conn.Close()
		return
	}

	roomsMu.Lock()
	room, ok := Rooms[client.RoomID]
	roomsMu.Unlock()

	if !ok {
		client.Conn.Close()
		return
	}

	room.mu.Lock()
	if _, exists := room.Clients[client]; exists {
		delete(room.Clients, client)
		log.Printf("Client %s left room %s", client.Username, client.RoomID)

		// Notify others if the leaving client is a player
		if client.Username != "HOST" {
			room.mu.Unlock() // Unlock before broadcasting to avoid deadlocks
			BroadcastToRoom(client.RoomID, Message{
				Type: "player_left",
				Payload: map[string]string{
					"username": client.Username,
				},
			})
		} else {
			room.mu.Unlock()
		}
	} else {
		room.mu.Unlock()
	}
	client.Conn.Close()
}

func JoinRoom(client *Client) {
	roomsMu.Lock()
	room, ok := Rooms[client.RoomID]
	roomsMu.Unlock()

	// 1. Check if room exists
	if !ok {
		log.Printf("Join failed: Room %s not found for client %s", client.RoomID, client.Username)
		client.Conn.WriteJSON(Message{
			Type: "error",
			Payload: map[string]string{
				"message": "Invalid Game PIN. Room not found.",
			},
		})
		client.Conn.Close()
		return
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	// 2. Check for duplicate username (if not HOST)
	if client.Username != "HOST" {
		for c := range room.Clients {
			if c.Username == client.Username {
				log.Printf("Join failed: Duplicate username %s in room %s", client.Username, client.RoomID)
				client.Conn.WriteJSON(Message{
					Type: "error",
					Payload: map[string]string{
						"message": "Username already taken in this room.",
					},
				})
				client.Conn.Close()
				return
			}
		}
	}

	room.Clients[client] = true
	// Collect existing players (excluding HOST)
	var players []string
	for c := range room.Clients {
		if c.Username != "HOST" {
			players = append(players, c.Username)
		}
	}

	log.Printf("Client %s joined room %s", client.Username, client.RoomID)

	// Send current player list to the joining client
	client.Conn.WriteJSON(Message{
		Type: "sync_players",
		Payload: map[string]interface{}{
			"players": players,
		},
	})

	// Only broadcast player_joined if not HOST
	if client.Username != "HOST" {
		// Broadcast to all clients in the room including the HOST
		msg := Message{
			Type: "player_joined",
			Payload: map[string]string{
				"username": client.Username,
			},
		}

		// Since we already have the room lock, we can't call BroadcastToRoom
		// because it will try to lock roomsMu and then room.mu (deadlock risk if other threads call handleConnection)
		// But BroadcastToRoom locks roomsMu first then room.mu.
		// Wait, JoinRoom is called with NO locks originally.
		// Let's manually broadcast here while we have room.mu lock.
		for c := range room.Clients {
			if c != client { // Don't send back to the joiner as they got sync_players
				c.Conn.WriteJSON(msg)
			}
		}
	}
}

func BroadcastToRoom(roomID string, msg Message) {
	roomsMu.Lock()
	room, ok := Rooms[roomID]
	roomsMu.Unlock()

	if !ok {
		return
	}

	room.mu.Lock()
	clients := make([]*Client, 0, len(room.Clients))
	for client := range room.Clients {
		clients = append(clients, client)
	}
	room.mu.Unlock()

	for _, client := range clients {
		if err := client.Conn.WriteJSON(msg); err != nil {
			log.Println("Write error:", err)
			LeaveRoom(client)
		}
	}
}
