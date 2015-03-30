package godray

import "math"

type Shape interface {
	Intersect(*Ray) (*Point, float64)
}

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

type Sphere struct {
	Center *Point
	Radius float64
}

func (s *Sphere) Normal(p *Point) *Vector {
	return p.subtract(s.Center).normalize()
}

func (s *Sphere) Intersect(r *Ray) (*Point, float64, *Vector) {
	l := r.V.normalize()
	o := r.P

	leftSide := l.dot(o.subtract(s.Center))

	diff := r.P.subtract(s.Center)
	sqrtSum := (leftSide * leftSide) - diff.dot(diff) + s.Radius*s.Radius

	if sqrtSum < 0 {
		return nil, 0, nil
	} else if sqrtSum == 0.0 {
		d := -1 * leftSide
		intersectPoint := o.add(r.V.Scale(d))
		return intersectPoint, d, s.Normal(intersectPoint)
	} else {
		d1 := -1*leftSide + math.Sqrt(sqrtSum)
		d2 := -1*leftSide - math.Sqrt(sqrtSum)

		pt1 := o.add(r.V.normalize().Scale(d1))
		pt2 := o.add(r.V.normalize().Scale(d2))

		if r.P.subtract(pt1).magnitude()-r.P.subtract(pt2).magnitude() >= 0 {
			return pt2, d2, s.Normal(pt2)
		} else {
			return pt1, d1, s.Normal(pt1)
		}
	}
}

func (a *Vector) add(b *Vector) *Vector {
	return &Vector{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func (a *Point) add(b *Vector) *Point {
	return &Point{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func (a *Point) subtract(b *Point) *Vector {
	return &Vector{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func (a *Vector) Scale(s float64) *Vector {
	return &Vector{a.X * s, a.Y * s, a.Z * s}
}

func (a *Vector) dot(b *Vector) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a *Vector) cross(b *Vector) *Vector {
	return &Vector{
		a.Y*b.Z - a.Z*b.Y,
		a.Z*b.X - a.X*b.Z,
		a.X*b.Y - a.Y*b.X,
	}
}

func (a *Vector) normalize() *Vector {
	m := a.magnitude()
	return &Vector{a.X / m, a.Y / m, a.Z / m}
}

func (a *Vector) magnitude() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
}
