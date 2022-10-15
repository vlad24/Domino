package game

import (
	"fmt"
)

type Move struct {
	tile        Tile
	appendRight bool
}

func (m Move) String() string {
	if m.appendRight {
		return fmt.Sprintf("<-%v", m.tile)
	} else {
		return fmt.Sprintf("%v->", m.tile)
	}
}
