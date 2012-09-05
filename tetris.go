package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/nsf/termbox-go"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	width      = 10
	height     = 15
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

type Vector struct {
	x, y int
}

func (first Vector) plus(second Vector) Vector {
	return Vector{first.x + second.x, first.y + second.y}
}
func (first Vector) minus(second Vector) Vector {
	return Vector{first.x - second.x, first.y - second.y}
}

// A VectorSet is a Set of Vectors -- the values of the map have the type struct{} so as to not use any space.
type VectorSet map[Vector]struct{}

// None is the element used as a value in a VectorSet to indicate the vector's (key's) presence in the set. It
// is an empty placeholder.
var None struct{} = struct{}{}

type Piece struct {
	blocks VectorSet
}

func (ps VectorSet) contains(p Vector) bool {
	_, ok := ps[p]
	return ok
}

type Direction int

const (
	Up Direction = iota + 1
	Down
	Left
	Right
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

// Load in the visual representations of pieces from the file containing the piece shapes and emit an array of
// all possible pieces.
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
gameLoop:
	for {
		switch event := termbox.PollEvent(); event.Type {
		case termbox.EventKey:
			switch event.Key {
			case termbox.KeyCtrlC:
				break gameLoop
			case termbox.KeyArrowLeft:
				game.board.Move(Left)
			case termbox.KeyArrowUp:
				game.board.Rotate()
			case termbox.KeyArrowRight:
				game.board.Move(Right)
			case termbox.KeyArrowDown:
				game.board.Move(Down)
			}
		case termbox.EventError:
			panic(event.Err)
		}
		game.board.Draw()
	}
}

func (game *Game) GeneratePiece() *Piece {
	return &game.pieces[rand.Intn(len(game.pieces))]
}

func (board *Board) moveIfPossible(translation Vector) bool {
	attemptedPosition := board.currentPosition.plus(translation)
	for point, _ := range board.currentPiece.blocks {
		attemptedPoint := point.plus(attemptedPosition)
		if attemptedPoint.x < 0 || attemptedPoint.x >= width ||
			attemptedPoint.y < 0 || attemptedPoint.y >= height ||
			board.cells.contains(attemptedPoint) {
			return false
		}
	}
	board.currentPosition = attemptedPosition
	return true
}

func (board *Board) Move(where Direction) bool {
	translation := Vector{0, 0}
	switch where {
	case Down:
		translation = Vector{0, 1}
	case Left:
		translation = Vector{-1, 0}
	case Right:
		translation = Vector{1, 0}
	}
	return board.moveIfPossible(translation)
}

func (board *Board) Rotate() bool {
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
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Print the borders. The internal cells (the board cells) are treated as pairs, so to keep them on even x
	// coordinates we'll put an empty column on the left side.
	termbox.SetCell(1, 0, 0x256D, termbox.ColorBlue, termbox.ColorDefault)
	termbox.SetCell(width*2+2, 0, 0x256E, termbox.ColorBlue, termbox.ColorDefault)
	termbox.SetCell(1, height+1, 0x2570, termbox.ColorBlue, termbox.ColorDefault)
	termbox.SetCell(width*2+2, height+1, 0x256F, termbox.ColorBlue, termbox.ColorDefault)
	for x := 2; x <= width*2+1; x++ {
		termbox.SetCell(x, 0, 0x2500, termbox.ColorBlue, termbox.ColorDefault)
		termbox.SetCell(x, height+1, 0x2500, termbox.ColorBlue, termbox.ColorDefault)
	}
	for y := 1; y <= height; y++ {
		termbox.SetCell(1, y, 0x2502, termbox.ColorBlue, termbox.ColorDefault)
		termbox.SetCell(width*2+2, y, 0x2502, termbox.ColorBlue, termbox.ColorDefault)
	}

	// Print the board contents. Each block will correspond to a side-by-side pair of cells in the termbox, so
	// that the visible blocks will be roughly square.
	for x := 1; x <= width; x++ {
		for y := 1; y <= height; y++ {
			if board.Filled(Vector{x - 1, y - 1}) {
				termbox.SetCell(x*2, y, ' ', termbox.ColorDefault, termbox.ColorGreen)
				termbox.SetCell(x*2+1, y, ' ', termbox.ColorDefault, termbox.ColorGreen)
			}
		}
	}

	termbox.Flush()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	game := NewGame()
	game.Start()

	termbox.Close()
	fmt.Println("Bye!")
}
