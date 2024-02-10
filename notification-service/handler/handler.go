package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Handler struct {
	upgrader websocket.Upgrader
}

func New() *Handler {
	return &Handler{
		upgrader: websocket.Upgrader{},
	}
}

type JsonPayload struct {
	Messsage string `json:"message"`
}

func (h *Handler) Echo(w http.ResponseWriter, r *http.Request) {
	upgrader := h.upgrader
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		var jsonPayload JsonPayload
		err := c.ReadJSON(&jsonPayload)
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", jsonPayload.Messsage)
		err = c.WriteJSON(jsonPayload)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
