package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

const maxPips = 6

func play(players map[byte]Player, state State) {
	for !state.isTerminalForGame() {
		playerNum := state.activePlayer
		opponentNum := opponent(playerNum)
		player := players[playerNum]
		opponent := players[opponentNum]
		log.Printf("%v is moving...", player.name())
		if !state.isTerminalForPlayer(playerNum) {
			move := player.move(state.moveCount+1, state.hands[playerNum], state.train, len(state.pile))
			log.Printf("%v moves: %v", player.name(), &move)
			if err := validate(playerNum, state, move); err != nil {
				log.Println(fmt.Errorf("%v makes invalid move: %v", player.name(), err))
				continue
			}
			state.shrinkHand(playerNum, move.tile)
			opponent.onOpponentMoved(move, state.train)
			if move.appendRight {
				state.growTrainRight(move.tile)
			} else {
				state.growTrainLeft(move.tile)
			}
			state.moveCount++
			state.activePlayer = opponentNum
		} else {
			nextTileNumber := -1
			for nextTileNumber < 0 || nextTileNumber >= len(state.pile) {
				nextTileNumber = player.selectToPick(len(state.pile))
			}
			log.Printf("%v grows hand: %d. Remaining pile: %d", player.name(), nextTileNumber, len(state.pile))
			newTile := state.pile[nextTileNumber]
			state.growHand(playerNum, newTile)
			state.shrinkPile(newTile)
			player.onHandGrown(state.hands[playerNum], newTile)
			players[opponentNum].onOpponentHandGrown()
		}
	}
	log.Printf("Game is over.\nPlayers hands: \n1:%v\n2:%v", state.hands[1], state.hands[2])
	winnerNumber, winnerScore, looserScore := conclude(state)
	if winnerNumber > 0 {
		winner := players[winnerNumber]
		log.Printf("%v won! %d points against %d points", winner.name(), winnerScore, looserScore)
		players[winnerNumber].onVictory(looserScore, winnerScore)
		players[opponent(winnerNumber)].onLoss(winnerScore, looserScore)
	} else {
		log.Printf("Draw! %d points both", winnerScore)
		for _, player := range players {
			player.onDraw(winnerScore)
		}
	}
}

func newGame(player1 Player, player2 Player) State {
	seed := time.Now().UnixNano()
	// seed = 1665853188102626733
	log.Printf("Game seed: %v\n", seed)
	rand.Seed(seed)
	tileIndex := make(map[int]Tile)
	tiles := make([]Tile, 0)
	serialId := 0
	for i := 0; i <= maxPips; i++ {
		for j := i; j <= maxPips; j++ {
			tile := Tile{byte(i), byte(j)}
			tiles = append(tiles, tile)
			tileIndex[serialId] = tile
			serialId++
		}
	}
	tilesAmount := len(tiles)
	rand.Shuffle(tilesAmount, func(i, j int) {
		tiles[i], tiles[j] = tiles[j], tiles[i]
	})
	hands := make(map[byte][]Tile, 2)
	permutation := rand.Perm(tilesAmount)
	const handSize = maxPips + 1
	hands[1] = make([]Tile, 0, handSize)
	for _, id := range permutation[:handSize] {
		hands[1] = append(hands[1], tileIndex[id])
	}
	hands[2] = make([]Tile, 0, handSize)
	for _, id := range permutation[handSize : 2*handSize] {
		hands[2] = append(hands[2], tileIndex[id])
	}
	firstToMove := firstToMove(hands)
	return InitialState(tiles, hands, firstToMove)
}

func opponent(playerNumber byte) byte {
	return 1 + (playerNumber % 2)
}

func firstToMove(hands map[byte][]Tile) byte {
	findSmallestDoublet := func(hand []Tile) int16 {
		var minDoubletEnd int16 = maxPips + 1
		for _, tile := range hand {
			if tile.isDoublet() && int16(tile.leftPips) < minDoubletEnd {
				minDoubletEnd = int16(tile.leftPips)
			}
		}
		if minDoubletEnd > maxPips {
			minDoubletEnd = -1
		}
		return minDoubletEnd
	}
	minDoubletPips1 := findSmallestDoublet(hands[1])
	minDoubletPips2 := findSmallestDoublet(hands[2])
	if minDoubletPips1 == 1 || (minDoubletPips1 >= 0 && minDoubletPips1 < minDoubletPips2) {
		return 1
	} else if minDoubletPips1 == 2 || (minDoubletPips2 >= 0 && minDoubletPips2 < minDoubletPips1) {
		return 2
	} else {
		return 1 + byte(rand.Intn(2))
	}
}

func validate(playerNumber byte, state State, move Move) (err error) {
	if !state.hasTile(playerNumber, move.tile) {
		return fmt.Errorf("player %v doesn't have tile %v", playerNumber, move.tile)
	}
	if state.isInitial() {
		return nil
	}
	if move.appendRight && !move.tile.hasPips(state.rightEnd) {
		return fmt.Errorf("player %d cannot append from right, right end: %v", playerNumber, state.rightEnd)
	}
	if !move.appendRight && !move.tile.hasPips(state.leftEnd) {
		return fmt.Errorf("player %d cannot append from left, left end: %v", playerNumber, state.leftEnd)
	}
	return nil
}

func conclude(state State) (winner byte, winnerScore uint16, looserScore uint16) {
	computeFinalTilesScore := func(tiles []Tile) uint16 {
		if len(tiles) == 1 && tiles[0].isDoublet() && tiles[0].leftPips == 0 {
			return 25
		}
		var sum uint16
		for _, tile := range tiles {
			sum += uint16(tile.leftPips + tile.rightPips)
		}
		return sum
	}
	score1 := computeFinalTilesScore(state.hands[1])
	score2 := computeFinalTilesScore(state.hands[2])
	if score1 < score2 {
		return 1, score1, score2
	} else if score2 < score1 {
		return 2, score2, score1
	} else {
		log.Printf("Draw! Both players have %d points", score1)
		return math.MaxInt8, score1, score2
	}
}

func PlayAgainstComputer(playerName string) {
	human := Human{playerName}
	computer := Computer{}
	players := map[byte]Player{1: human, 2: computer}
	state := newGame(computer, human)
	play(players, state)
}
