package voting

import (
	"fmt"

	"github.com/teamteamdev/ugractf-2024-x/tasks/moderator/internal/store"
)

type Candidate struct {
	ID       int64
	Name     string
	Username string
	Votes    int
}

func (c Candidate) Header() string {
	if c.Username != "" {
		return fmt.Sprintf("%s (%s)", c.Name, c.Username)
	} else {
		return c.Name
	}
}

func saveCandidates(db *store.DB, candidates []Candidate) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE chat_members SET vote = NULL`)
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO voting (id) VALUES (?)`)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, candidate := range candidates {
		if _, err := stmt.Exec(candidate.ID); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
