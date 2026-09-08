package repo

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

type SaveUser struct {
	Email string
	Hash  string
	ID    string
}

func (s *Storage) UserData(ctx context.Context, res SaveUser) error {
	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO users (email, hash, id) VALUES ( ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, res.Email, res.Hash, res.ID)
	if err != nil {
		return err
	}
	return nil

}

func NewStorage(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, hash TEXT, email TEXT, id TEXT)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(); err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}
