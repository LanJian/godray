package godray

import (
	// "fmt"
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
