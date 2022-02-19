package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

const dbFile = "wordle.db"

var (
	IdleTimeout = 60 * time.Second
	repo        *sqliteRepo
)

func main() {
	var err error
	repo, err = newRepo(dbFile)
	if err != nil {
		log.Fatal(err)
	}

	ssh.Handle(handler)
	server := &ssh.Server{
		Addr:        ":22",
		IdleTimeout: IdleTimeout,
	}
	log.Fatal(server.ListenAndServe())
}

func handler(s ssh.Session) {
	ctx := context.Background()
	game := NewGame(wordOfTheDay())

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		log.Printf("failed to make old state: %v", err.Error())
		return
	}
	defer terminal.Restore(0, oldState)

	term := terminal.NewTerminal(s, "")
	term.SetPrompt("> ")

	userID := userKey(s)
	games, err := repo.ListGames(ctx, userID)
	if err != nil {
		log.Printf("failed to list games for user %s: %v", userID, err.Error())
		return
	}

	if len(games) > 0 {
		mostRecentGame := games[len(games)-1]
		if mostRecentGame.Answer == wordOfTheDay() {
			RenderStats(ctx, s, term, &mostRecentGame)
			return
		}
	}

	// Render the initial game board
	Render(s, term, game)

	for {
		word, err := term.ReadLine()
		if err == io.EOF {
			fmt.Printf("got EOF: %v", err)
			return
		}

		err, win := game.Guess(word)
		if err != nil {
			if errors.Is(err, ErrGameOver) {
				// Ran out of guesses, game over
				Render(s, term, game)
				time.Sleep(time.Millisecond * 500)
				Warn(s, term, game.Answer)
				repo.SaveGame(ctx, userID, game)
				RenderStats(ctx, s, term, game)
				return
			}
			Warn(s, term, err.Error())
		}

		Render(s, term, game)

		if win {
			time.Sleep(time.Millisecond * 500)
			WarnGreen(s, term, "Winner!\n")
			repo.SaveGame(ctx, userID, game)
			RenderStats(ctx, s, term, game)
			return
		}
	}
}

func RenderStats(ctx context.Context, s ssh.Session, term *terminal.Terminal, game *Game) {
	userID := userKey(s)
	games, _ := repo.ListGames(ctx, userID)

	Render(s, term, game)
	Print(s, term, "\n    Statistics\n")
	Print(s, term, fmt.Sprintf("played..................%d\n", Played(games)))
	Print(s, term, fmt.Sprintf("win %%...................%d\n", WinPercent(games)))
	Print(s, term, fmt.Sprintf("current streak..........%d\n", CurrentStreak(games)))
	Print(s, term, fmt.Sprintf("max streak..............%d\n", MaxStreak(games)))
	Print(s, term, "guess distribution.......\n")
	for i, val := range GuessDistribution(games) {
		Print(s, term, fmt.Sprintf("    %d...................%d\n", i+1, val))
	}

	var (
		now   = time.Now()
		hours = (24 - now.Hour()) - 1
		mins  = 60 - now.Minute()
	)
	Print(s, term, fmt.Sprintf("\nNext Wordle in %d hours %d mins\n", hours, mins))
}

func Render(s ssh.Session, term *terminal.Terminal, game *Game) {
	Clear(s)
	Print(s, term, "    Wordle\n")

	for _, word := range game.Guesses {

	GUESSLETTERS:
		for guessIndex, c := range word {
			var (
				letter      = string(c)
				letterBoxed = fmt.Sprintf("[%s]", letter)
			)

			letterInWord := false
			for answerIndex, answerLetter := range strings.Split(game.Answer, "") {
				if letter == answerLetter && guessIndex == answerIndex {
					PrintGreen(s, term, letterBoxed)
					continue GUESSLETTERS
				}

				if letter == answerLetter {
					letterInWord = true
				}
			}

			if letterInWord {
				PrintYellow(s, term, letterBoxed)
			} else {
				Print(s, term, letterBoxed)
			}
		}
		Print(s, term, "\n") // Newline for each word.
	}

	for i := len(game.Guesses); i < MaxGuesses; i++ {
		// Print rows of empty boxes for each remaining guess.
		Print(s, term, "[ ][ ][ ][ ][ ]\n")
	}
}

func Warn(s ssh.Session, term *terminal.Terminal, text string) {
	Clear(s)
	PrintRed(s, term, text)
	time.Sleep(time.Second * 2)
	Clear(s)
}

func WarnGreen(s ssh.Session, term *terminal.Terminal, text string) {
	Clear(s)
	PrintGreen(s, term, text)
	time.Sleep(time.Second * 2)
	Clear(s)
}

func Print(s ssh.Session, term *terminal.Terminal, text string) {
	io.WriteString(s, text)
}

func PrintRed(s ssh.Session, term *terminal.Terminal, text string) {
	text = fmt.Sprintf("%s%s%s", term.Escape.Red, text, term.Escape.Reset)
	Print(s, term, text)
}

func PrintGreen(s ssh.Session, term *terminal.Terminal, text string) {
	text = fmt.Sprintf("%s%s%s", term.Escape.Green, text, term.Escape.Reset)
	Print(s, term, text)
}

func PrintYellow(s ssh.Session, term *terminal.Terminal, text string) {
	text = fmt.Sprintf("%s%s%s", term.Escape.Yellow, text, term.Escape.Reset)
	Print(s, term, text)
}

func Clear(s ssh.Session) {
	io.WriteString(s, "\033[H\033[2J")
}
