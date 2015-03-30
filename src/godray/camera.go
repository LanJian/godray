package godray

import "math"

type camera struct {
	position *Point
	View     *Vector
	up       *Vector
}

func NewCamera(position *Point, view *Vector, up *Vector) *camera {
	return &camera{position, view.normalize(), up.normalize()}
}

func (c *camera) side() *Vector {
	return c.View.cross(c.up)
}

func (c *camera) GetRayTo(screen *Screen, u int, v int) *Ray {
	w, h := float64(screen.Width), float64(screen.Height)
	fovy := float64(screen.Fov)

	d := (h / 2) / math.Tan(fovy/2)
	sideVector := c.side().Scale(float64(u) - w/2)
	upVector := c.up.Scale(float64(v) - h/2)
	viewVector := c.View.Scale(d)

	return &Ray{
		c.position,
		sideVector.add(upVector).add(viewVector).normalize(),
	}
}
