package godray

import (
	"math"
)

const epsilon = 1e-4

type Object interface {
	Intersect(*Ray) (*Point, float64, *Vector)
	Material() *Material
	SetMaterial(*Material)
	Normal(*Point) *Vector
}

type Sphere struct {
	Center   *Point
	Radius   float64
	material *Material
}

type Intersection struct {
	Point    *Point
	Distance float64
	Normal   *Vector
}

// Sphere methods
func NewSphere(center *Point, radius float64, material *Material) Sphere {
	return Sphere{center, radius, material}
}

func (s Sphere) Material() *Material {
	return s.material
}

func (s Sphere) SetMaterial(material *Material) {
	s.material = material
}

func (s Sphere) Normal(p *Point) *Vector {
	return p.Subtract(s.Center).Normalize()
}

func (s Sphere) Intersect(r *Ray) (*Point, float64, *Vector) {
	l := r.V.Normalize()
	o := r.P

	leftSide := l.Dot(o.Subtract(s.Center))

	diff := r.P.Subtract(s.Center)
	sqrtSum := (leftSide * leftSide) - diff.Dot(diff) + s.Radius*s.Radius

	if sqrtSum == 0.0 {
		d := -1 * leftSide
		intersectPoint := o.Add(r.V.Scale(d))

		if d > epsilon {
			return intersectPoint, d, s.Normal(intersectPoint)
		}
	} else if sqrtSum > 0 {
		d1 := -1*leftSide + math.Sqrt(sqrtSum)
		d2 := -1*leftSide - math.Sqrt(sqrtSum)

		pt1 := o.Add(r.V.Normalize().Scale(d1))
		pt2 := o.Add(r.V.Normalize().Scale(d2))

		if r.P.Subtract(pt1).Magnitude()-r.P.Subtract(pt2).Magnitude() >= 0 {
			if d2 > epsilon {
				return pt2, d2, s.Normal(pt2)
			}
		} else if d1 > epsilon {
			return pt1, d1, s.Normal(pt1)
		}
	}

	return nil, math.MaxFloat64, nil
}
