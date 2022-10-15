package game

import "math"

type State struct {
	tiles        []Tile
	pile         []Tile
	train        []Tile
	leftEnd      byte
	rightEnd     byte
	hands        map[byte][]Tile
	activePlayer byte
	moveCount    uint16
}

func InitialState(tiles []Tile, hands map[byte][]Tile, activePlayer byte) State {
	tilesAmount := len(tiles)
	pile := make([]Tile, 0, tilesAmount)
	held := make(map[string]bool)
	for _, t := range append(hands[1], hands[2]...) {
		held[t.id()] = true
	}
	for _, t := range tiles {
		if !held[t.id()] {
			pile = append(pile, t)
		}
	}
	train := make([]Tile, 0, tilesAmount)
	return State{
		tiles:        tiles,
		pile:         pile,
		train:        train,
		leftEnd:      math.MaxInt8,
		rightEnd:     math.MaxInt8,
		hands:        hands,
		activePlayer: activePlayer,
		moveCount:    0,
	}
}

func (s *State) hasTile(playerNumber byte, tile Tile) bool {
	hand := s.hands[playerNumber]
	hasTile := false
	for i := 0; !hasTile && i < len(hand); i++ {
		hasTile = hand[i].isSameAs(tile)
	}
	return hasTile
}

func (s *State) growHand(playerNumber byte, tile Tile) {
	s.hands[playerNumber] = append(s.hands[playerNumber], tile)
}

func (s *State) shrinkHand(playerNumber byte, tile Tile) {
	removeIndex := -1
	hand := s.hands[playerNumber]
	for i := 0; i < len(hand); i++ {
		if hand[i].isSameAs(tile) {
			removeIndex = i
			break
		}
	}
	s.hands[playerNumber] = append(hand[:removeIndex], hand[removeIndex+1:]...)
}

func (s *State) shrinkPile(tile Tile) {
	removeIndex := -1
	for i := 0; i < len(s.pile); i++ {
		if s.pile[i].isSameAs(tile) {
			removeIndex = i
			break
		}
	}
	s.pile = append(s.pile[:removeIndex], s.pile[removeIndex+1:]...)
}

func (s *State) growTrainRight(tile Tile) {
	if s.isInitial() {
		s.leftEnd = tile.leftPips
		s.rightEnd = tile.rightPips
	} else {
		if s.rightEnd != tile.leftPips {
			tile.flip()
		}
		s.rightEnd = tile.rightPips
	}
	s.train = append(s.train, tile)
}

func (s *State) growTrainLeft(tile Tile) {
	if s.isInitial() {
		s.leftEnd = tile.leftPips
		s.rightEnd = tile.rightPips
	} else {
		if s.leftEnd != tile.rightPips {
			tile.flip()
		}
		s.leftEnd = tile.leftPips
	}
	s.train = append([]Tile{tile}, s.train...)
}

func (s *State) isTerminalForPlayer(playerNumber byte) bool {
	if s.isInitial() {
		return false
	}
	hand := s.hands[playerNumber]
	moveExists := false
	for i := 0; !moveExists && i < len(hand); i++ {
		tile := hand[i]
		moveExists = tile.hasPips(s.leftEnd) || tile.hasPips(s.rightEnd)
	}
	return !moveExists
}

func (s *State) isTerminalForGame() bool {
	return !s.isInitial() &&
		s.isTerminalForPlayer(s.activePlayer) &&
		(len(s.hands[s.activePlayer]) == 0 || len(s.pile) == 0)
}

func (s *State) isInitial() bool {
	return len(s.train) == 0
}
