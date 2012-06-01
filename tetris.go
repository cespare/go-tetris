package main

import (
	"fmt"
	"strings"
)

const (
	width  = 10
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
	cells [height][width]bool
	currentPiece *Piece
	currentX, currentY int
}

type Game struct {
	board *Board
	nextPiece *Piece
}

type Piece interface {
	Initialize(*Board) bool
	Move(*Board, direction) bool
	Rotate(*Board, direction) bool
	Filled(x, y int) bool
}

func (board *Board) Filled(x, y int) bool {
	return board.cells[y][x] && board.currentPiece.Filled(x - board.currentX, y - board.currentY)
}

func (board *Board) Draw() {
	fmt.Println("+" + strings.Repeat("-", width) + "+")
	for _, row := range board.cells {
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
