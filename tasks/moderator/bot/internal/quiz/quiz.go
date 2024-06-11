package quiz

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/teamteamdev/ugractf-2024-x/tasks/moderator/internal/store"
)

func getId(r *http.Request) (int64, error) {
	value := r.URL.Query().Get("id")
	if value == "" {
		return 0, fmt.Errorf("id is not provided")
	}

	return strconv.ParseInt(value, 10, 64)
}

var (
	internalServerErrorText = []byte(http.StatusText(http.StatusInternalServerError))

	MinScore = flag.Int("min-score", 100, "number of quiz questions to pass")
)

func internalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(internalServerErrorText)
}

func SetupHTTPHandlers(db *store.DB, doneCallback func(id int64)) {
	http.HandleFunc("/question", func(w http.ResponseWriter, r *http.Request) {
		id, err := getId(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var seed int
		row := db.QueryRow(`SELECT seed FROM quizes WHERE id = ?`, id)
		if err := row.Scan(&seed); err == sql.ErrNoRows {
			seed = rand.Int()
			_, err = db.Exec(`INSERT OR REPLACE INTO quizes (id, seed) VALUES (?, ?)`, id, seed)
			if err != nil {
				log.Printf("failed insert new seed: %v", err)
				internalServerError(w)
				return
			}
		} else if err != nil {
			log.Printf("failed select seed: %v", err)
			internalServerError(w)
			return
		}

		q := generateQuestion(seed)

		if err := json.NewEncoder(w).Encode(q); err != nil {
			log.Printf("failed encode question: %v", err)
		}
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		id, err := getId(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		answer, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("failed read body: %v", err)
			internalServerError(w)
			return
		}
		if !utf8.Valid(answer) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid answer"))
			return
		}

		row := db.QueryRow(`SELECT seed FROM quizes WHERE id = ?`, id)
		var seed int
		if err := row.Scan(&seed); err == sql.ErrNoRows {
			w.WriteHeader(http.StatusPreconditionFailed)
			w.Write([]byte("no generated question"))
			return
		} else if err != nil {
			log.Printf("failed read answers: %v", err)
			internalServerError(w)
			return
		}

		q := generateQuestion(seed)
		if !q.IsCorrectAnswer(string(answer)) {
			w.WriteHeader(http.StatusTeapot)
			w.Write([]byte("incorrect answer"))
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Printf("failed create transaction: %v", err)
			internalServerError(w)
			return
		}
		_, err = tx.Exec(`DELETE FROM quizes WHERE id = ? AND seed = ?`, id, seed)
		if err != nil {
			tx.Rollback()
			log.Printf("failed delete quiz: %v", err)
			internalServerError(w)
			return
		}
		_, err = tx.Exec(`UPDATE users SET score = score + 1 WHERE id = ?`, id)
		if err != nil {
			tx.Rollback()
			log.Printf("failed update score: %v", err)
			internalServerError(w)
			return
		}
		var score int
		row = tx.QueryRow(`SELECT score FROM users WHERE id = ?`, id)
		if err := row.Scan(&score); err != nil {
			tx.Rollback()
			log.Printf("failed select score: %v", err)
			internalServerError(w)
			return
		}
		if err := tx.Commit(); err != nil {
			log.Printf("failed commit tx: %v", err)
			internalServerError(w)
			return
		}

		if score >= *MinScore {
			if doneCallback != nil {
				go doneCallback(id)
			}
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("correct answer"))
	})
}
