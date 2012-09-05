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
	initialX   = 4
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
	game.board.currentPosition = Vector{initialX, 0}
	game.nextPiece = game.GeneratePiece()
	return game
}

func (game *Game) Start() {
	game.board.Draw()
gameLoop:
	for {
		gameOver := false
		switch event := termbox.PollEvent(); event.Type {
		// Movement: arrow keys or vim controls (h, j, k, l)
		// Exit: 'q' or ctrl-c.
		case termbox.EventKey:
			if event.Ch == 0 { // A special key combo was pressed
				switch event.Key {
				case termbox.KeyCtrlC:
					break gameLoop
				case termbox.KeyArrowLeft:
					gameOver = game.Move(Left)
				case termbox.KeyArrowUp:
					game.Rotate()
				case termbox.KeyArrowRight:
					gameOver = game.Move(Right)
				case termbox.KeyArrowDown:
					gameOver = game.Move(Down)
				}
			} else {
				switch event.Ch {
				case 'q':
					break gameLoop
				case 'h':
					gameOver = game.Move(Left)
				case 'k':
					game.Rotate()
				case 'l':
					gameOver = game.Move(Right)
				case 'j':
					gameOver = game.Move(Down)
				}
			}
		case termbox.EventError:
			panic(event.Err)
		}
		if gameOver {
			break gameLoop
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

func (board *Board) mergeCurrentPiece() {
	for point, _ := range board.currentPiece.blocks {
		board.cells[point.plus(board.currentPosition)] = None
	}
}

// Anchor the current piece to the board and generate a new piece. Returns whether the new piece overlaps with
// existing pieces (indicating that the game is over).
func (game *Game) anchor() bool {
	game.board.mergeCurrentPiece()
	game.board.currentPiece = game.nextPiece
	game.board.currentPosition = Vector{initialX, 0}
	game.nextPiece = game.GeneratePiece()

	for point, _ := range game.board.currentPiece.blocks {
		if game.board.cells.contains(point.plus(game.board.currentPosition)) {
			return false
		}
	}
	return true
}

// Attempt to move. Returns whether the game ends as a result.
func (game *Game) Move(where Direction) bool {
	translation := Vector{0, 0}
	switch where {
	case Down:
		translation = Vector{0, 1}
	case Left:
		translation = Vector{-1, 0}
	case Right:
		translation = Vector{1, 0}
	}
	// Attempt to make the move.
	moved := game.board.moveIfPossible(translation)

	gameOver := false
	// Perform anchoring if we tried to move down but we were unsuccessful.
	if where == Down && !moved {
		gameOver = !game.anchor()
	}
	return gameOver
}

func (game *Game) Rotate() bool {
	return true
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
