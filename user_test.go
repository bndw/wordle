package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWinPercent(t *testing.T) {
	stats := UserStats{
		Games: []Game{
			{Answer: "water", Won: true},
			{Answer: "waste", Won: true},
			{Answer: "treat", Won: false},
			{Answer: "onion", Won: false},
		},
	}

	assert.Equal(t, 50, stats.WinPercent())
}

func TestCurrentStreak(t *testing.T) {
	stats := UserStats{
		Games: []Game{
			{Answer: "water", Won: true},
			{Answer: "waste", Won: false},
			{Answer: "treat", Won: true},
			{Answer: "onion", Won: true},
		},
	}

	assert.Equal(t, 2, stats.CurrentStreak())
}

func TestMaxStreak(t *testing.T) {
	var tests = []struct {
		Name      string
		Stats     UserStats
		MaxStreak int
	}{
		{
			Name:      "one game win",
			Stats:     UserStats{Games: []Game{{Answer: "water", Won: true}}},
			MaxStreak: 1,
		},
		{
			Name:      "one game lose",
			Stats:     UserStats{Games: []Game{{Answer: "water", Won: false}}},
			MaxStreak: 0,
		},
		{
			Name: "many games",
			Stats: UserStats{
				Games: []Game{
					{Answer: "water", Won: true},
					{Answer: "waste", Won: false},
					{Answer: "treat", Won: true},
					{Answer: "onion", Won: true},
					{Answer: "treat", Won: true},
					{Answer: "treat", Won: true},
					{Answer: "onion", Won: true},
					{Answer: "onion", Won: true},
					{Answer: "waste", Won: false},
					{Answer: "treat", Won: true},
					{Answer: "onion", Won: true},
					{Answer: "onion", Won: true},
				},
			},
			MaxStreak: 6,
		},
		{
			Name: "many games 2",
			Stats: UserStats{
				Games: []Game{
					{Answer: "water", Won: true},
					{Answer: "waste", Won: false},
					{Answer: "treat", Won: true},
					{Answer: "onion", Won: true},
					{Answer: "treat", Won: true},
					{Answer: "treat", Won: true},
					{Answer: "onion", Won: true},
					{Answer: "waste", Won: false},
					{Answer: "treat", Won: true},
					{Answer: "onion", Won: true},
					{Answer: "onion", Won: true},
					{Answer: "treat", Won: true},
					{Answer: "treat", Won: true},
					{Answer: "onion", Won: true},
					{Answer: "onion", Won: true},
				},
			},
			MaxStreak: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t, tt.MaxStreak, tt.Stats.MaxStreak())
		})
	}

}
