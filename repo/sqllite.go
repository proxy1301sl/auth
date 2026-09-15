package repo

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

type User struct {
	Email string
	Hash  string
	ID    string
	Role  string
}

type Response struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (s *Storage) UserData(ctx context.Context, res User) error {
	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO users (email, hash, id, role) VALUES ( ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, res.Email, res.Hash, res.ID, res.Role)
	if err != nil {
		return err
	}
	return nil

}

func NewStorage(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(5)

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (hash TEXT,  email TEXT UNIQUE, id TEXT, role TEXT NOT NULL DEFAULT 'user' )")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(); err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s *Storage) GetUser(ctx context.Context, email string) (User, error) {
	var u User
	rows := s.db.QueryRowContext(ctx, "SELECT id, email, hash, role FROM users WHERE email = ?", email)
	err := rows.Scan(&u.ID, &u.Email, &u.Hash, &u.Role)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (s *Storage) GetUserByID(ctx context.Context, id string) (User, error) {
	var u User
	rows := s.db.QueryRowContext(ctx, `SELECT id, email, hash, role FROM users WHERE id = ?`, id)
	err := rows.Scan(&u.ID, &u.Email, &u.Hash, &u.Role)
	if err != nil {
		return u, err
	}
	return u, nil

}
