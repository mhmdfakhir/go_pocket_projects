package main

import (
	"bufio"
	"fmt"
	"learngo/gordle/gordle"
	"os"
)

const maxAttempts = 6

func main() {
	corpus, err := gordle.
		ReadCorpus("corpus/english.txt")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to read corpus: %s", err)
		return
	}

	// Create the game.
	g, err := gordle.New(bufio.NewReader(os.Stdin), corpus, maxAttempts)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to start game: %s", err)
		return
	}

	// Run the game! It will end when it's over
	g.Play()
}
