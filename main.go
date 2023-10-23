package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
	TileSize     = 32
)

type Game struct {
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, ":)")
	for i := int32(0); i < ScreenWidth/TileSize; i++ {
		for j := int32(0); j < ScreenWidth/TileSize; j++ {
			if (i+j)%2 == 0 {
				vector.DrawFilledRect(screen, float32(i*TileSize), float32(j*TileSize), TileSize, TileSize, color.RGBA{0, 0, 255, 255}, false)
			} else {
				vector.DrawFilledRect(screen, float32(i*TileSize), float32(j*TileSize), TileSize, TileSize, color.RGBA{0, 0, 100, 255}, false)
			}
		}
	}
	vector.DrawFilledCircle(screen, 250, 250, 50, color.White, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Snake Tetris")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
