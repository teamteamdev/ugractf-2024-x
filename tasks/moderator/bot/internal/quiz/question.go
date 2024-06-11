package quiz

import (
	"encoding/json"
	"fmt"
)

type Answer struct {
	Text      string
	IsCorrect bool
}

func (a Answer) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Text)
}

func (a *Answer) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		a.IsCorrect = true
		return json.Unmarshal(data, &a.Text)
	} else {
		parsed := make(map[string]any, 2)
		if err := json.Unmarshal(data, &parsed); err != nil {
			return err
		}

		if len(parsed) != 2 {
			return fmt.Errorf("invalid number of keys in answer: %v", parsed)
		}

		if text, ok := parsed["text"]; ok {
			if a.Text, ok = text.(string); !ok {
				return fmt.Errorf("text field isn't string in answer: %v", parsed)
			}
		} else {
			return fmt.Errorf("text field not in answer: %v", parsed)
		}

		if isCorrect, ok := parsed["is_correct"]; ok {
			if a.IsCorrect, ok = isCorrect.(bool); !ok {
				return fmt.Errorf("is_correct field isn't bool in answer: %v", parsed)
			}
		} else {
			return fmt.Errorf("is_correct field not in answer: %v", parsed)
		}
	}

	return nil
}

type Question struct {
	Text    string   `json:"text"`
	Answers []Answer `json:"answers"`
}

func (q *Question) IsCorrectAnswer(answer string) bool {
	for _, ans := range q.Answers {
		if ans.IsCorrect && answer == ans.Text {
			return true
		}
	}

	return false
}
