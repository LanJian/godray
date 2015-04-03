package godray

import "image/color"

type Light struct {
	Position *Point
	Ambient  color.Color
	Diffuse  color.Color
	Specular color.Color
}
