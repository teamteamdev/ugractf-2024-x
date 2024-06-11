package voting

import (
	"errors"
	"flag"

	"github.com/teamteamdev/ugractf-2024-x/tasks/moderator/internal/store"
)

var (
	MinActivity   = flag.Int("min-activity", 100, "min activity score to be a candidate")
	MinCandidates = flag.Int("min-candidates", 3, "min number of active candidates to start voting")
	MaxCandidates = flag.Int("max-candidates", 10, "max number candidates for voting epoch")

	ErrTooFewCandidates = errors.New("too few candidate count")
)

func StartVoting(db *store.DB) ([]Candidate, error) {
	rows, err := db.Query(`
		SELECT u.id, name, username
		FROM chat_members c
		JOIN users u ON u.id = c.id
		WHERE activity > ?
		ORDER BY RANDOM()
		LIMIT ?
	`, *MinActivity, *MaxCandidates)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candidates []Candidate

	for rows.Next() {
		var c Candidate
		if err := rows.Scan(&c.ID, &c.Name, &c.Username); err != nil {
			return nil, err
		}
		candidates = append(candidates, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(candidates) < *MinCandidates {
		return candidates, ErrTooFewCandidates
	}

	if err := saveCandidates(db, candidates); err != nil {
		return nil, err
	}

	return candidates, nil
}

func EndVoting(db *store.DB) ([]Candidate, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(`
		SELECT u.id, u.name, u.username, (
				SELECT count(1)
				FROM chat_members c
				WHERE c.vote = v.id
			) as votes
		FROM voting v
		JOIN users u ON u.id = v.id
		ORDER BY votes DESC
	`)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	var candidates []Candidate
	for rows.Next() {
		var c Candidate
		if err := rows.Scan(&c.ID, &c.Name, &c.Username, &c.Votes); err != nil {
			tx.Rollback()
			return nil, err
		}
		candidates = append(candidates, c)
	}
	if err := rows.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}

	if _, err := tx.Exec(`DELETE FROM voting`); err != nil {
		tx.Rollback()
		return nil, err
	}

	if _, err := tx.Exec(`UPDATE chat_members SET vote = NULL`); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return candidates, nil
}
