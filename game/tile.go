package game

import (
	"fmt"
	"math"
)

type Tile struct {
	leftPips  byte
	rightPips byte
}

func (t Tile) id() string {
	minPips := t.leftPips
	maxPips := t.rightPips
	if t.leftPips > t.rightPips {
		minPips, maxPips = maxPips, minPips
	}
	return fmt.Sprintf("%d_%d", minPips, maxPips)
}

func (t Tile) isDoublet() bool {
	return t.leftPips == t.rightPips
}

func (t Tile) hasPips(pips byte) bool {
	return t.leftPips == pips || t.rightPips == pips
}

func (t Tile) pipsOtherThan(pips byte) byte {
	if t.leftPips == pips {
		return t.rightPips
	} else if t.rightPips == pips {
		return t.leftPips
	} else {
		return math.MaxInt8
	}
}

func (t Tile) isSameAs(other Tile) bool {
	return t.hasPips(other.leftPips) && t.hasPips(other.rightPips)
}

func (t Tile) flip() {
	t.leftPips, t.rightPips = t.rightPips, t.leftPips
}

func (t Tile) String() string {
	return fmt.Sprintf("[%d %d]", t.leftPips, t.rightPips)
}
