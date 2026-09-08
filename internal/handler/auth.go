package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"unicode/utf8"

	"github.com/google/uuid"
	hash "https/github.com/proxy1301sl/auth/internal"
	"https/github.com/proxy1301sl/auth/repo"
)

type RequestAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(s *repo.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestAuth
		id := uuid.NewString()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := ValidateLogin(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		encr, err := hash.HashedPassword(req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res := repo.SaveUser{
			Email: req.Email,
			Hash:  encr,
			ID:    id,
		}
		err = s.UserData(r.Context(), res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func ValidateLogin(req *RequestAuth) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if utf8.RuneCountInString(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return errors.New("invalid email")
	}
	return nil
}
