package game

type Player interface {
	name() string
	move(moveNumber uint16, hand []Tile, train []Tile, pileSize int) Move
	selectToPick(pileSize int) int
	onHandGrown(hand []Tile, newTile Tile)
	onOpponentMoved(move Move, train []Tile)
	onOpponentHandGrown()
	onVictory(looserScore uint16, playerScore uint16)
	onLoss(winnerScore uint16, playerScore uint16)
	onDraw(drawScore uint16)
}
