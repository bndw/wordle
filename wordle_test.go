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
