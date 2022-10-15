package game

import (
	"math/rand"
)

type Computer struct {
}

func (p Computer) name() string {
	return "Alpha-Dominator"
}

func (p Computer) move(moveNumber uint16, hand []Tile, train []Tile, pileSize int) Move {
	var bestTile *Tile = nil
	var appendRight bool = false
	var left int16 = -1
	var right int16 = -1
	if len(train) > 0 {
		left = int16(train[0].leftPips)
		right = int16(train[len(train)-1].rightPips)
	}
	// If at least half of pips are played more than 2 times, we can start analyze pips usage frequences
	var frequencyThresh byte = 3
	enoughOccurenceData := false
	pipsPlayedOften := 0
	occurenceMap := p.occurenceMap(hand, train)
	for _, pipFreqency := range occurenceMap {
		if pipFreqency >= frequencyThresh {
			pipsPlayedOften++
		}
		enoughOccurenceData = pipsPlayedOften >= maxPips/2
		if enoughOccurenceData {
			break
		}
	}
	if enoughOccurenceData {
		bestTile, appendRight = p.bestOccurenceTile(occurenceMap, hand, left, right)
	}
	if bestTile == nil {
		// If occurence cannot help, then use greedy approach
		bestTile, appendRight = p.bestPipSumTile(hand, train, left, right)
	}
	bestTile, appendRight = p.bestPipSumTile(hand, train, left, right)
	return Move{*bestTile, appendRight}
}

func (p Computer) selectToPick(pileSize int) int {
	return rand.Intn(pileSize)
}

func (p Computer) onHandGrown(hand []Tile, newTile Tile) {
}

func (p Computer) onOpponentMoved(move Move, train []Tile) {
}

func (p Computer) onOpponentHandGrown() {
}

func (p Computer) onVictory(looserScore uint16, playerScore uint16) {
}

func (p Computer) onLoss(winnerScore uint16, playerScore uint16) {
}

func (p Computer) onDraw(drawScore uint16) {
}

func (p Computer) bestPipSumTile(hand []Tile, train []Tile, left int16, right int16) (*Tile, bool) {
	matchAny := left < 0 || right < 0
	var maxSum byte = 0
	appendRight := false
	var maxTile *Tile = nil
	for _, tile := range hand {
		tile := tile
		if matchAny || tile.hasPips(byte(left)) || tile.hasPips(byte(right)) {
			var sum byte = tile.leftPips + tile.rightPips
			if sum >= maxSum {
				maxSum = sum
				maxTile = &tile
				appendRight = matchAny || tile.hasPips(byte(right))
			}
		}
	}
	return maxTile, appendRight
}

func (p Computer) occurenceMap(hand []Tile, train []Tile) map[byte]byte {
	occurence := make(map[byte]byte, maxPips+1)
	for _, v := range append(hand, train...) {
		if v.isDoublet() {
			occurence[v.leftPips]++
		} else {
			occurence[v.leftPips]++
			occurence[v.rightPips]++
		}
	}
	return occurence
}

func (p Computer) bestOccurenceTile(occurence map[byte]byte, hand []Tile, leftEnd int16, rightEnd int16) (*Tile, bool) {
	if leftEnd < 0 || rightEnd < 0 {
		return nil, false
	}
	var left, right byte = byte(leftEnd), byte(rightEnd)
	trainScore := occurence[left] + occurence[right]
	maxScore := trainScore
	var bestAttackTile *Tile = nil
	appendRight := false
	for _, tile := range hand {
		tile := tile
		if tile.hasPips(left) {
			newScore := trainScore - occurence[left] + occurence[tile.pipsOtherThan(left)]
			if newScore >= maxScore {
				maxScore = newScore
				bestAttackTile = &tile
				appendRight = false
			}
		} else if tile.hasPips(right) {
			newScore := trainScore - occurence[right] + occurence[tile.pipsOtherThan(right)]
			if newScore >= maxScore {
				maxScore = newScore
				bestAttackTile = &tile
				appendRight = true
			}
		}
	}
	return bestAttackTile, appendRight
}
