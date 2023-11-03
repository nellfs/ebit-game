package common

const TileSize = 32

func Lerp(start, end, t float32) float32 {
	return start + t*(end-start)
}

type Vector2 struct {
	X, Y float32
}

func NewVector2(x, y float32) Vector2 {
	return Vector2{x, y}
}
