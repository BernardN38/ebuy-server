package application

import (
	"time"

	"github.com/BernardN38/ebuy-server/authentication-service/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

func SetupRouter(h *handler.Handler, tm *jwtauth.JWTAuth) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/api/v1/auth/health", h.CheckHealth)
	r.Post("/api/v1/auth/users", h.CreatUser)
	r.Get("/api/v1/auth/users/{userId}", h.GetUser)
	r.Post("/api/v1/auth/users/login", h.LoginUser)
	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tm))
		r.Use(jwtauth.Authenticator(tm))
		r.Get("/api/v1/auth/protected", h.Protected)
	})
	return r
}
