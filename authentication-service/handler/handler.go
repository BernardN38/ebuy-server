package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/BernardN38/ebuy-server/authentication-service/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	authService  *service.AuthService
	validator    *validator.Validate
	tokenManager *TokenManager
}

func New(service *service.AuthService, jwtAuth *jwtauth.JWTAuth) *Handler {
	v := validator.New()
	tm := NewTokenManger(jwtAuth)
	return &Handler{
		authService:  service,
		validator:    v,
		tokenManager: tm,
	}
}

func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}{
		Name:   "auth-service",
		Status: "up",
	})
}
func (h *Handler) Protected(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	w.Write([]byte(fmt.Sprintf("you accessed protected route, claims:%v", claims["user_id"])))
	w.WriteHeader(http.StatusOK)
}
func (h *Handler) CreatUser(w http.ResponseWriter, r *http.Request) {
	var creatUserPayload CreatUserPayload
	err := json.NewDecoder(r.Body).Decode(&creatUserPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.validator.Struct(creatUserPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.authService.CreateUser(r.Context(), service.CreateUserParams{
		Username: creatUserPayload.Username,
		Email:    creatUserPayload.Email,
		Password: creatUserPayload.Password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	userIdint, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "bad user id", http.StatusBadRequest)
		return
	}
	user, err := h.authService.GetUser(r.Context(), userIdint)
	if err != nil {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "error encoding user", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginPayload LoginUserPayload
	err := json.NewDecoder(r.Body).Decode(&loginPayload)
	if err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	userId, err := h.authService.VerifyUser(r.Context(), loginPayload.Email, loginPayload.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not verify user", http.StatusUnauthorized)
		return
	}
	tokenString, err := h.tokenManager.CreateToken(userId)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	cookie := CreateJWTCookie(tokenString)
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}
