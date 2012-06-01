package main

import (
	"fmt"
	"strings"
)

const (
	width  = 20
	height = 10
)

type direction int

const (
	up direction = iota + 1
	down
	left
	right
)

type Board struct {
	Cells [height][width]bool
}

type Piece interface {
	Initialize(*Board) bool
	Move(*Board, direction) bool
	Rotate(*Board, direction) bool
}

func (board *Board) Draw() {
	fmt.Println("+" + strings.Repeat("-", width) + "+")
	for _, row := range board.Cells {
		fmt.Printf("|")
		for _, cell := range row {
			if cell {
				fmt.Printf("#")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Println("|")
	}
	fmt.Println("+" + strings.Repeat("-", width) + "+")
}

func main() {
	board := new(Board)
	board.Draw()
}
