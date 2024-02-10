package application

import (
	"flag"
	"log"
	"net/http"

	"github.com/BernardN38/ebuy-server/notification-service/handler"
)

var addr = flag.String("addr", "localhost:8083", "http service address")

type app struct {
}

func Run() {
	flag.Parse()
	log.SetFlags(0)

	h := handler.New()
	http.HandleFunc("/echo", h.Echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
