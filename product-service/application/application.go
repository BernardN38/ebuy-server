package application

import (
	"database/sql"
	"embed"
	"log"
	"net/http"

	"github.com/BernardN38/ebuy-server/handler"
	"github.com/BernardN38/ebuy-server/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/lib/pq"
)

type Application struct {
	server *Server
}

type Server struct {
	router *chi.Mux
	port   string
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func NewServer(port string, jwtSecret string, h *handler.Handler) *Server {
	tokenManager := jwtauth.New("HS256", []byte(jwtSecret), nil)
	r := SetupRouter(h, tokenManager)
	return &Server{
		router: r,
		port:   port,
	}
}
func New() *Application {
	config, err := getEnvConfig()
	if err != nil {
		log.Fatalln("unable to get env config with error:", err)
		return &Application{}
	}
	//connect to db
	db, err := sql.Open("postgres", config.PostgresDsn)
	if err != nil {
		log.Fatalln(err)
		return &Application{}
	}
	//run db migrations
	err = RunDatabaseMigrations(db)
	if err != nil {
		log.Fatalln(err)
		return &Application{}
	}

	productService := service.New(db)
	handler := handler.New(productService)
	server := NewServer(config.Port, config.JwtSecret, handler)

	return &Application{
		server: server,
	}
}

func (a *Application) Run() {
	//start server
	log.Printf("listening on port %s", a.server.port)
	log.Fatal(http.ListenAndServe(a.server.port, a.server.router))
}
