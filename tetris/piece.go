package tetris

import (
	"github.com/nsf/termbox-go"
)

// A particular rotational instance of a piece.
type PieceInstance []Vector

// A Tetris piece, including the all the possible rotations, a color, and an index indicating the current
// rotation of the piece.
type Piece struct {
	rotations       []PieceInstance
	currentRotation int
	initialLocation Vector
	color           termbox.Attribute
}

// Find the current PieceInstance of this piece.
func (p *Piece) instance() PieceInstance {
	return p.rotations[p.currentRotation]
}

// Go to the next rotation.
func (p *Piece) rotate() {
	p.currentRotation = (p.currentRotation + 1) % len(p.rotations)
}

// Go to the previous rotation.
func (p *Piece) unrotate() {
	p.currentRotation = (p.currentRotation - 1) % len(p.rotations)
	if p.currentRotation < 0 {
		p.currentRotation += len(p.rotations)
	}
}

// This has all the hard-coded tetris pieces.
// TODO: It might be nice to have a way to parse these from a configuration file. Maybe the format could look
// something like this:
//
// ##   #
//	## ##
//	   #
//
// ##
// ##
//
// etc.
func tetrisPieces() []Piece {
	// ##
	// ##
	return []Piece{Piece{[]PieceInstance{[]Vector{Vector{0, 0}, Vector{1, 0}, Vector{0, 1}, Vector{1, 1}}},
		0, Vector{4, 0}, termbox.ColorYellow},
		// ##
		//  ##
		Piece{[]PieceInstance{[]Vector{Vector{0, 0}, Vector{1, 0}, Vector{1, 1}, Vector{2, 1}},
			[]Vector{Vector{1, 0}, Vector{0, 1}, Vector{1, 1}, Vector{0, 2}},
		}, 0, Vector{3, 0}, termbox.ColorRed},
		//  ##
		// ##
		Piece{[]PieceInstance{[]Vector{Vector{1, 0}, Vector{2, 0}, Vector{0, 1}, Vector{1, 1}},
			[]Vector{Vector{0, 0}, Vector{0, 1}, Vector{1, 1}, Vector{1, 2}},
		}, 0, Vector{3, 0}, termbox.ColorGreen},
		// ###
		//  #
		Piece{[]PieceInstance{[]Vector{Vector{0, 0}, Vector{1, 0}, Vector{2, 0}, Vector{1, 1}},
			[]Vector{Vector{1, 0}, Vector{0, 1}, Vector{1, 1}, Vector{1, 2}},
			[]Vector{Vector{1, 0}, Vector{0, 1}, Vector{1, 1}, Vector{2, 1}},
			[]Vector{Vector{0, 0}, Vector{0, 1}, Vector{1, 1}, Vector{0, 2}},
		}, 0, Vector{3, 0}, termbox.ColorMagenta},
		// ###
		// #
		Piece{[]PieceInstance{[]Vector{Vector{0, 1}, Vector{1, 1}, Vector{2, 1}, Vector{0, 2}},
			[]Vector{Vector{0, 0}, Vector{1, 0}, Vector{1, 1}, Vector{1, 2}},
			[]Vector{Vector{2, 0}, Vector{0, 1}, Vector{1, 1}, Vector{2, 1}},
			[]Vector{Vector{1, 0}, Vector{1, 1}, Vector{1, 2}, Vector{2, 2}},
		}, 0, Vector{3, -1}, termbox.ColorWhite},
		// ###
		//   #
		Piece{[]PieceInstance{[]Vector{Vector{0, 1}, Vector{1, 1}, Vector{2, 1}, Vector{2, 2}},
			[]Vector{Vector{1, 0}, Vector{1, 1}, Vector{1, 2}, Vector{0, 2}},
			[]Vector{Vector{0, 1}, Vector{1, 1}, Vector{2, 1}, Vector{0, 0}},
			[]Vector{Vector{1, 0}, Vector{2, 0}, Vector{1, 1}, Vector{1, 2}},
		}, 0, Vector{3, -1}, termbox.ColorBlue},
		// ####
		Piece{[]PieceInstance{[]Vector{Vector{0, 1}, Vector{1, 1}, Vector{2, 1}, Vector{3, 1}},
			[]Vector{Vector{1, 0}, Vector{1, 1}, Vector{1, 2}, Vector{1, 3}},
		}, 0, Vector{3, -1}, termbox.ColorCyan},
	}
}
