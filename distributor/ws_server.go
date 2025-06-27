package distributor

import (
	"log"
	"net/http"
	"strconv"
	"work-distributor-system/models"

	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade the HTTP connection to a WebSocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // For development purposes
}

// WorkerWebSocket upgrades the HTTP connection to a WebSocket for a worker.
// It registers the worker and streams tasks to them as JSON.
func WorkerWebSocket(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("id")
	userID, _ := strconv.Atoi(userIDStr)

	// Upgrading connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	taskChan := make(chan models.Task) // channel to stream tasks to this worker
	RegisterWorker(userID, taskChan)   // register this worker and their task channel

	// Starting goroutine to send tasks over WebSocket
	go func() {
		for task := range taskChan {
			err := conn.WriteJSON(task)
			if err != nil {
				log.Println("Error sending task:", err)
				return
			}
		}
	}()

	// Reading from WebSocket to detect disconnection
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error (closing):", err)
			delete(workers, userID) // removing disconnected worker
			close(taskChan)
			return
		}
	}
}

// Storing active client WebSocket connections
var clientConns = make(map[uint]*websocket.Conn)

// ClientWebSocket sets up a WebSocket connection for clients to receive task and their status updates
func ClientWebSocket(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("id")
	userIDUint64, _ := strconv.Atoi(userIDStr)
	userID := uint(userIDUint64)

	/// Upgrade the client connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Client WebSocket upgrade error:", err)
		return
	}

	// Registering the client connection
	clientConns[userID] = conn
	log.Printf("Client %d connected via WebSocket", userID)

	// Keep the connection alive until error occurs
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Client %d disconnected", userID)
			conn.Close()
			delete(clientConns, userID)
			break
		}
	}
}

// ClientConn returns the WebSocket connection for a given client ID
func ClientConn(clientID uint) (*websocket.Conn, bool) {
	conn, ok := clientConns[clientID]
	return conn, ok
}
