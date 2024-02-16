package service

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Service interface {
	SendMessageToUser(userID int32, payload []byte)
}

type NotificationService struct {
	clients *sync.Map
}

func New(clientMap *sync.Map) Service {
	return &NotificationService{
		clients: clientMap,
	}
}

// Function to send a message to a specific user
func (n *NotificationService) SendMessageToUser(userID int32, payload []byte) {
	// Lock the mutex before accessing the map
	if clientInterface, ok := n.clients.Load(userID); ok {
		client, ok := clientInterface.(*websocket.Conn)
		if !ok {
			// Log error if the type assertion fails
			log.Printf("Invalid type for client: %T", clientInterface)
		}

		err := client.WriteJSON(payload)
		if err != nil {
			log.Println("write:", err)
		}
	} else {
		log.Printf("User %v not found", userID)
	}
}
