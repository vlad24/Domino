package game

import (
	"fmt"
	"log"
	"strings"
)

type Human struct {
	humanName string
}

func (p Human) name() string {
	return p.humanName
}

func (p Human) move(moveNumber uint16, hand []Tile, train []Tile, pileSize int) Move {
	var left int16 = -1
	var right int16 = -1
	if len(train) > 0 {
		left = int16(train[0].leftPips)
		right = int16(train[len(train)-1].rightPips)
	}
	fmt.Printf("========= Move %d\n", moveNumber)
	fmt.Printf("=== Pile size: %d\n", pileSize)
	fmt.Printf("=== Board:\n%v\n\n", prettyTrainStr(train))
	fmt.Printf("=== Hand: \n%v\n", prettyHandStr(hand, left, right))
	fmt.Printf("=========\n")
	var input int
	var nextTileNumber int
	var appendRight bool
	for {
		fmt.Printf("Pick a tile (a..%c) with asterisk. Use upper case to append right, lower for left:\n", 'a'+len(hand)-1)
		fmt.Scanf("%c\n", &input)
		if input >= 'a' && input <= 'z' {
			nextTileNumber = input - 'a'
		} else if input >= 'A' && input <= 'Z' {
			nextTileNumber = input - 'A'
			appendRight = true
		}
		if nextTileNumber < 0 || nextTileNumber >= len(hand) {
			log.Printf("Wrong tile number %v (%c), must be from [0, %d)\n", nextTileNumber, input, len(hand))
			continue
		}
		break
	}
	return Move{hand[nextTileNumber], appendRight}
}

func (p Human) selectToPick(pileSize int) int {
	var tileNumber int
	for {
		fmt.Printf("No tiles suitable for move. Pick new tile from pile [0..%v]: ", pileSize-1)
		fmt.Scanf("%d", &tileNumber)
		if tileNumber < 0 || tileNumber >= pileSize {
			log.Printf("Wrong tile number, must be within 0 and %v\n", pileSize-1)
			continue
		}
		break
	}
	return tileNumber
}

func (p Human) onHandGrown(hand []Tile, newTile Tile) {
	fmt.Printf("You picked a tile :%v.\nYour new hand is: %v\n", newTile, hand)
}

func (p Human) onOpponentMoved(move Move, train []Tile) {
}

func (p Human) onOpponentHandGrown() {
	log.Printf("Opponent took one more tile")
}

func (p Human) onVictory(looserScore uint16, playerScore uint16) {
}

func (p Human) onLoss(winnerScore uint16, playerScore uint16) {
}

func (p Human) onDraw(drawScore uint16) {
}

func prettyTrainStr(train []Tile) string {
	tileStrings := make([]string, 0, len(train))
	for i := 0; i < len(train); i++ {
		tileStrings = append(tileStrings, fmt.Sprintf("%v", train[i]))
	}
	return strings.Join(tileStrings, "")
}

func prettyHandStr(hand []Tile, allowedLeft int16, allowedRight int16) string {
	tileStrings := make([]string, 0, len(hand))
	unbound := allowedLeft < 0 || allowedRight < 0
	for i := 0; i < len(hand); i++ {
		tile := hand[i]
		tileRef := 'a' + i
		leftMarker := ""
		rightMarker := ""
		if unbound {
			leftMarker = "*"
			rightMarker = "*"
		} else {
			if tile.hasPips(byte(allowedLeft)) {
				leftMarker = "*"
			}
			if tile.hasPips(byte(allowedRight)) {
				rightMarker = "*"
			}
		}
		tileStrings = append(tileStrings, fmt.Sprintf("%v%c%v:%v", leftMarker, tileRef, rightMarker, hand[i]))
	}
	return strings.Join(tileStrings, ", ")
}
