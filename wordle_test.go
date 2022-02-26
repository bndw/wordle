package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame(t *testing.T) {
	var tests = []struct {
		Name    string
		Answer  string
		Guesses []string
	}{
		{"first guess", "water", []string{"water"}},
		{"second guess", "water", []string{"teeth", "water"}},
		{"third guess", "water", []string{"teeth", "salad", "water"}},
		{"fourth guess", "water", []string{"teeth", "salad", "cheese", "water"}},
		{"fourth guess", "water", []string{"teeth", "salad", "cheese", "grape", "water"}},
		{"fifth guess", "water", []string{"teeth", "salad", "cheese", "grape", "chili", "water"}},
		{"loss", "water", []string{"teeth", "salad", "soupa", "trails", "orange"}},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			game := NewGame(tt.Answer)

			for _, word := range tt.Guesses {
				_, won := game.Guess(word)
				shouldWin := tt.Answer == word
				assert.Equal(t, shouldWin, won)
			}
		})
	}
}

func TestWinPercent(t *testing.T) {
	games := Games{
		{Answer: "water", Won: true},
		{Answer: "waste", Won: true},
		{Answer: "treat", Won: false},
		{Answer: "onion", Won: false},
	}

	assert.Equal(t, 50, games.WinPercent())
}

func TestCurrentStreak(t *testing.T) {
	games := Games{
		{Answer: "water", Won: true},
		{Answer: "waste", Won: false},
		{Answer: "treat", Won: true},
		{Answer: "onion", Won: true},
	}

	assert.Equal(t, 2, games.CurrentStreak())
}

func TestMaxStreak(t *testing.T) {
	var tests = []struct {
		Name      string
		Games     Games
		MaxStreak int
	}{
		{
			Name:      "one game win",
			Games:     Games{{Answer: "water", Won: true}},
			MaxStreak: 1,
		},
		{
			Name:      "one game lose",
			Games:     Games{{Answer: "water", Won: false}},
			MaxStreak: 0,
		},
		{
			Name: "many games",
			Games: Games{
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
			Games: Games{
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
			assert.Equal(t, tt.MaxStreak, tt.Games.MaxStreak())
		})
	}

}
