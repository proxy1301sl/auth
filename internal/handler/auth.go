package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"unicode/utf8"

	"github.com/google/uuid"
	authctx "https/github.com/proxy1301sl/auth"
	hash "https/github.com/proxy1301sl/auth/internal"
	"https/github.com/proxy1301sl/auth/repo"
	"https/github.com/proxy1301sl/auth/token"
)

type RequestAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(s *repo.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestAuth
		id := uuid.NewString()
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
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
		res := repo.User{
			Email: req.Email,
			Hash:  encr,
			ID:    id,
			Role:  "",
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

func Login(s *repo.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestAuth
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := ValidateLogin(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := s.GetUser(r.Context(), req.Email)
		if err != nil {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}
		validate := hash.CheckPasswordHash(req.Password, user.Hash)
		if !validate {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}
		token, err := token.GenerateJWT(user.ID, user.Role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := map[string]string{"token": token}
		json.NewEncoder(w).Encode(res)

	}
}

func Profile(s *repo.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := authctx.FromUserId(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		user, err := s.GetUserByID(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp := repo.Response{
			ID:    userID,
			Email: user.Email,
			Role:  user.Role,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)

	}
}
