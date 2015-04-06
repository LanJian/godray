package godray

import "math"

type camera struct {
	Position *Point
	View     *Vector
	Up       *Vector
}

func NewCamera(position *Point, view *Vector, up *Vector) *camera {
	return &camera{position, view.Normalize(), up.Normalize()}
}

func (c *camera) side() *Vector {
	return c.View.Cross(c.Up)
}

func (c *camera) GetRayTo(screen *Screen, u int, v int) *Ray {
	w, h := float64(screen.Width), float64(screen.Height)
	fovy := float64(screen.Fov)

	d := (h / 2) / math.Tan(fovy/2)
	sideVector := c.side().Scale(float64(u) - w/2)
	upVector := c.Up.Scale(h/2 - float64(v))
	viewVector := c.View.Scale(d)

	return &Ray{
		c.Position,
		sideVector.Add(upVector).Add(viewVector).Normalize(),
	}
}
