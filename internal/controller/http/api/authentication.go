package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Xacor/gophermart/internal/controller/usecase"
	"github.com/Xacor/gophermart/internal/entity"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type authRoutes struct {
	a usecase.Auth
	l *zap.Logger
}

func newAuthRoutes(handler chi.Router, a usecase.Auth, l *zap.Logger) {
	r := &authRoutes{a, l}

	handler.Post("/register", r.Register)
	handler.Post("/login", r.Authenticate)
}

type tokenResponse struct {
	Token string `json:"token"`
}

func (ar *authRoutes) Register(w http.ResponseWriter, r *http.Request) {
	var (
		user entity.User
		buf  bytes.Buffer
	)

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := ar.a.Register(r.Context(), user)
	if errors.Is(err, usecase.ErrUserExists) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	token, err := ar.a.Authenticate(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resp := tokenResponse{token}

	json, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(json)
}

func (ar *authRoutes) Authenticate(w http.ResponseWriter, r *http.Request) {
	var (
		user entity.User
		buf  bytes.Buffer
	)

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := ar.a.Authenticate(r.Context(), user)
	if errors.Is(err, usecase.ErrInvalidCredentials) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resp := tokenResponse{token}

	json, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(json)
}
