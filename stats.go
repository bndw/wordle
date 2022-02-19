package main

import (
	"sort"
)

func Played(games []Game) int {
	return len(games)
}

func WinPercent(games []Game) int {
	wins := 0
	for _, g := range games {
		if g.Won {
			wins += 1
		}
	}
	return int(float32(wins) / float32(len(games)) * 100)
}

func CurrentStreak(games []Game) int {
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

func MaxStreak(games []Game) int {
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

func GuessDistribution(games []Game) []int {
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
