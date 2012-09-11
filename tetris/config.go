package tetris

import (
	"github.com/nsf/termbox-go"
)

const (
	// The width of the game board in game cells (each game cell is two terminal cells wide).
	width = 10
	// The height of the game board.
	height = 18
	// The background color of the game. It's necessary to set this to ensure that the colors work well with any
	// terminal background color.
	backgroundColor = termbox.ColorBlack
)
