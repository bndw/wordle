package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gliderlabs/ssh"
)

var cache = map[string]*UserStats{}

type UserStats struct {
	Games []Game
}

func LoadStats(s ssh.Session) *UserStats {
	key := userKey(s)
	user, ok := cache[key]
	if !ok {
		fmt.Printf("new user: %s\n", key)
		u := &UserStats{}
		cache[key] = u
		return u
	}

	fmt.Printf("existing user: %s\n", key)
	return user
}

func SaveStats(s ssh.Session, us *UserStats) {
	key := userKey(s)
	cache[key] = us
}

func userKey(s ssh.Session) string {
	parts := strings.Split(s.RemoteAddr().String(), ":")
	ip := parts[0]
	return fmt.Sprintf("%s|%s", s.User(), ip)
}

func (s *UserStats) Played() int {
	return len(s.Games)
}

func (s *UserStats) WinPercent() int {
	wins := 0
	for _, g := range s.Games {
		if g.Won {
			wins += 1
		}
	}
	return int(float32(wins) / float32(len(s.Games)) * 100)
}

func (s *UserStats) CurrentStreak() int {
	streak := 0
	for i := len(s.Games) - 1; i >= 0; i-- {
		if s.Games[i].Won {
			streak += 1
		} else {
			break
		}
	}
	return streak
}

func (s *UserStats) MaxStreak() int {
	var (
		streaks = []int{0}
		streak  = 0
	)
	for _, game := range s.Games {
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

func (s *UserStats) GuessDistribution() []int {
	dist := []int{0, 0, 0, 0, 0, 0}
	for _, game := range s.Games {
		// the index of the guess count is zero-based
		i := len(game.Guesses) - 1
		dist[i] = dist[i] + 1
	}

	return dist
}
