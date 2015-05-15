package main

import . "./godray"

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

var i *Vector = &Vector{1, 0, 0}
var j *Vector = &Vector{0, 1, 0}
var k *Vector = &Vector{0, 0, 1}
var o *Point = &Point{0, 0, 0}

const MAX_DEPTH int = 10

func getClosestIntersection(ray *Ray, objects []Object) (*Intersection, Object) {
	intersections := make([]*Intersection, len(objects))

	for i, object := range objects {
		point, dist, n := object.Intersect(ray)
		intersections[i] = &Intersection{point, dist, n}
	}

	var closestIntersection *Intersection
	var closestObject Object
	for i, intersection := range intersections {
		if closestIntersection == nil ||
			intersection.Distance < closestIntersection.Distance {
			closestIntersection = intersection
			closestObject = objects[i]
		}
	}

	return closestIntersection, closestObject
}

func raytrace(ray *Ray, lights []*Light, objects []Object, depth int) *Color {
	var ambientTerm *Color = Black
	var diffuseTerm *Color = Black
	var specularTerm *Color = Black

	if depth <= 0 {
		return nil
	}

	closestIntersection, closestObject := getClosestIntersection(ray, objects)
	intersection := closestIntersection.Point
	n := closestIntersection.Normal

	vv := ray.V.Normalize().Scale(-1)

	// Phong equaion calculation
	// for details, please refer to https://en.wikipedia.org/wiki/Phong_reflection_model
	if intersection != nil {
		for _, light := range lights {
			ambientTerm = ambientTerm.Add(closestObject.Material().Ambient.
				Multiply(light.Ambient))

			l := light.Position.Subtract(intersection).Normalize()

			rayToLight := &Ray{
				intersection,
				l,
			}

			obstruction, _ := getClosestIntersection(rayToLight, objects)
			// shadow calculation - if obstructed, we skip coloring.
			if obstruction.Point != nil {
				continue
			}

			r := n.Scale(2 * l.Dot(n)).Subtract(l)

			if l.Dot(n) > 0 {
				diffuseTerm = diffuseTerm.Add(light.Diffuse.Scale(l.Dot(n)).
					Multiply(closestObject.Material().Diffuse))
			}

			if r.Dot(vv) > 0 {
				specularTerm = specularTerm.Add(light.Specular.
					Scale(math.Pow(r.Dot(vv), closestObject.Material().Shininess)).
					Multiply(closestObject.Material().Specular))
			}
		}

		// recurse
		newRay := &Ray{intersection, n.Scale(vv.Dot(n)).Scale(2).Subtract(vv).Normalize()}
		// reflection calculation - treating the reflected color as a light source
		reflectedColor := raytrace(newRay, lights, objects, depth-1)

		if reflectedColor != nil {
			if newRay.V.Dot(n) > 0 {
				diffuseTerm = diffuseTerm.Add(reflectedColor.Scale(newRay.V.Dot(n)).
					Multiply(closestObject.Material().Diffuse).Scale(closestObject.Material().Reflectivity))
			}

			if newRay.V.Dot(vv) > 0 {
				specularTerm = specularTerm.Add(reflectedColor.
					Scale(math.Pow(newRay.V.Dot(vv), closestObject.Material().Shininess)).
					Multiply(closestObject.Material().Specular).Scale(closestObject.Material().Reflectivity))
			}
		}

		return ambientTerm.Add(diffuseTerm).Add(specularTerm)
	}

	return nil
}

func main() {
	eye := o
	camera := NewCamera(eye.Add(k.Scale(10)), k.Scale(-1), j)
	screen := &Screen{Width: 800, Height: 600, Fov: 45}

	// lights
	lights := []*Light{
		&Light{
			&Point{10, 4, 2},
			White.Scale(0.1),
			White,
			White,
		},
	}

	// objects
	objects := []Object{
		NewSphere(&Point{5, 0, -2}, 2, &Material{
			&Color{color.RGBA{250, 60, 60, 255}},
			&Color{color.RGBA{250, 60, 60, 255}},
			White,
			20,
			0.4,
		}),
		NewSphere(&Point{-2, 2, -4}, 3, &Material{
			&Color{color.RGBA{60, 250, 60, 255}},
			&Color{color.RGBA{60, 250, 60, 255}},
			White,
			50,
			0.9,
		}),
		NewSphere(&Point{0, -10, -4}, 8, &Material{
			Blue,
			&Color{color.RGBA{100, 100, 200, 255}},
			White,
			50,
			0.9,
		}),
		NewSphere(&Point{0, 4, 0}, 1, &Material{
			&Color{color.RGBA{250, 225, 150, 255}},
			&Color{color.RGBA{250, 225, 150, 255}},
			White,
			30,
			0.6,
		}),
	}

	// Uncomment these to test the rays and intersection.
	// miss := &Ray{&Point{0, 5, 0}, &Vector{0, 0, -4}}
	// _, _, n := sphere.Intersect(hit)
	// intersection, t := sphere.Intersect(miss)
	// fmt.Println(n)

	out, err := os.Create("./output.png")
	if err != nil {
		os.Exit(1)
	}

	imgRect := image.Rect(0, 0, screen.Width, screen.Height)
	img := image.NewRGBA(imgRect)
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.ZP, draw.Src)

	for u := 0; u < screen.Width; u++ {
		for v := 0; v < screen.Height; v++ {
			ray := camera.GetRayTo(screen, u, v)
			illumination := raytrace(ray, lights, objects, MAX_DEPTH)
			if illumination == nil {
				illumination = Black
			}

			fill := &image.Uniform{color.RGBA{
				illumination.R,
				illumination.G,
				illumination.B,
				illumination.A,
			}}

			draw.Draw(img, image.Rect(u, v, u+1, v+1), fill, image.ZP, draw.Src)
		}
	}

	err = png.Encode(out, img)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
