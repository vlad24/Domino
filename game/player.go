package game

type Player interface {
	name() string
	move(moveNumber int, hand []*Tile, train []*Tile, pileSize int) Move
	selectToPick(pileSize int) int
	onHandGrown(hand []*Tile, newTile *Tile)
	onOpponentMoved(move *Move, train []*Tile)
	onOpponentHandGrown()
	onVictory(looserScore int, playerScore int)
	onLoss(winnerScore int, playerScore int)
	onDraw(drawScore int)
}
