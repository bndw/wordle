package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

const dbFile = "wordle.db"

var (
	IdleTimeout = 5 * time.Minute
	repo        *sqliteRepo
)

func main() {
	hostKey := flag.String("key", "key.pem", "key")
	port := flag.String("port", "22", "port")
	flag.Parse()

	var err error
	repo, err = newRepo(dbFile)
	if err != nil {
		log.Fatal(err)
	}

	server := &ssh.Server{
		Addr:        fmt.Sprintf(":%s", *port),
		IdleTimeout: IdleTimeout,
		Handler:     handler,
	}

	hostKeyPEM, err := os.ReadFile(*hostKey)
	if err != nil {
		log.Fatal(err)
	}
	hostKeySigner, err := gossh.ParsePrivateKey(hostKeyPEM)
	if err != nil {
		log.Fatal(err)
	}
	server.AddHostKey(hostKeySigner)

	fmt.Printf("listening on :%s\n", *port)
	log.Fatal(server.ListenAndServe())
}

func handler(s ssh.Session) {
	var (
		ctx  = s.Context()
		term = terminal.NewTerminal(s, "")
		game = NewGame(wordOfTheDay())
		user = userKey(s)
	)
	term.SetPrompt("> ")

	log.Printf("player connected: %s\n", user)
	defer func() {
		log.Printf("player disconnected: %s\n", user)
	}()

	games, err := repo.ListGames(ctx, user)
	if err != nil {
		log.Printf("failed to list games for user %s: %v", user, err.Error())
		return
	}

	if len(games) > 0 {
		lastGame := games[0]
		if lastGame.Answer == wordOfTheDay() {
			RenderStats(s, term, &lastGame)
			return
		}
	}

	// Render the initial game board
	Render(s, term, game)

	for {
		word, err := term.ReadLine()
		if err != nil {
			log.Printf("read line err: %v", err)
			return
		}

		err, win := game.Guess(word)
		switch {
		case win:
			// Win, game over
			Render(s, term, game)
			time.Sleep(time.Millisecond * 700)
			WarnGreen(s, term, "Winner!\n")
			repo.SaveGame(ctx, user, game)
			RenderStats(s, term, game)
			return
		case err != nil && errors.Is(err, ErrGameOver):
			// Lose, game over
			Render(s, term, game)
			time.Sleep(time.Millisecond * 700)
			Warn(s, term, game.Answer)
			repo.SaveGame(ctx, user, game)
			RenderStats(s, term, game)
			return
		case err != nil:
			// General error, warn and keep going
			Warn(s, term, err.Error())
			fallthrough
		default:
			// Keep going
			Render(s, term, game)
		}
	}
}

func RenderStats(s ssh.Session, term *terminal.Terminal, game *Game) {
	user := userKey(s)
	games, _ := repo.ListGames(s.Context(), user)

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
	time.Sleep(time.Millisecond * 500)
	Clear(s)
}

func WarnGreen(s ssh.Session, term *terminal.Terminal, text string) {
	Clear(s)
	PrintGreen(s, term, text)
	time.Sleep(time.Millisecond * 500)
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
