package main

import (
	"fmt"
	"vlad24/domino/game"
)

func main() {
	var playerName string = "Anonymous Human"
	fmt.Printf("Hello, player! Input your name: ")
	fmt.Scanf("%v", &playerName)
	game.PlayAgainstComputer(playerName)
}
