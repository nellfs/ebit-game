package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
	TileSize     = 32
)

var (
	whiteImage = ebiten.NewImage(3, 3)

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func Lerp(start, end, t float32) float32 {
	return start + t*(end-start)
}

type Snake struct {
	X, Y, Next_X, Next_Y, Speed_X, Speed_Y, Scale float32
}

type Fruit struct {
	X, Y, Scale float32
}

type Game struct {
	FramesCouter int32
	Snake        []Snake
	Fruit        Fruit
}

func NewGame() ebiten.Game {
	g := &Game{}
	g.Init()

	return g
}

func (g *Game) Init() {
	g.FramesCouter = 0
	g.Snake = []Snake{{X: 128 + 16, Y: 128, Next_X: 128 + 16, Next_Y: 128 + 16, Scale: 16},
		{X: 128 - 16, Y: 128 + 16, Next_X: 0, Next_Y: 0, Scale: 16},
		{X: 128 - 48, Y: 128 + 16, Next_X: 128, Next_Y: 128, Scale: 16},
		{X: 128 - 48, Y: 128 + 16, Next_X: 128, Next_Y: 128, Scale: 16}}
	g.Fruit = Fruit{X: 256, Y: 256, Scale: TileSize}
}

func (g *Game) MoveSnake() {
	for i := len(g.Snake) - 1; i > 0; i-- {
		g.Snake[i].Next_X = g.Snake[i-1].Next_X
		g.Snake[i].Next_Y = g.Snake[i-1].Next_Y
	}
}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.Snake[0].Speed_Y = 0
		g.Snake[0].Speed_X = 32
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.Snake[0].Speed_Y = 0
		g.Snake[0].Speed_X = -32
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.Snake[0].Speed_X = 0
		g.Snake[0].Speed_Y = 32
	} else if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.Snake[0].Speed_X = 0
		g.Snake[0].Speed_Y = -32
	}

	if g.FramesCouter%9 == 0 {
		g.MoveSnake()
		g.Snake[0].Next_X += g.Snake[0].Speed_X
		g.Snake[0].Next_Y += g.Snake[0].Speed_Y
		g.FramesCouter = 0
	}

	for i := 0; i < len(g.Snake); i++ {
		g.Snake[i].X = Lerp(g.Snake[i].X, g.Snake[i].Next_X, 0.18)
		g.Snake[i].Y = Lerp(g.Snake[i].Y, g.Snake[i].Next_Y, 0.18)
	}

	if int32(g.Snake[0].Next_X-16) == int32(g.Fruit.X) && int32(g.Snake[0].Next_Y-16) == int32(g.Fruit.Y) {
		g.Fruit.X = float32(rand.Int31n(25)) * 32
		g.Fruit.Y = float32(rand.Int31n(19)) * 32
		g.Snake = append(g.Snake, Snake{g.Snake[len(g.Snake)-1].X, g.Snake[len(g.Snake)-1].Y, g.Snake[len(g.Snake)-1].Next_X, g.Snake[len(g.Snake)-1].Next_Y, 0, 0, 16})
	}

	g.FramesCouter++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := int32(0); i < ScreenWidth/TileSize; i++ {
		for j := int32(0); j < ScreenWidth/TileSize; j++ {
			if (i+j)%2 == 0 {
				vector.DrawFilledRect(screen, float32(i*TileSize), float32(j*TileSize), TileSize, TileSize, color.RGBA{42, 42, 210, 255}, false)
			} else {
				vector.DrawFilledRect(screen, float32(i*TileSize), float32(j*TileSize), TileSize, TileSize, color.RGBA{32, 40, 145, 255}, false)
			}
		}
	}
	for i := 0; i < len(g.Snake); i++ {
		if i == 0 {
			vector.DrawFilledCircle(screen, g.Snake[i].X, g.Snake[i].Y, g.Snake[i].Scale, color.RGBA{255, 255, 255, 255}, true)
			continue
		}
		vector.DrawFilledCircle(screen, g.Snake[i].X, g.Snake[i].Y, g.Snake[i].Scale, color.RGBA{225 - uint8(i*5), 225 - uint8(i*5), 225 - uint8(i*5), 255}, true)
	}

	vector.DrawFilledRect(screen, g.Fruit.X, g.Fruit.Y, g.Fruit.Scale, g.Fruit.Scale, color.RGBA{255, 0, 0, 255}, false)

	//test

	msg := fmt.Sprintf(`FPS: %0.2f, TPS: %0.2f`, ebiten.ActualFPS(), ebiten.ActualTPS())
	ebitenutil.DebugPrint(screen, msg)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Snake Tetris")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
