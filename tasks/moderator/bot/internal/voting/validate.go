package voting

import "errors"

var (
	ErrInvalidVotes = errors.New("invalid votes")
	ErrAlmostEqual  = errors.New("almost equal")
)

func Validate(candidates []Candidate) error {
	if len(candidates) < 2 {
		return nil
	}

	total := 0
	for _, c := range candidates {
		total += c.Votes
	}

	maxDiff := 2 * float64(total) / 10
	minDiff := float64(total) / 10

	if float64(candidates[0].Votes-candidates[1].Votes) < minDiff {
		return ErrAlmostEqual
	}

	for i, c := range candidates[1:] {
		diff := float64(candidates[i].Votes - c.Votes)
		if diff > maxDiff {
			return ErrInvalidVotes
		}
	}

	return nil
}
