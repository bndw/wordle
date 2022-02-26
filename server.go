package main

import (
	"errors"
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

func newServer(repo *sqliteRepo, hostKey, port string) (*ssh.Server, error) {
	server := &ssh.Server{
		Addr:        fmt.Sprintf(":%s", port),
		IdleTimeout: time.Minute * 5,
		Handler:     newHandler(repo),
	}

	hostKeyPEM, err := os.ReadFile(hostKey)
	if err != nil {
		return nil, err
	}
	hostKeySigner, err := gossh.ParsePrivateKey(hostKeyPEM)
	if err != nil {
		return nil, err
	}
	server.AddHostKey(hostKeySigner)

	return server, nil
}

func newHandler(repo *sqliteRepo) func(ssh.Session) {
	return func(s ssh.Session) {
		var (
			ctx        = s.Context()
			term       = terminal.NewTerminal(s, "")
			user       = userKey(s)
			todaysWord = wordOfTheDay()
			game       = NewGame(todaysWord)
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
			if lastGame.Answer == todaysWord {
				if lastGame.IsDone() {
					// Today's game is already complete
					renderStats(s, term, &lastGame, games)
					return
				} else {
					// Continue the unfinished game
					game = &lastGame
				}
			}
		}

		// Render the initial game board
		render(s, term, game)

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
				render(s, term, game)
				time.Sleep(time.Millisecond * 700)
				warnGreen(s, term, "Winner!\n")
				repo.SaveGame(ctx, user, game)
				renderStats(s, term, game, games)
				return
			case err != nil && errors.Is(err, ErrGameOver):
				// Lose, game over
				render(s, term, game)
				time.Sleep(time.Millisecond * 700)
				warn(s, term, game.Answer)
				repo.SaveGame(ctx, user, game)
				renderStats(s, term, game, games)
				return
			case err != nil:
				// General error, warn and keep going
				warn(s, term, err.Error())
				repo.SaveGame(ctx, user, game)
				fallthrough
			default:
				// Keep going
				repo.SaveGame(ctx, user, game)
				render(s, term, game)
			}
		}
	}
}

func render(s ssh.Session, term *terminal.Terminal, game *Game) {
	clear(s)
	print(s, term, "    Wordle\n")

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
					printGreen(s, term, letterBoxed)
					continue GUESSLETTERS
				}

				if letter == answerLetter {
					letterInWord = true
				}
			}

			if letterInWord {
				printYellow(s, term, letterBoxed)
			} else {
				print(s, term, letterBoxed)
			}
		}
		print(s, term, "\n") // Newline for each word.
	}

	for i := len(game.Guesses); i < MaxGuesses; i++ {
		// Print rows of empty boxes for each remaining guess.
		print(s, term, "[ ][ ][ ][ ][ ]\n")
	}
}

func renderStats(s ssh.Session, term *terminal.Terminal, game *Game, games Games) {
	render(s, term, game)
	print(s, term, "\n    Statistics\n")
	print(s, term, fmt.Sprintf("played..................%d\n", games.Played()))
	print(s, term, fmt.Sprintf("win %%...................%d\n", games.WinPercent()))
	print(s, term, fmt.Sprintf("current streak..........%d\n", games.CurrentStreak()))
	print(s, term, fmt.Sprintf("max streak..............%d\n", games.MaxStreak()))
	print(s, term, "guess distribution.......\n")
	for i, val := range games.GuessDistribution() {
		print(s, term, fmt.Sprintf("    %d...................%d\n", i+1, val))
	}

	var (
		now   = time.Now()
		hours = (24 - now.Hour()) - 1
		mins  = 60 - now.Minute()
	)
	print(s, term, fmt.Sprintf("\nNext Wordle in %d hours %d mins\n", hours, mins))
}

func warn(s ssh.Session, term *terminal.Terminal, text string) {
	clear(s)
	printRed(s, term, text)
	time.Sleep(time.Millisecond * 500)
	clear(s)
}

func warnGreen(s ssh.Session, term *terminal.Terminal, text string) {
	clear(s)
	printGreen(s, term, text)
	time.Sleep(time.Millisecond * 500)
	clear(s)
}

func print(s ssh.Session, term *terminal.Terminal, text string) {
	io.WriteString(s, text)
}

func printRed(s ssh.Session, term *terminal.Terminal, text string) {
	text = fmt.Sprintf("%s%s%s", term.Escape.Red, text, term.Escape.Reset)
	print(s, term, text)
}

func printGreen(s ssh.Session, term *terminal.Terminal, text string) {
	text = fmt.Sprintf("%s%s%s", term.Escape.Green, text, term.Escape.Reset)
	print(s, term, text)
}

func printYellow(s ssh.Session, term *terminal.Terminal, text string) {
	text = fmt.Sprintf("%s%s%s", term.Escape.Yellow, text, term.Escape.Reset)
	print(s, term, text)
}

func clear(s ssh.Session) {
	io.WriteString(s, "\033[H\033[2J")
}
