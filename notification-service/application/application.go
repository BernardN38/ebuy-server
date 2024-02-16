package application

import (
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/BernardN38/ebuy-server/notification-service/handler"
	"github.com/BernardN38/ebuy-server/notification-service/messaging"
	"github.com/BernardN38/ebuy-server/notification-service/service"
	"github.com/go-chi/jwtauth/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

var addr = flag.String("addr", "localhost:8083", "http service address")
var jwtSecret = flag.String("jwtSecret", "qwertyuiopasdfghjklzxcvbnm123456qwertyuiopasdfghjklzxcvbnm123456", "http service address")
var rabbitUrl = flag.String("rabbitUrl", "amqp://guest:guest@localhost", "http service address")

type app struct {
}

func Run() {
	flag.Parse()
	log.SetFlags(0)
	clientMap := &sync.Map{}
	h := handler.New(clientMap)
	s := service.New(clientMap)

	conn, err := amqp.Dial(*rabbitUrl)
	if err != nil {
		log.Fatalln(err)
	}
	mr, err := messaging.New(conn, s)
	if err != nil {
		log.Fatalln(err)
	}
	go func() { log.Fatal(mr.ListenForMessages()) }()
	tokenManager := jwtauth.New("HS512", []byte(*jwtSecret), nil)
	router := SetupRouter(h, tokenManager)
	log.Printf("listening on port %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, router))
}
