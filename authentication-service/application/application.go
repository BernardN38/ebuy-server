package application

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/BernardN38/ebuy-server/authentication-service/handler"
	"github.com/BernardN38/ebuy-server/authentication-service/service"
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

func NewServer(port string, jwtSecret string, h *handler.Handler, tokenManager *jwtauth.JWTAuth) *Server {
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

	// Connect to the database
	db, err := sql.Open("postgres", config.PostgresDsn)
	if err != nil {
		log.Fatalln("unable to connect to the database:", err)
		return &Application{}
	}
	defer db.Close()
	// Check if the database exists
	if err := createDatabaseIfNotExists(db, config.DbName); err != nil {
		log.Fatalln("unable to create or check the database:", err)
		return &Application{}
	}

	// Connect to the specific database
	db, err = sql.Open("postgres", config.PostgresDsn+" dbname="+config.DbName)
	if err != nil {
		log.Fatalln("unable to connect to the specific database:", err)
		return &Application{}
	}
	// Run database migrations
	if err := RunDatabaseMigrations(db); err != nil {
		log.Fatalln("unable to run database migrations:", err)
		return &Application{}
	}

	authService := service.New(db)

	jwtAuth := jwtauth.New("HS512", []byte(config.JwtSecret), nil)
	handler := handler.New(authService, jwtAuth)
	server := NewServer(config.Port, config.JwtSecret, handler, jwtAuth)

	return &Application{
		server: server,
	}
}

func (a *Application) Run() {
	// Start server
	log.Printf("listening on port %s", a.server.port)
	log.Fatal(http.ListenAndServe(a.server.port, a.server.router))
}

func createDatabaseIfNotExists(db *sql.DB, dbName string) error {
	result, err := db.Exec(fmt.Sprintf("select 1 from pg_database where datname = '%s'", dbName))
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return err
		}
	}

	return err
}
