package main

import (
	"fmt"
	"github.com/cespare/go-tetris/tetris"
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	game := tetris.NewGame()
	game.Start()

	termbox.Close()
	fmt.Println("Bye!")
}
