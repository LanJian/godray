package godray

import "image/color"

type Material struct {
	Diffuse   color.Color
	Specular  color.Color
	Shininess float64
}
