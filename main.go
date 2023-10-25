package main

import (
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

func Lerp(start, end, t float32) float32 {
	return start + t*(end-start)
}

type Snake struct {
	X, Y, Next_X, Next_Y, Scale float32
}

type Fruit struct {
	X, Y, Scale float32
}

type Game struct {
	Snake []Snake
	Fruit Fruit
}

func NewGame() ebiten.Game {
	g := &Game{}
	g.Init()

	return g
}

func (g *Game) Init() {
	g.Snake = []Snake{{X: 128 + 16, Y: 128, Next_X: 128 + 16, Next_Y: 128 + 16, Scale: 16},
		{X: 128 - 16, Y: 128 + 16, Next_X: 0, Next_Y: 0, Scale: 16},
		{X: 128 - 48, Y: 128 + 16, Next_X: 128, Next_Y: 128, Scale: 16},
		{X: 128 - 48, Y: 128 + 16, Next_X: 128, Next_Y: 128, Scale: 16}}
	g.Fruit = Fruit{X: 256, Y: 256, Scale: TileSize}
}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		moveSnake(g)
		g.Snake[0].Next_X += 32
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		moveSnake(g)
		g.Snake[0].Next_X -= 32
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		moveSnake(g)
		g.Snake[0].Next_Y += 32
	} else if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		moveSnake(g)
		g.Snake[0].Next_Y -= 32
	}

	for i := 0; i < len(g.Snake); i++ {
		g.Snake[i].X = Lerp(g.Snake[i].X, g.Snake[i].Next_X, 0.2)
		g.Snake[i].Y = Lerp(g.Snake[i].Y, g.Snake[i].Next_Y, 0.2)
	}

	if int32(g.Snake[0].Next_X-16) == int32(g.Fruit.X) && int32(g.Snake[0].Next_Y-16) == int32(g.Fruit.Y) {
		g.Fruit.X = float32(rand.Int31n(25)) * 32
		g.Fruit.Y = float32(rand.Int31n(19)) * 32
		g.Snake = append(g.Snake, Snake{g.Snake[len(g.Snake)-1].X, g.Snake[len(g.Snake)-1].Y, g.Snake[len(g.Snake)-1].Next_X, g.Snake[len(g.Snake)-1].Next_Y, 16})
	}

	return nil
}

func moveSnake(g *Game) {
	for i := len(g.Snake) - 1; i > 0; i-- {
		g.Snake[i].Next_X = g.Snake[i-1].Next_X
		g.Snake[i].Next_Y = g.Snake[i-1].Next_Y
	}
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
	for i := 0; i < len(g.Snake); i++ {
		vector.DrawFilledCircle(screen, g.Snake[i].X, g.Snake[i].Y, g.Snake[i].Scale, color.RGBA{250 - uint8(i*10), 230 - uint8(i*10), 94, 255}, true)
	}

	vector.DrawFilledRect(screen, g.Fruit.X, g.Fruit.Y, g.Fruit.Scale, g.Fruit.Scale, color.RGBA{255, 0, 0, 255}, false)
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
