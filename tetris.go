package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

const (
	width      = 10
	height     = 18
	initialX   = 4
	piecesFile = "./pieces.txt"
)

type Vector struct {
	x, y int
}

func (first Vector) plus(second Vector) Vector {
	return Vector{first.x + second.x, first.y + second.y}
}
func (first Vector) equals(second Vector) bool {
	return first.x == second.x && first.y == second.y
}

// A particular rotational instance of a piece.
type PieceInstance []Vector

type Piece struct {
	rotations       []PieceInstance
	currentRotation int
}

func (p *Piece) instance() PieceInstance {
	return p.rotations[p.currentRotation]
}

func (p *Piece) rotate() {
	p.currentRotation = (p.currentRotation + 1) % len(p.rotations)
}

func (p *Piece) unrotate() {
	p.currentRotation = (p.currentRotation - 1) % len(p.rotations)
	if p.currentRotation < 0 {
		p.currentRotation += len(p.rotations)
	}
}

func TetrisPieces() []Piece {
	return []Piece{Piece{[]PieceInstance{[]Vector{Vector{0, 0}, Vector{1, 0}, Vector{0, 1}, Vector{1, 1}}}, 0},
		Piece{[]PieceInstance{[]Vector{Vector{0, 0}, Vector{1, 0}, Vector{1, 1}, Vector{2, 1}},
			[]Vector{Vector{1, 0}, Vector{0, 1}, Vector{1, 1}, Vector{0, 2}},
		}, 0},
		Piece{[]PieceInstance{[]Vector{Vector{1, 0}, Vector{2, 0}, Vector{0, 1}, Vector{1, 1}},
			[]Vector{Vector{0, 0}, Vector{0, 1}, Vector{1, 1}, Vector{1, 2}},
		}, 0},
		Piece{[]PieceInstance{[]Vector{Vector{1, 0}, Vector{0, 1}, Vector{1, 1}, Vector{2, 1}},
			[]Vector{Vector{0, 0}, Vector{0, 1}, Vector{1, 1}, Vector{0, 2}},
			[]Vector{Vector{0, 0}, Vector{1, 0}, Vector{2, 0}, Vector{1, 1}},
			[]Vector{Vector{1, 0}, Vector{0, 1}, Vector{1, 1}, Vector{1, 2}},
		}, 0},
		Piece{[]PieceInstance{[]Vector{Vector{1, 0}, Vector{1, 1}, Vector{1, 2}, Vector{2, 2}},
			[]Vector{Vector{0, 1}, Vector{1, 1}, Vector{2, 1}, Vector{0, 2}},
			[]Vector{Vector{0, 0}, Vector{1, 0}, Vector{1, 1}, Vector{1, 2}},
			[]Vector{Vector{2, 0}, Vector{0, 1}, Vector{1, 1}, Vector{2, 1}},
		}, 0},
		Piece{[]PieceInstance{[]Vector{Vector{1, 0}, Vector{1, 1}, Vector{1, 2}, Vector{0, 2}},
			[]Vector{Vector{0, 1}, Vector{1, 1}, Vector{2, 1}, Vector{0, 0}},
			[]Vector{Vector{1, 0}, Vector{2, 0}, Vector{1, 1}, Vector{1, 2}},
			[]Vector{Vector{0, 1}, Vector{1, 1}, Vector{2, 1}, Vector{2, 2}},
		}, 0},
		Piece{[]PieceInstance{[]Vector{Vector{1, 0}, Vector{1, 1}, Vector{1, 2}, Vector{1, 3}},
			[]Vector{Vector{0, 1}, Vector{1, 1}, Vector{2, 1}, Vector{3, 1}},
		}, 0},
	}
}

// A VectorSet is a Set of Vectors -- the values of the map have the type struct{} so as to not use any space.
type VectorSet map[Vector]struct{}

// None is the element used as a value in a VectorSet to indicate the vector's (key's) presence in the set. It
// is an empty placeholder.
var None struct{} = struct{}{}

func (ps VectorSet) contains(p Vector) bool {
	_, ok := ps[p]
	return ok
}

func (ps VectorSet) add(p Vector) {
	ps[p] = None
}

func (ps VectorSet) delete(p Vector) {
	delete(ps, p)
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
	board           *Board
	nextPiece       *Piece
	pieces          []Piece
	over            bool
	dropDelayMillis int
	ticker          *time.Ticker
}

func NewGame() *Game {
	game := new(Game)
	game.pieces = TetrisPieces()
	game.board = NewBoard()
	game.board.currentPiece = game.GeneratePiece()
	game.board.currentPosition = Vector{initialX, 0}
	game.nextPiece = game.GeneratePiece()
	game.over = false
	// Start off the delay at 3/4 of a second.
	game.dropDelayMillis = 750
	game.startTicker()
	return game
}

func (game *Game) startTicker() {
	game.ticker = time.NewTicker(time.Duration(game.dropDelayMillis) * time.Millisecond)
}

func (game *Game) stopTicker() {
	game.ticker.Stop()
}

type GameEvent int

const (
	MoveLeft GameEvent = iota
	MoveRight
	MoveDown
	Rotate
	QuickDrop
	Quit
	NoEvent // An event that doesn't cause a change to game state; e.g., a window resize.
)

func (game *Game) Start() {
	game.board.Draw()
gameLoop:
	for {
		eventChan := make(chan GameEvent, 1)
		go func() { eventChan <- waitForUserEvent() }()
		var event GameEvent
		select {
		case event = <-eventChan:
		case <-game.ticker.C:
			event = MoveDown
		}
		switch event {
		case MoveLeft:
			game.Move(Left)
		case MoveRight:
			game.Move(Right)
		case MoveDown:
			game.Move(Down)
		case QuickDrop:
			game.QuickDrop()
		case Rotate:
			game.Rotate()
		case Quit:
			break gameLoop
		}
		if game.over {
			break gameLoop
		}
		game.board.Draw()
	}
}

func waitForTick(ticker *time.Ticker) GameEvent {
	<-ticker.C
	return MoveDown
}

func waitForUserEvent() GameEvent {
	switch event := termbox.PollEvent(); event.Type {
	// Movement: arrow keys or vim controls (h, j, k, l)
	// Exit: 'q' or ctrl-c.
	case termbox.EventKey:
		if event.Ch == 0 { // A special key combo was pressed
			switch event.Key {
			case termbox.KeyCtrlC:
				return Quit
			case termbox.KeyArrowLeft:
				return MoveLeft
			case termbox.KeyArrowUp:
				return Rotate
			case termbox.KeyArrowRight:
				return MoveRight
			case termbox.KeyArrowDown:
				return QuickDrop
			}
		} else {
			switch event.Ch {
			case 'q':
				return Quit
			case 'h':
				return MoveLeft
			case 'k':
				return Rotate
			case 'l':
				return MoveRight
			case 'j':
				return QuickDrop
			}
		}
	case termbox.EventError:
		panic(event.Err)
	}
	return NoEvent
}

func (game *Game) GeneratePiece() *Piece {
	return &game.pieces[rand.Intn(len(game.pieces))]
}

func (board *Board) currentPieceInCollision() bool {
	for _, point := range board.currentPiece.instance() {
		attemptedPoint := point.plus(board.currentPosition)
		if attemptedPoint.x < 0 || attemptedPoint.x >= width ||
			attemptedPoint.y < 0 || attemptedPoint.y >= height ||
			board.cells.contains(attemptedPoint) {
			return true
		}
	}
	return false
}

func (board *Board) moveIfPossible(translation Vector) bool {
	position := board.currentPosition
	board.currentPosition = board.currentPosition.plus(translation)
	if board.currentPieceInCollision() {
		board.currentPosition = position
		return false
	}
	return true
}

func (board *Board) mergeCurrentPiece() {
	for _, point := range board.currentPiece.instance() {
		board.cells.add(point.plus(board.currentPosition))
	}
}

// Check whether a horizontal row is complete.
func (board *Board) rowComplete(y int) bool {
	for x := 0; x < width; x++ {
		if !board.cells.contains(Vector{x, y}) {
			return false
		}
	}
	return true
}

// Clear a single row and move every above cell down.
func (board *Board) collapseRow(rowY int) {
	for y := rowY - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			if board.cells.contains(Vector{x, y}) {
				board.cells.add(Vector{x, y + 1})
			} else {
				board.cells.delete(Vector{x, y + 1})
			}
		}
	}
	// Clear the top row completely
	for x := 0; x < width; x++ {
		board.cells.delete(Vector{x, 0})
	}
}

// Clear any complete rows and move the above blocks down.
func (board *Board) clearRows() {
	y := height - 1
	for y >= 0 {
		for board.rowComplete(y) {
			board.collapseRow(y)
		}
		y -= 1
	}
}

// Anchor the current piece to the board, clears lines, and generate a new piece. Sets the 'game over' state
// if the new piece overlaps existing pieces.
func (game *Game) anchor() {
	game.board.mergeCurrentPiece()
	game.board.clearRows()

	game.board.currentPiece = game.nextPiece
	game.board.currentPosition = Vector{initialX, 0}
	game.nextPiece = game.GeneratePiece()

	if game.board.currentPieceInCollision() {
		game.over = true
	}
}

// Attempt to move.
func (game *Game) Move(where Direction) {
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

	// Perform anchoring if we tried to move down but we were unsuccessful.
	if where == Down && !moved {
		game.anchor()
	}
}

// Drop the piece all the way and anchor it.
func (game *Game) QuickDrop() {
	// Move down as far as possible
	for game.board.moveIfPossible(Vector{0, 1}) {
	}
	game.anchor()
}

func (game *Game) Rotate() {
	game.board.currentPiece.rotate()
	if game.board.currentPieceInCollision() {
		game.board.currentPiece.unrotate()
	}
}

func (board *Board) Filled(position Vector) bool {
	if board.cells.contains(position) {
		return true
	}
	if board.currentPiece == nil {
		return false
	}
	for _, point := range board.currentPiece.instance() {
		if point.plus(board.currentPosition).equals(position) {
			return true
		}
	}
	return false
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
