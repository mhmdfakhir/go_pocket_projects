package gordle

import (
	"fmt"
	"os"
	"strings"
)

// hint describes the validity of a character in a word.
type hint byte

const (
	absentCharacter hint = iota
	wrongPosition
	correctPosition
)

// feedback is a list of hints, one per character of the word
type feedback []hint

// Equal determines equality of two feedbacks.
func (fb feedback) Equal(other feedback) bool {
	if len(fb) != len(other) {
		return false
	}
	for index, value := range fb {
		if value != other[index] {
			return false
		}
	}
	return true
}

// String implements the Stringer interface.
func (h hint) String() string {
	switch h {
	case absentCharacter:
		return "⏹️" // grey square
	case wrongPosition:
		return "🟡" // yellow circle
	case correctPosition:
		return "💚" // green heart
	default:
		// This should never happen.
		return "💔" // red broken heart
	}
}

// String implements the Stringer interface for a slice of hints.
func (fb feedback) String() string {
	sb := strings.Builder{}
	for _, h := range fb {
		sb.WriteString(h.String())
	}
	return sb.String()
}

// computeFeedback verifies every character of the guess against
// the solution.
func computeFeedback(guess, solution []rune) feedback {
	// initialise holders for marks
	result := make(feedback, len(guess))
	usedBoolMap := make([]bool, len(solution))

	if len(guess) != len(solution) {
		_, _ = fmt.Fprintf(os.Stderr, "Internal error! Guess and solution"+
			" have different lengths: %d vs %d", len(guess), len(solution))
		return result
	}

	// check for correct letters
	for posInGuess, character := range guess {
		if character == solution[posInGuess] {
			result[posInGuess] = correctPosition
			usedBoolMap[posInGuess] = true
		}
	}

	// look for letters in the wrong position
	for posInGuess, guessChar := range guess {
		if result[posInGuess] != absentCharacter {
			// The character has already been marked, ignore it.
			continue
		}

		for posInSolution, charInSolution := range solution {
			if usedBoolMap[posInSolution] {
				// The letter of the solution is already assigned
				// to a letter of the guess.
				// Skip to the next letter of the solution.
				continue
			}

			if guessChar == charInSolution {
				result[posInGuess] = wrongPosition
				usedBoolMap[posInSolution] = true
				// Skip to the next letter of the guess.
				break
			}
		}
	}

	return result
}
