package tetris

import (
	"github.com/nsf/termbox-go"
)

var (
	headerHeight       = 5
	previewHeight      = 6
	sidebarWidth       = 20
	instructionsHeight = 10

	// The internal cells (the board cells) are treated as pairs, so to keep them on even x coordinates we'll
	// put an empty column on the left side.
	totalHeight = headerHeight + height + instructionsHeight + 2
	totalWidth  = (width * 2) + sidebarWidth + 1
)

// A board cell is two terminal cells wide, for squaritude. Only need to set the whole bg color (for filling
// in a cell).
func setBoardCell(x, y int, color termbox.Attribute) {
	termbox.SetCell(x*2+2, headerHeight+y+2, ' ', termbox.ColorDefault, color)
	termbox.SetCell(x*2+3, headerHeight+y+2, ' ', termbox.ColorDefault, color)
}

// Print a message in white text.
func printString(x, y int, message string) {
	for i, ch := range message {
		termbox.SetCell(x+i, y, ch, termbox.ColorWhite, termbox.ColorDefault)
	}
}

// Print a message vertically in white text.
func printStringVertical(x, y int, message string) {
	for i, ch := range message {
		termbox.SetCell(x, y+i, ch, termbox.ColorWhite, termbox.ColorDefault)
	}
}

// Print a box-drawing border character.
func printBorderCharacter(x, y int, ch rune) {
	termbox.SetCell(x, y, ch, termbox.ColorBlue, termbox.ColorDefault)
}

// Print the current score in big ascii art digits
var digitToAsciiArt = map[int][]string{0: []string{" __ ", "/  \\", "\\__/"},
1: []string{"    ", " /| ", "  | "},
2: []string{" __ ", "  _)", " /__"},
3: []string{" __ ", "  _)", " __)"},
4: []string{"    ", "|__|", "   |"},
5: []string{"  __", " |_ ", " __)"},
6: []string{" __ ", "/__ ", "\\__)"},
7: []string{" ___", "   /", "  / "},
8: []string{" __ ", "(__)", "(__)"},
9: []string{" __ ", "(__\\", " __/"},
	}

func drawDigitAsAscii(x, y, digit int) {
	for i, line := range digitToAsciiArt[digit] {
		printString(x, y+i, line)
	}
}

/*
// See http://en.wikipedia.org/wiki/Box-drawing_character for unicode characters.

+---------------------------------------+
|                 header                |
+-----------------------+---------------+
|                       |               |
|                       |   preview     |
|                       |               |
|                       |               |
|        board          +---------------+
|   (width x height)    |               |
|                       |               |
|                       |    score      |
|                       |               |
|                       |               |
+-----------------------+---------------+
|                                       |
|             instructions              |
|                                       |
+---------------------------------------+

*/
func drawStaticBoardParts() {
	// Print the borders.
	for x := 2; x < totalWidth+2; x++ {
		printBorderCharacter(x, 0, '─')
		printBorderCharacter(x, headerHeight+1, '─')
		printBorderCharacter(x, headerHeight+height+2, '─')
		printBorderCharacter(x, totalHeight+1, '─')
	}
	for x := width + 2; x < totalWidth+2; x++ {
		printBorderCharacter(x, headerHeight+previewHeight+2, '─')
	}
	for y := 1; y < totalHeight+1; y++ {
		printBorderCharacter(1, y, '│')
		printBorderCharacter(totalWidth+2, y, '│')
	}
	// Bold borders around the board
	for x := 2; x < (width*2)+2; x++ {
		printBorderCharacter(x, headerHeight+1, '━')
		printBorderCharacter(x, headerHeight+height+2, '━')
	}
	for y := headerHeight + 2; y < headerHeight+height+2; y++ {
		printBorderCharacter(1, y, '┃')
		printBorderCharacter((width*2)+2, y, '┃')
	}
	// Print the various corners
	printBorderCharacter(1, 0, '┌')
	printBorderCharacter(totalWidth+2, 0, '┐')
	printBorderCharacter(totalWidth+2, totalHeight+1, '┘')
	printBorderCharacter(1, totalHeight+1, '└')
	printBorderCharacter(1, headerHeight+1, '┢')
	printBorderCharacter((width*2)+2, headerHeight+1, '┱')
	printBorderCharacter(totalWidth+2, headerHeight+1, '┤')
	printBorderCharacter((width*2)+2, headerHeight+previewHeight+2, '┠')
	printBorderCharacter(totalWidth+2, headerHeight+previewHeight+2, '┤')
	printBorderCharacter(1, headerHeight+height+2, '┡')
	printBorderCharacter((width*2)+2, headerHeight+height+2, '┹')
	printBorderCharacter(totalWidth+2, headerHeight+height+2, '┤')

	// Print the header logo
	header := []string{"",
		"   ____         _____    _        _     ",
		"  / ___| ___   |_   _|__| |_ _ __(_)___ ",
		" | |  _ / _ \\    | |/ _ \\ __| '__| / __|",
		" | |_| | (_) |   | |  __/ |_| |  | \\__ \\",
		"  \\____|\\___/    |_|\\___|\\__|_|  |_|___/",
	}
	for i, line := range header {
		printString(2, i, line)
	}

	// Print the "NEXT" text vertically
	printStringVertical((width*2)+5, headerHeight+3, "NEXT")

	// Print the "SCORE" header
	printString((width*2)+10, headerHeight+previewHeight+4, "SCORE")

	// Print instructions below the game board.
	instructions := []string{"Controls:",
		"",
		"Move left       left arrow or 'h'",
		"Move right      right arrow or 'l'",
		"Move down       down arrow or 'j'",
		"Rotate piece    up arrow or 'k'",
		"Quick drop      space",
		"Quit            ctrl-c or 'q'",
	}
	for i, message := range instructions {
		printString(4, headerHeight+height+4+i, message)
	}
}
