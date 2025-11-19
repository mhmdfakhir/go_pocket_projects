package gordle

import (
	"errors"
	"slices"
	"strings"
	"testing"
)

func TestGameAsk(t *testing.T) {
	tt := map[string]struct {
		input string
		want  []rune
	}{
		"5 characters in english": {
			input: "HELLO",
			want:  []rune("HELLO"),
		},
		"5 characters in arabic": {
			input: "مرحبا",
			want:  []rune("مرحبا"),
		},
		"5 characters in japanese": {
			input: "こんにちは",
			want:  []rune("こんにちは"),
		},
		"3 characters in japanese": {
			input: "こんに\nこんにちは",
			want:  []rune("こんにちは"),
		},
		"2 characters in arabic": {
			input: "مرحبا",
			want:  []rune("مرحبا"),
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			g, _ := New(strings.NewReader(tc.input), []string{string(tc.want)}, 0)

			got := g.ask()
			if !slices.Equal(got, tc.want) {
				t.Errorf("got = %v, want %v", string(got), string(tc.want))
			}
		})
	}
}

// Example solution has length 5 here
func TestGameValidateGuess(t *testing.T) {
	tt := map[string]struct {
		word     []rune
		expected error
	}{
		"nominal": {
			word:     []rune("GUESS"),
			expected: nil,
		},
		"too long": {
			word:     []rune("POCKET"),
			expected: errInvalidWordLength,
		},
		"shorter than 5 chars": {
			word:     []rune("POET"),
			expected: errInvalidWordLength,
		},
		"shorter than 5 runes in Arabic": {
			word:     []rune("رررر"),
			expected: errInvalidWordLength,
		},
		"longer than 5 runes in Arabic": {
			word:     []rune("رررررررر"),
			expected: errInvalidWordLength,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			// "three" has len = 5 to reflect a solution of length 5
			// don't need it to be that way
			g, _ := New(nil, []string{"THREE"}, 0)

			err := g.validateGuess(tc.word)
			if !errors.Is(err, tc.expected) {
				t.Errorf("%c, expected %q, got %q",
					tc.word, tc.expected, err)
			}
		})
	}
}

func TestComputeFeedback(t *testing.T) {
	tt := map[string]struct {
		guess            string
		solution         string
		expectedFeedback feedback
	}{
		"nominal": {
			guess:            "HERTZ",
			solution:         "HERTZ",
			expectedFeedback: feedback{correctPosition, correctPosition, correctPosition, correctPosition, correctPosition},
		},
		"double character": {
			guess:            "HELLO",
			solution:         "HELLO",
			expectedFeedback: feedback{correctPosition, correctPosition, correctPosition, correctPosition, correctPosition},
		},
		"double character with wrong answer": {
			guess:            "HELLL",
			solution:         "HELLO",
			expectedFeedback: feedback{correctPosition, correctPosition, correctPosition, correctPosition, absentCharacter},
		},
		"five identical, but only two are there": {
			guess:            "LLLLL",
			solution:         "HELLO",
			expectedFeedback: feedback{absentCharacter, absentCharacter, correctPosition, correctPosition, absentCharacter},
		},
		"two identical, but not in the right position (from left to right)": {
			guess:            "HLLEO",
			solution:         "HELLO",
			expectedFeedback: feedback{correctPosition, wrongPosition, correctPosition, wrongPosition, correctPosition},
		},
		"three identical, but not in the right position (from left to right)": {
			guess:            "HLLLO",
			solution:         "HELLO",
			expectedFeedback: feedback{correctPosition, absentCharacter, correctPosition, correctPosition, correctPosition},
		},
		"one correct, one incorrect, one absent (left of the correct)": {
			guess:            "LLLWW",
			solution:         "HELLO",
			expectedFeedback: feedback{wrongPosition, absentCharacter, correctPosition, absentCharacter, absentCharacter},
		},
		"swapped characters": {
			guess:            "HOLLE",
			solution:         "HELLO",
			expectedFeedback: feedback{correctPosition, wrongPosition, correctPosition, correctPosition, wrongPosition},
		},
		"absent character": {
			guess:            "HULFO",
			solution:         "HELFO",
			expectedFeedback: feedback{correctPosition, absentCharacter, correctPosition, correctPosition, correctPosition},
		},
		"absent character and incorrect": {
			guess:            "HULPP",
			solution:         "HELPO",
			expectedFeedback: feedback{correctPosition, absentCharacter, correctPosition, correctPosition, absentCharacter},
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			fb := computeFeedback([]rune(tc.guess), []rune(tc.solution))
			if !tc.expectedFeedback.Equal(fb) {
				t.Errorf(
					"guess: %q, got the wrong feedback, wanted %v, got %v",
					tc.guess, tc.expectedFeedback, fb)
			}
		})
	}
}

