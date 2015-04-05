package godray

import (
	// "fmt"
	"image/color"
	"math"
)

type Vector struct {
	X, Y, Z float64
}

type Point struct {
	X, Y, Z float64
}

type Ray struct {
	P *Point
	V *Vector
}

func (a *Vector) Add(b *Vector) *Vector {
	return &Vector{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func (a *Point) Add(b *Vector) *Point {
	return &Point{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func (a *Point) Subtract(b *Point) *Vector {
	return &Vector{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func (a *Vector) Subtract(b *Vector) *Vector {
	return a.Add(b.Scale(-1))
}

func (a *Vector) Scale(s float64) *Vector {
	return &Vector{a.X * s, a.Y * s, a.Z * s}
}

func (a *Vector) Dot(b *Vector) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a *Vector) Cross(b *Vector) *Vector {
	return &Vector{
		a.Y*b.Z - a.Z*b.Y,
		a.Z*b.X - a.X*b.Z,
		a.X*b.Y - a.Y*b.X,
	}
}

func (a *Vector) Normalize() *Vector {
	m := a.Magnitude()
	return &Vector{a.X / m, a.Y / m, a.Z / m}
}

func (a *Vector) Magnitude() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
}

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
