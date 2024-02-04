package application

import (
	"time"

	"github.com/BernardN38/ebuy-server/product-service/handler"
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

	r.Get("/api/v1/products/health", h.CheckHealth)
	r.Get("/api/v1/products/{productId}", h.GetProduct)
	r.Get("/api/v1/products/recent", h.GetRecentProducts)
	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tm))
		r.Use(jwtauth.Authenticator(tm))
		r.Post("/api/v1/products", h.CreateProduct)
		r.Patch("/api/v1/products/{productId}", h.PatchProduct)
		r.Delete("/api/v1/products/{productId}", h.DeleteProduct)
	})
	return r
}
