package tetris

import (
	"github.com/nsf/termbox-go"
)

// A particular rotational instance of a piece.
type PieceInstance []Vector

type Piece struct {
	rotations       []PieceInstance
	currentRotation int
	color           termbox.Attribute
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
	return []Piece{Piece{[]PieceInstance{[]Vector{Vector{0, 0}, Vector{1, 0}, Vector{0, 1}, Vector{1, 1}}},
		0, termbox.ColorYellow},
		Piece{[]PieceInstance{[]Vector{Vector{0, 0}, Vector{1, 0}, Vector{1, 1}, Vector{2, 1}},
			[]Vector{Vector{1, 0}, Vector{0, 1}, Vector{1, 1}, Vector{0, 2}},
		}, 0, termbox.ColorRed},
		Piece{[]PieceInstance{[]Vector{Vector{1, 0}, Vector{2, 0}, Vector{0, 1}, Vector{1, 1}},
			[]Vector{Vector{0, 0}, Vector{0, 1}, Vector{1, 1}, Vector{1, 2}},
		}, 0, termbox.ColorGreen},
		Piece{[]PieceInstance{[]Vector{Vector{1, 0}, Vector{0, 1}, Vector{1, 1}, Vector{2, 1}},
			[]Vector{Vector{0, 0}, Vector{0, 1}, Vector{1, 1}, Vector{0, 2}},
			[]Vector{Vector{0, 0}, Vector{1, 0}, Vector{2, 0}, Vector{1, 1}},
			[]Vector{Vector{1, 0}, Vector{0, 1}, Vector{1, 1}, Vector{1, 2}},
		}, 0, termbox.ColorMagenta},
		Piece{[]PieceInstance{[]Vector{Vector{1, 0}, Vector{1, 1}, Vector{1, 2}, Vector{2, 2}},
			[]Vector{Vector{0, 1}, Vector{1, 1}, Vector{2, 1}, Vector{0, 2}},
			[]Vector{Vector{0, 0}, Vector{1, 0}, Vector{1, 1}, Vector{1, 2}},
			[]Vector{Vector{2, 0}, Vector{0, 1}, Vector{1, 1}, Vector{2, 1}},
		}, 0, termbox.ColorWhite},
		Piece{[]PieceInstance{[]Vector{Vector{1, 0}, Vector{1, 1}, Vector{1, 2}, Vector{0, 2}},
			[]Vector{Vector{0, 1}, Vector{1, 1}, Vector{2, 1}, Vector{0, 0}},
			[]Vector{Vector{1, 0}, Vector{2, 0}, Vector{1, 1}, Vector{1, 2}},
			[]Vector{Vector{0, 1}, Vector{1, 1}, Vector{2, 1}, Vector{2, 2}},
		}, 0, termbox.ColorBlue},
		Piece{[]PieceInstance{[]Vector{Vector{1, 0}, Vector{1, 1}, Vector{1, 2}, Vector{1, 3}},
			[]Vector{Vector{0, 1}, Vector{1, 1}, Vector{2, 1}, Vector{3, 1}},
		}, 0, termbox.ColorCyan},
	}
}
