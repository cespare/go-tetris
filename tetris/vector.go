package tetris

type Vector struct {
	x, y int
}

func (first Vector) plus(second Vector) Vector {
	return Vector{first.x + second.x, first.y + second.y}
}
func (first Vector) equals(second Vector) bool {
	return first.x == second.x && first.y == second.y
}
