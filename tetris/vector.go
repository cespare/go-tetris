package tetris

// A two-dimensional integer-valued vector.
type Vector struct {
	x, y int
}

// Add two vectors.
func (first Vector) plus(second Vector) Vector {
	return Vector{first.x + second.x, first.y + second.y}
}

// Determine whether two vectors are the same.
func (first Vector) equals(second Vector) bool {
	return first.x == second.x && first.y == second.y
}
