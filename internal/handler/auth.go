package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"unicode/utf8"

	"github.com/google/uuid"
	"https/github.com/proxy1301sl/auth/repo"
	"https/github.com/proxy1301sl/auth/token"
	authctx "https/github.com/proxy1301sl/auth/utils"
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
			var MaxErr *http.MaxBytesError
			if errors.As(err, &MaxErr) {
				http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := ValidateLogin(&req); err != nil {
			http.Error(w, "", http.StatusBadRequest)
		}
		user, err := s.GetUser(r.Context(), req.Email)
		if err != nil {
			http.Error(w, "email already exist", http.StatusBadRequest)
			return
		}
		encr, err := authctx.HashedPassword(req.Password)
		if err != nil {
			http.Error(w, "problem with token", http.StatusBadRequest)
			return
		}
		res := repo.User{
			Email: user.Email,
			Hash:  encr,
			ID:    id,
		}
		err = s.UserData(r.Context(), res)
		if err != nil {
			http.Error(w, "invalid email or password", http.StatusInternalServerError)
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
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			var MaxErr *http.MaxBytesError
			if errors.As(err, &MaxErr) {
				http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
				return
			}
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
		validate := authctx.CheckPasswordHash(req.Password, user.Hash)
		if !validate {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}
		tokenString, err := token.GenerateJWT(user.ID, user.Role)
		if err != nil {
			http.Error(w, "failed to create token", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := map[string]string{"token": tokenString}
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
