package game

import (
	"log"
	"math/rand"
)

type Computer struct {
}

func (p Computer) name() string {
	return "Alpha-Dominator"
}

func (p Computer) move(moveNumber int, hand []*Tile, train []*Tile, pileSize int) Move {
	left := -1
	right := -1
	if len(train) != 0 {
		left = train[0].leftPips
		right = train[len(train)-1].rightPips
	}
	var bestTile *Tile = nil
	var appendRight bool = false
	// If at least half of pips are played more than 2 times, we can start analyze pips usage frequences
	frequencyThresh := 3
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
	if bestTile == nil {
		log.Printf("\nTRAIN:%v\nMY HAND: %v.\nEnough occurence data: %v,\n Occurence map: %v", train, hand, enoughOccurenceData, occurenceMap)
	}
	return Move{bestTile, appendRight}
}

func (p Computer) selectToPick(pileSize int) int {
	return rand.Intn(pileSize)
}

func (p Computer) onHandGrown(hand []*Tile, newTile *Tile) {
}

func (p Computer) onOpponentMoved(move *Move, train []*Tile) {
}

func (p Computer) onOpponentHandGrown() {
}

func (p Computer) onVictory(looserScore int, playerScore int) {
}

func (p Computer) onLoss(winnerScore int, playerScore int) {
}

func (p Computer) onDraw(drawScore int) {
}

func (p Computer) bestPipSumTile(hand []*Tile, train []*Tile, left int, right int) (*Tile, bool) {
	matchAny := left < 0 || right < 0
	maxSum := -1
	appendRight := false
	var maxTile *Tile = nil
	for _, v := range hand {
		if matchAny || (v.hasPips(left) || v.hasPips(right)) {
			sum := v.leftPips + v.rightPips
			if sum > maxSum {
				maxSum = sum
				maxTile = v
				appendRight = matchAny || v.hasPips(right)
			}
		}
	}
	return maxTile, appendRight
}

func (p Computer) occurenceMap(hand []*Tile, train []*Tile) map[int]int {
	occurence := make(map[int]int, maxPips+1)
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

func (p Computer) bestOccurenceTile(occurence map[int]int, hand []*Tile, left int, right int) (*Tile, bool) {
	trainScore := occurence[left] + occurence[right]
	maxScore := trainScore
	var bestAttackTile *Tile = nil
	appendRight := false
	for _, tile := range hand {
		if tile.hasPips(left) {
			newScore := trainScore - occurence[left] + occurence[tile.pipsOtherThan(left)]
			if newScore >= maxScore {
				maxScore = newScore
				bestAttackTile = tile
				appendRight = false
			}
		} else if tile.hasPips(right) {
			newScore := trainScore - occurence[right] + occurence[tile.pipsOtherThan(right)]
			if newScore >= maxScore {
				maxScore = newScore
				bestAttackTile = tile
				appendRight = true
			}
		}
	}
	return bestAttackTile, appendRight
}
