package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	// "math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

type Vector2 struct {
	X float32
	Y float32
}

type SnakePart struct {
	Position     Vector2
	NextPosition Vector2
	Direction    Vector2
	Scale        float32
}

type Fruit struct {
	Position Vector2
	Scale    float32
}

type Game struct {
	FramesCouter   int32
	Snake          []SnakePart
	Fruit          Fruit
	Input          *Input
	MovementBuffer []Vector2

	MoveTime int32
	Time     float32

	aa       bool
	counter  int
	vertices []ebiten.Vertex
	indices  []uint16
}

func NewGame() ebiten.Game {
	g := &Game{}
	g.Init()

	return g
}
func NewVector2(x, y float32) Vector2 {
	return Vector2{x, y}
}

func (g *Game) Init() {
	g.FramesCouter = 0

	g.Snake = []SnakePart{{Position: NewVector2(16, 16), NextPosition: NewVector2(128-16, 128-16), Direction: NewVector2(0, 0), Scale: TileSize / 2},
		{Position: NewVector2(16, 16), NextPosition: NewVector2(16, 16), Direction: NewVector2(0, 0), Scale: TileSize / 2},
		{Position: NewVector2(16, 16), NextPosition: NewVector2(16, 16), Direction: NewVector2(0, 0), Scale: TileSize / 2},
		{Position: NewVector2(16, 16), NextPosition: NewVector2(16, 16), Direction: NewVector2(0, 0), Scale: TileSize / 2},
	}
	g.Fruit = Fruit{Position: NewVector2(256, 256), Scale: TileSize}

	g.MoveTime = 8
	whiteImage.Fill(color.White)

}

func (g *Game) MoveSnake() {
	for i := len(g.Snake) - 1; i > 0; i-- {
		g.Snake[i].NextPosition = g.Snake[i-1].NextPosition
	}
}

func Lerp(start, end, t float32) float32 {
	return start + t*(end-start)
}

func (g *Game) Update() error {
	if g.Input.MoveRight() {
		g.Snake[0].Direction = NewVector2(TileSize, 0)
		g.MovementBuffer = append(g.MovementBuffer, g.Snake[0].Direction)
	} else if g.Input.MoveLeft() {
		g.Snake[0].Direction = NewVector2(-TileSize, 0)
		g.MovementBuffer = append(g.MovementBuffer, g.Snake[0].Direction)
	}

	if g.Input.MoveDown() {
		g.Snake[0].Direction = NewVector2(0, TileSize)
		g.MovementBuffer = append(g.MovementBuffer, g.Snake[0].Direction)
	} else if g.Input.MoveUp() {
		g.Snake[0].Direction = NewVector2(0, -TileSize)
		g.MovementBuffer = append(g.MovementBuffer, g.Snake[0].Direction)
	}

	// Time to Move!
	if g.FramesCouter%g.MoveTime == 0 {
		g.MoveSnake()
		if len(g.MovementBuffer) > 0 {
			if len(g.MovementBuffer) > 1 {
				g.MovementBuffer = g.MovementBuffer[1:]
			}
			if len(g.MovementBuffer) > 5 {
				g.MovementBuffer = g.MovementBuffer[1:5]
			}
			g.Snake[0].NextPosition.X += g.MovementBuffer[0].X
			g.Snake[0].NextPosition.Y += g.MovementBuffer[0].Y
		}
		g.FramesCouter = 0
	}

	//Lerp
	for i := 0; i < len(g.Snake); i++ {
		g.Time = 1
		if i == 0 || i == len(g.Snake)-1 {
			g.Time = 0.3
		}
		if i == len(g.Snake)-1 {
			g.Snake[i].Position.X = Lerp(g.Snake[i].Position.X, g.Snake[i-1].Position.X, g.Time)
			g.Snake[i].Position.Y = Lerp(g.Snake[i].Position.Y, g.Snake[i-1].Position.Y, g.Time)
			continue
		}
		g.Snake[i].Position.X = Lerp(g.Snake[i].Position.X, g.Snake[i].NextPosition.X, g.Time)
		g.Snake[i].Position.Y = Lerp(g.Snake[i].Position.Y, g.Snake[i].NextPosition.Y, g.Time)

	}

	if int32(g.Snake[0].NextPosition.X-16) == int32(g.Fruit.Position.X) && int32(g.Snake[0].NextPosition.Y-16) == int32(g.Fruit.Position.Y) {
		g.Fruit.Position.X = float32(rand.Int31n(25)) * TileSize
		g.Fruit.Position.Y = float32(rand.Int31n(19)) * TileSize
		g.Snake = append(g.Snake, SnakePart{NewVector2(g.Snake[len(g.Snake)-1].NextPosition.X, g.Snake[len(g.Snake)-1].NextPosition.Y), NewVector2(g.Snake[len(g.Snake)-1].NextPosition.X, g.Snake[len(g.Snake)-1].NextPosition.Y), NewVector2(0, 0), 16})
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
		if i != 0 {
			// vector.StrokeLine(screen, g.Snake[i].Position.X, g.Snake[i].Position.Y, g.Snake[i-1].Position.X, g.Snake[i-1].Position.Y, 32, color.White, true)
			g.drawLine(screen, NewVector2(g.Snake[i].Position.X, g.Snake[i].Position.Y), NewVector2(g.Snake[i-1].Position.X, g.Snake[i-1].Position.Y))
		}
		// vector.DrawFilledCircle(screen, g.Snake[i].Position.X, g.Snake[i].Position.Y, g.Snake[i].Scale, color.White, true)
	}

	vector.DrawFilledRect(screen, g.Fruit.Position.X, g.Fruit.Position.Y, g.Fruit.Scale, g.Fruit.Scale, color.RGBA{255, 0, 0, 255}, false)

	//test

	//end

	msg := fmt.Sprintf(`FPS: %0.2f, TPS: %0.2f`, ebiten.ActualFPS(), ebiten.ActualTPS())
	ebitenutil.DebugPrint(screen, msg)

}

func (g *Game) drawLine(screen *ebiten.Image, current Vector2, next Vector2) {
	var path vector.Path
	path.LineTo(current.X, current.Y)
	path.LineTo(next.X, next.Y)

	// Draw the main line in white.
	op := &vector.StrokeOptions{LineCap: vector.LineCapRound, LineJoin: vector.LineJoinRound, MiterLimit: 0, Width: 32}
	vs, is := path.AppendVerticesAndIndicesForStroke(g.vertices[:0], g.indices[:0], op)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = 1
		vs[i].ColorG = 1
		vs[i].ColorB = 1
		vs[i].ColorA = 1
	}
	screen.DrawTriangles(vs, is, whiteSubImage, &ebiten.DrawTrianglesOptions{
		AntiAlias: true,
	})

	// Draw the center line in red.
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
