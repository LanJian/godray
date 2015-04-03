package godray

import (
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
	color.Color
}

func multiply(int1, int2 uint32) uint8 {
	return uint8(float64(int1) * float64(int2) / 255)
}

func (c1 Color) Multiply(c2 color.Color) color.Color {
	multiply := func(int1, int2 uint32) uint8 {
		return uint8(float64(int1) * float64(int2) / 255)
	}

	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	return color.RGBA{
		multiply(r1, r2),
		multiply(g1, g2),
		multiply(b1, b2),
		multiply(a1, a2),
	}
}

func (c1 Color) Scale(scalar uint32) color.Color {
	r1, g1, b1, a1 := c1.RGBA()

	return color.RGBA{
		multiply(r1, scalar),
		multiply(g1, scalar),
		multiply(b1, scalar),
		multiply(a1, scalar),
	}
}
