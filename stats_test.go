package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWinPercent(t *testing.T) {
	games := []Game{
		{Answer: "water", Won: true},
		{Answer: "waste", Won: true},
		{Answer: "treat", Won: false},
		{Answer: "onion", Won: false},
	}

	assert.Equal(t, 50, WinPercent(games))
}

func TestCurrentStreak(t *testing.T) {
	games := []Game{
		{Answer: "water", Won: true},
		{Answer: "waste", Won: false},
		{Answer: "treat", Won: true},
		{Answer: "onion", Won: true},
	}

	assert.Equal(t, 2, CurrentStreak(games))
}

func TestMaxStreak(t *testing.T) {
	var tests = []struct {
		Name      string
		Games     []Game
		MaxStreak int
	}{
		{
			Name:      "one game win",
			Games:     []Game{{Answer: "water", Won: true}},
			MaxStreak: 1,
		},
		{
			Name:      "one game lose",
			Games:     []Game{{Answer: "water", Won: false}},
			MaxStreak: 0,
		},
		{
			Name: "many games",
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
			MaxStreak: 6,
		},
		{
			Name: "many games 2",
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
			MaxStreak: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t, tt.MaxStreak, MaxStreak(tt.Games))
		})
	}

}
