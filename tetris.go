package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"math/rand"
	"time"
)

const (
	width      = 10
	height     = 10
	piecesFile = "./pieces.txt"
)

func readLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

type Direction int

var None struct{} = struct{}{}

type Vector struct {
	x, y int
}

func (first Vector) plus(second Vector) Vector {
	return Vector{first.x + second.x, first.y + second.y}
}
func (first Vector) minus(second Vector) Vector {
	return Vector{first.x - second.x, first.y - second.y}
}

type VectorSet map[Vector]struct{}

type Piece struct {
	blocks VectorSet
}

func (ps VectorSet) contains(p Vector) bool {
	_, ok := ps[p]
	return ok
}

const (
	up Direction = iota + 1
	down
	left
	right
)

type Board struct {
	cells           VectorSet
	currentPiece    *Piece
	currentPosition Vector
}

func NewBoard() *Board {
	board := new(Board)
	board.cells = make(VectorSet)
	return board
}

type Game struct {
	board     *Board
	nextPiece *Piece
	pieces    []Piece
}

func loadPieces() (pieces []Piece, err error) {
	lines, err := readLines(piecesFile)
	waiting := true
	blocks := make(VectorSet)
	pieces = make([]Piece, 0)
	x, y := 0, 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if waiting {
				continue
			}
			pieces = append(pieces, Piece{blocks})
			blocks = make(VectorSet)
			x, y = 0, 0
			waiting = true
		} else {
			waiting = false
			for _, char := range line {
				if char != ' ' {
					blocks[Vector{x, y}] = None
				}
				x++
			}
			y++
			x = 0
		}
	}
	if !waiting {
		pieces = append(pieces, Piece{blocks})
	}
	return
}

func NewPiece(points []Vector) Piece {
	pointSet := make(VectorSet)
	for _, point := range points {
		pointSet[point] = None
	}
	return Piece{pointSet}
}

func NewGame() *Game {
	game := new(Game)
	pieces, _ := loadPieces()
	game.pieces = pieces
	game.board = NewBoard()
	game.board.currentPiece = game.GeneratePiece()
	game.nextPiece = game.GeneratePiece()
	return game
}

func (game *Game) Start() {
	game.board.Draw()
	for {
		game.board.Draw()
	}
}

func (game *Game) GeneratePiece() *Piece {
	return &game.pieces[rand.Intn(len(game.pieces))]
}

func (board *Board) Move(where Direction) bool {
	translation := Vector{0, 0}
	switch where {
	case up:
		translation = Vector{0, -1}
	case down:
		translation = Vector{0, 1}
	case left:
		translation = Vector{-1, 0}
	case right:
		translation = Vector{1, 0}
	}
	board.currentPosition = board.currentPosition.plus(translation)
	return true
}

func (board *Board) Rotate(Direction) bool {
	return true
}

func (board *Board) MergeCurrentPiece() {
	return
}

func (board *Board) Filled(position Vector) bool {
	if board.cells.contains(position) {
		return true
	}
	if board.currentPiece == nil {
		return false
	}
	return board.currentPiece.blocks.contains(position.minus(board.currentPosition))
}

func (board *Board) Draw() {
	fmt.Println("+" + strings.Repeat("-", width) + "+")
	for y := 0; y < height; y++ {
		fmt.Printf("|")
		for x := 0; x < height; x++ {
			if board.Filled(Vector{x, y}) {
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
	rand.Seed(time.Now().UnixNano())
	game := NewGame()
	game.Start()
}
