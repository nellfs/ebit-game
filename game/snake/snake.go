package snake

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/nellfs/ebit-game/common"
	"github.com/nellfs/ebit-game/game"
)

type SnakeBody struct {
	Position     common.Vector2
	NextPosition common.Vector2
	Direction    common.Vector2
	Scale        float32
}

type Snake struct {
	SnakeBody      []SnakeBody
	Input          *game.Input
	movementBuffer []common.Vector2
}

type Fruit struct {
	Position common.Vector2
	Scale    float32
}

func (s *Snake) MoveSnake() {
	for i := len(s.SnakeBody) - 1; i > 0; i-- {
		s.SnakeBody[i].NextPosition = s.SnakeBody[i-1].NextPosition
	}

	if len(s.movementBuffer) > 0 {
		if len(s.movementBuffer) > 1 {
			s.movementBuffer = s.movementBuffer[1:]
		}
		if len(s.movementBuffer) > 4 {
			s.movementBuffer = s.movementBuffer[1:4]
		}
		s.SnakeBody[0].NextPosition.X += s.movementBuffer[0].X
		s.SnakeBody[0].NextPosition.Y += s.movementBuffer[0].Y
	}
}

func (s *Snake) Update() {

	//Snake Input
	if s.Input.MoveRight() {
		s.SnakeBody[0].Direction = common.NewVector2(common.TileSize, 0)
		s.movementBuffer = append(s.movementBuffer, s.SnakeBody[0].Direction)
	} else if s.Input.MoveLeft() {
		s.SnakeBody[0].Direction = common.NewVector2(-common.TileSize, 0)
		s.movementBuffer = append(s.movementBuffer, s.SnakeBody[0].Direction)
	}

	if s.Input.MoveDown() {
		s.SnakeBody[0].Direction = common.NewVector2(0, common.TileSize)
		s.movementBuffer = append(s.movementBuffer, s.SnakeBody[0].Direction)
	} else if s.Input.MoveUp() {
		s.SnakeBody[0].Direction = common.NewVector2(0, -common.TileSize)
		s.movementBuffer = append(s.movementBuffer, s.SnakeBody[0].Direction)
	}

	//Snake Lerp
	for i := 0; i < len(s.SnakeBody); i++ {
		lerpTime := float32(1.0)
		if i == 0 || i == len(s.SnakeBody)-1 {
			lerpTime = 0.3
		}
		if i == len(s.SnakeBody)-1 {
			s.SnakeBody[i].Position.X = common.Lerp(s.SnakeBody[i].Position.X, s.SnakeBody[i-1].Position.X, lerpTime)
			s.SnakeBody[i].Position.Y = common.Lerp(s.SnakeBody[i].Position.Y, s.SnakeBody[i-1].Position.Y, lerpTime)
			continue
		}
		s.SnakeBody[i].Position.X = common.Lerp(s.SnakeBody[i].Position.X, s.SnakeBody[i].NextPosition.X, lerpTime)
		s.SnakeBody[i].Position.Y = common.Lerp(s.SnakeBody[i].Position.Y, s.SnakeBody[i].NextPosition.Y, lerpTime)
	}
}

func (s *Snake) Draw(screen *ebiten.Image) {
	for i := 0; i < len(s.SnakeBody); i++ {
		if i != 0 {
			s.drawLine(screen, common.NewVector2(s.SnakeBody[i].Position.X, s.SnakeBody[i].Position.Y), common.NewVector2(s.SnakeBody[i-1].Position.X, s.SnakeBody[i-1].Position.Y))
		}
	}

}

func (s *Snake) drawLine(screen *ebiten.Image, current common.Vector2, next common.Vector2) {
	var path vector.Path
	path.LineTo(current.X, current.Y)
	path.LineTo(next.X, next.Y)

	options := &vector.StrokeOptions{LineCap: vector.LineCapRound, LineJoin: vector.LineJoinRound, MiterLimit: 0, Width: 32}
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, options)

	screen.DrawTriangles(vs, is, screen, &ebiten.DrawTrianglesOptions{
		AntiAlias: true,
	})
}
