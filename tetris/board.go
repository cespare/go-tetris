package tetris

import (
	"github.com/nsf/termbox-go"
)

// A map from a point on a board to the color of that cell.
type ColorMap map[Vector]termbox.Attribute

// Returns whether a vector is a member of the color map.
func (cm ColorMap) contains(v Vector) bool {
	_, ok := cm[v]
	return ok
}

// A Board represents the state of a tetris game board, including the current piece that is descending and the
// blocks which already exist on the board.
type Board struct {
	cells           ColorMap
	currentPiece    *Piece
	currentPosition Vector
}

// Create a new empty tetris board with no current piece.
func newBoard() *Board {
	board := new(Board)
	board.cells = make(ColorMap)
	return board
}

// Finds whether the current piece is in collision (going over the edge, or overlapping existing occupied
// blocks). This is useful for testing for collision when moving or rotating by speculatively making the move,
// seeing if it collides, and moving back.
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

// Moves the current piece to another location, if possible. The current piece is updated if this is
// successful; otherwise, the piece is left unmoved. This method returns a boolean indicating whether the move
// was successful.
func (board *Board) moveIfPossible(translation Vector) bool {
	position := board.currentPosition
	board.currentPosition = board.currentPosition.plus(translation)
	if board.currentPieceInCollision() {
		board.currentPosition = position
		return false
	}
	return true
}

// Merge the blocks of the current piece into the game board and remove the current piece.
func (board *Board) mergeCurrentPiece() {
	for _, point := range board.currentPiece.instance() {
		board.cells[point.plus(board.currentPosition)] = board.currentPiece.color
	}
	board.currentPiece = nil
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
			if color, ok := board.cells[Vector{x, y}]; ok {
				board.cells[Vector{x, y + 1}] = color
			} else {
				delete(board.cells, Vector{x, y + 1})
			}
		}
	}
	// Clear the top row completely
	for x := 0; x < width; x++ {
		delete(board.cells, Vector{x, 0})
	}
}

// Clear any complete rows and move the above blocks down. Returns the number of cleared rows.
func (board *Board) clearRows() {
	rowsCleared := 0
	y := height - 1
	for y >= 0 {
		for board.rowComplete(y) {
			rowsCleared += 1
			board.collapseRow(y)
		}
		y -= 1
	}
}

// Find all completed rows.
func (board *Board) clearedRows() []int {
	cleared := make([]int, 0)
	for y := 0; y < height; y++ {
		if board.rowComplete(y) {
			cleared = append(cleared, y)
		}
	}
	return cleared
}

// Finds the color of a particular board cell. It returns the background color if the cell is empty.
func (board *Board) CellColor(position Vector) termbox.Attribute {
	if color, ok := board.cells[position]; ok {
		return color
	}
	if board.currentPiece == nil {
		return backgroundColor
	}
	for _, point := range board.currentPiece.instance() {
		if point.plus(board.currentPosition).equals(position) {
			return board.currentPiece.color
		}
	}
	return backgroundColor
}
