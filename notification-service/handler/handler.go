package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/jwtauth/v5"

	"github.com/gorilla/websocket"
)

type Handler struct {
	upgrader websocket.Upgrader
	clients  *sync.Map
}

func New(clientMap *sync.Map) *Handler {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			clientMap.Range(func(key, value interface{}) bool {
				fmt.Printf("Key: %v\n", key)
				return true // continue iterating
			})
		}
	}()
	return &Handler{
		upgrader: websocket.Upgrader{},
		clients:  clientMap,
	}
}

type JsonPayload struct {
	UserId  int32  `json:"userId"`
	Message string `json:"message"`
}

func (h *Handler) Echo(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	userID, ok := claims["user_id"].(float64)
	if !ok {
		log.Print("invalid cookie user id")
		http.Error(w, "bad token", http.StatusBadRequest)
		return
	}
	upgrader := h.upgrader
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	// Lock the mutex before accessing the map

	h.clients.Store(int32(userID), c)

	for {
		var jsonPayload JsonPayload
		err := c.ReadJSON(&jsonPayload)
		if err != nil {
			log.Println("read:", err)
			// break
			continue
		}
		log.Printf("recv: %s", jsonPayload.Message)
		// Example of targeting a specific websocket based on userID
		h.sendMessageToUser(int32(jsonPayload.UserId), []byte("hello there"))
	}
}

// Function to send a message to a specific user
func (h *Handler) sendMessageToUser(userID int32, payload []byte) {
	// Lock the mutex before accessing the map
	if clientInterface, ok := h.clients.Load(userID); ok {
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
