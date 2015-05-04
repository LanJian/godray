package godray

type Material struct {
	Ambient      *Color
	Diffuse      *Color
	Specular     *Color
	Shininess    float64
	Reflectivity float64
}
