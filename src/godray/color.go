package godray

import (
	"image/color"
	"math"
)

// colors
var (
	Red   = &Color{color.RGBA{255, 0, 0, 255}}
	Green = &Color{color.RGBA{0, 255, 0, 255}}
	Blue  = &Color{color.RGBA{0, 0, 255, 255}}
	White = &Color{color.RGBA{255, 255, 255, 255}}
	Black = &Color{color.RGBA{0, 0, 0, 255}}
)

type ColorAlgebra interface {
	Multiply(c color.Color) color.Color
	Scale(factor float64) color.Color
}

type Color struct {
	color.RGBA
}

func multiply(a, b float64) uint8 {
	//fmt.Println("x ", a, b, a*b/255, clamp(0, 255, a*b/255))
	return clamp(0, 255, a*b/255)
}

func clamp(min, max, num float64) uint8 {
	return uint8(math.Min(math.Max(num, min), max))
}

func (c1 Color) Multiply(c2 *Color) *Color {
	r1, g1, b1 := c1.R, c1.G, c1.B
	r2, g2, b2 := c2.R, c2.G, c2.B

	return &Color{
		color.RGBA{
			multiply(float64(r1), float64(r2)),
			multiply(float64(g1), float64(g2)),
			multiply(float64(b1), float64(b2)),
			//multiply(float64(a1), float64(a2)),
			255,
		},
	}
}

func (c1 Color) Add(c2 *Color) *Color {
	r1, g1, b1 := c1.R, c1.G, c1.B
	r2, g2, b2 := c2.R, c2.G, c2.B

	return &Color{
		color.RGBA{
			clamp(0, 255, float64(r1)+float64(r2)),
			clamp(0, 255, float64(g1)+float64(g2)),
			clamp(0, 255, float64(b1)+float64(b2)),
			//uint8(a1 + a2),
			255,
		},
	}
}

func (c1 Color) Scale(scalar float64) *Color {
	r1, g1, b1 := c1.R, c1.G, c1.G

	return &Color{
		color.RGBA{
			clamp(0, 255, float64(r1)*scalar),
			clamp(0, 255, float64(g1)*scalar),
			clamp(0, 255, float64(b1)*scalar),
			255,
		},
	}
}
