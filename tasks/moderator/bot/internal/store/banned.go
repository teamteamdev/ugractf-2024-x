package store

import (
	"database/sql"
)

func IsBanned(db *DB, id int64) (bool, error) {
	var one int
	row := db.QueryRow(`SELECT 1 FROM banned WHERE id = ?`, id)
	if err := row.Scan(&one); err == sql.ErrNoRows {
		return false, nil
	} else if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

func Ban(db *DB, id int64) error {
	_, err := db.Exec(`INSERT INTO banned (id) VALUES (?) ON CONFLICT DO NOTHING`, id)
	return err
}
