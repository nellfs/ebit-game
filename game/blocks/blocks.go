package blocks

import "github.com/nellfs/ebit-game/common"

type Piece struct {
	Position     common.Vector2
	NextPosition common.Vector2
	Blocks       [][]bool
}
