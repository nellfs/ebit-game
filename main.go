package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/nellfs/ebit-game/common"
	"github.com/nellfs/ebit-game/game/blocks"
	"github.com/nellfs/ebit-game/game/snake"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
	TileSize     = 32
)

type Game struct {
	FramesCouter   int32
	Snake          snake.Snake
	Fruit          snake.Fruit
	Piece          blocks.Piece
	MovementBuffer []common.Vector2

	MoveTime int32
	Time     float32

	vertices []ebiten.Vertex
	indices  []uint16
}

func NewGame() ebiten.Game {
	g := &Game{}
	g.Init()

	return g
}

func (g *Game) Init() {
	g.FramesCouter = 0

	g.Snake.SnakeBody = []snake.SnakeBody{{Position: common.NewVector2(16, 16), NextPosition: common.NewVector2(128-16, 128-16), Direction: common.NewVector2(0, 0), Scale: TileSize / 2},
		{Position: common.NewVector2(16, 16), NextPosition: common.NewVector2(16, 16), Direction: common.NewVector2(0, 0), Scale: TileSize / 2},
		{Position: common.NewVector2(16, 16), NextPosition: common.NewVector2(16, 16), Direction: common.NewVector2(0, 0), Scale: TileSize / 2},
		{Position: common.NewVector2(16, 16), NextPosition: common.NewVector2(16, 16), Direction: common.NewVector2(0, 0), Scale: TileSize / 2},
	}
	g.Fruit = snake.Fruit{Position: common.NewVector2(256, 256), Scale: TileSize}
	g.Piece.Blocks = [][]bool{
		{false, true, false},
		{false, true, false},
		{false, true, false},
		{false, true, false},
		{true, true, true}}

	g.MoveTime = 8
}

func (g *Game) MovePiece() {
	g.Piece.NextPosition.Y += TileSize
}

func (g *Game) Update() error {

	//Movement
	if g.FramesCouter%g.MoveTime == 0 {
		g.MovePiece()
		g.Snake.MoveSnake()
		g.FramesCouter = 0
	}

	g.Snake.Update()

	//Piece
	g.Piece.Position.X = common.Lerp(g.Piece.Position.X, g.Piece.NextPosition.X, g.Time)
	g.Piece.Position.Y = common.Lerp(g.Piece.Position.Y, g.Piece.NextPosition.Y, g.Time)

	g.FramesCouter++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	//Map
	for i := int32(0); i < ScreenWidth/TileSize; i++ {
		for j := int32(0); j < ScreenWidth/TileSize; j++ {
			if (i+j)%2 == 0 {
				vector.DrawFilledRect(screen, float32(i*TileSize), float32(j*TileSize), TileSize, TileSize, color.RGBA{42, 42, 210, 255}, false)
			} else {
				vector.DrawFilledRect(screen, float32(i*TileSize), float32(j*TileSize), TileSize, TileSize, color.RGBA{35, 41, 169, 255}, false)
			}
		}
	}

	g.Snake.Draw(screen)

	//Fruit
	vector.DrawFilledRect(screen, g.Fruit.Position.X, g.Fruit.Position.Y, g.Fruit.Scale, g.Fruit.Scale, color.RGBA{255, 0, 0, 255}, false)

	//Blocks
	for y := 0; y < len(g.Piece.Blocks); y++ {
		for x := 0; x < len(g.Piece.Blocks[y]); x++ {
			if g.Piece.Blocks[y][x] == true {
				vector.DrawFilledRect(screen, float32(x*TileSize)+g.Piece.Position.X, float32(y*TileSize)+g.Piece.Position.Y, g.Fruit.Scale, g.Fruit.Scale, color.RGBA{255, 0, 255, 255}, false)
			}
		}
	}

	//Debug
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
