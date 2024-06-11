package store

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
	mu sync.Mutex
}

type Tx struct {
	*sql.Tx
	db       *DB
	unlocked bool
}

func (w *DB) Exec(query string, args ...any) (sql.Result, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.DB.Exec(query, args...)
}

func (w *DB) Begin() (*Tx, error) {
	w.mu.Lock()
	tx, err := w.DB.Begin()
	if err != nil {
		w.mu.Unlock()
		return nil, err
	}
	return &Tx{Tx: tx, db: w}, nil
}

func (tx *Tx) Rollback() error {
	err := tx.Tx.Rollback()
	if !tx.unlocked {
		tx.db.mu.Unlock()
		tx.unlocked = true
	}
	return err
}

func (tx *Tx) Commit() error {
	err := tx.Tx.Commit()
	if !tx.unlocked {
		tx.db.mu.Unlock()
		tx.unlocked = true
	}
	return err
}

func NewDB(uri string) (*DB, error) {
	db, err := sql.Open("sqlite3", uri)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			username TEXT NOT NULL,
			score INTEGER NOT NULL DEFAULT 0,
			invited INTEGER NOT NULL DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS quizes (
			id INTEGER PRIMARY KEY,
			seed INTEGER NOT NULL
		);

		CREATE TABLE IF NOT EXISTS banned (
			id INTEGER PRIMARY KEY
		);

		CREATE TABLE IF NOT EXISTS chat_members (
			id INTEGER PRIMARY KEY,
			activity INTEGER NOT NULL DEFAULT 0,
			vote INTEGER DEFAULT NULL
		);

		CREATE TABLE IF NOT EXISTS voting (
			id INTEGER PRIMARY KEY
		);

		DELETE FROM voting;
		UPDATE chat_members SET vote = NULL;
	`)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &DB{db, sync.Mutex{}}, nil
}
