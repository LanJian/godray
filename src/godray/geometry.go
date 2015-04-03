package godray

import "math"

type Object interface {
	Intersect(*Ray) (*Point, float64)
	Material() *Material
	SetMaterial(*Material)
}

type Sphere struct {
	Center   *Point
	Radius   float64
	material *Material
}

func NewSphere(center *Point, radius float64, material *Material) *Sphere {
	return &Sphere{center, radius, material}
}

func (s *Sphere) Material() *Material {
	return s.material
}

func (s *Sphere) SetMaterial(material *Material) {
	s.material = material
}

func (s *Sphere) Normal(p *Point) *Vector {
	return p.Subtract(s.Center).Normalize()
}

func (s *Sphere) Intersect(r *Ray) (*Point, float64, *Vector) {
	l := r.V.Normalize()
	o := r.P

	leftSide := l.Dot(o.Subtract(s.Center))

	diff := r.P.Subtract(s.Center)
	sqrtSum := (leftSide * leftSide) - diff.Dot(diff) + s.Radius*s.Radius

	if sqrtSum < 0 {
		return nil, 0, nil
	} else if sqrtSum == 0.0 {
		d := -1 * leftSide
		intersectPoint := o.Add(r.V.Scale(d))
		return intersectPoint, d, s.Normal(intersectPoint)
	} else {
		d1 := -1*leftSide + math.Sqrt(sqrtSum)
		d2 := -1*leftSide - math.Sqrt(sqrtSum)

		pt1 := o.Add(r.V.Normalize().Scale(d1))
		pt2 := o.Add(r.V.Normalize().Scale(d2))

		if r.P.Subtract(pt1).Magnitude()-r.P.Subtract(pt2).Magnitude() >= 0 {
			return pt2, d2, s.Normal(pt2)
		} else {
			return pt1, d1, s.Normal(pt1)
		}
	}
}
