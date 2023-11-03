package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Input struct {
}

func (i *Input) MoveLeft() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyLeft)
}

func (i *Input) MoveRight() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyRight)
}

func (i *Input) MoveUp() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyUp)
}

func (i *Input) MoveDown() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyDown)
}

func (i *Input) Space() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace)
}
