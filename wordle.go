package main

import (
	"fmt"
	"sort"
	"time"
)

const (
	MaxGuesses = 6
	WordLength = 5
)

var (
	ErrGameOver = fmt.Errorf("game over")
)

type Game struct {
	ID       int64
	Answer   string
	Guesses  []string
	Started  time.Time
	Finished time.Time
	Won      bool
}

type Games []Game

func NewGame(answer string) *Game {
	return &Game{
		Answer:  answer,
		Guesses: make([]string, 0, MaxGuesses),
		Started: time.Now(),
	}
}

// Guess guesses a word and returns true if correct
func (g *Game) Guess(word string) (error, bool) {
	if g.IsDone() {
		return ErrGameOver, false
	}

	if len(word) != WordLength {
		return fmt.Errorf("word must be %d letters", WordLength), false
	}

	if !allowedGuess(word) {
		return fmt.Errorf("invalid word %q", word), false
	}

	g.Guesses = append(g.Guesses, word)
	g.Won = word == g.Answer

	if g.IsDone() {
		g.Finished = time.Now()
		return ErrGameOver, g.Won
	}
	return nil, g.Won
}

func (g *Game) IsDone() bool {
	return g.Won || len(g.Guesses) == MaxGuesses
}

func (g *Game) String() string {
	return g.Render()
}

func (g *Game) Render() string {
	board := ""
	rowTemplate := "[%v][%v][%v][%v][%v]\n"
	emptyrow := "[ ][ ][ ][ ][ ]\n"

	for _, word := range g.Guesses {
		letters := []interface{}{}
		for wordI, c := range word {
			letter := fmt.Sprintf("[%s]", string(c))

			// Color the letter
			for answerI, answerC := range g.Answer {
				if answerC == c {
					// The letter is in the word
					if answerI == wordI {
						// And in the correct spot.
						// Green
						letter = fmt.Sprintf("\033[32;0;0m%s\033[0m", string(c))
					} else {
						// Yellow
						letter = fmt.Sprintf("\033[33;0;0m%s\033[0m", string(c))
					}
				}
			}

			letters = append(letters, letter)
		}
		board += fmt.Sprintf(rowTemplate, letters...)
	}

	for i := len(g.Guesses); i < MaxGuesses; i++ {
		board += emptyrow
	}

	return board
}

func (games Games) Played() int {
	return len(games)
}

func (games Games) WinPercent() int {
	wins := 0
	for _, g := range games {
		if g.Won {
			wins += 1
		}
	}
	return int(float32(wins) / float32(len(games)) * 100)
}

func (games Games) CurrentStreak() int {
	streak := 0
	for i := len(games) - 1; i >= 0; i-- {
		if games[i].Won {
			streak += 1
		} else {
			break
		}
	}
	return streak
}

func (games Games) MaxStreak() int {
	var (
		streaks = []int{0}
		streak  = 0
	)
	for _, game := range games {
		if game.Won {
			streak += 1
		} else {
			streaks = append(streaks, streak)
			streak = 0
		}
	}

	if streak != 0 {
		streaks = append(streaks, streak)
	}

	sort.Ints(streaks)
	return streaks[len(streaks)-1]
}

func (games Games) GuessDistribution() []int {
	dist := []int{0, 0, 0, 0, 0, 0}
	for _, game := range games {
		if !game.Won {
			continue
		}
		// the index of the guess count is zero-based
		i := len(game.Guesses) - 1
		dist[i] = dist[i] + 1
	}

	return dist
}
