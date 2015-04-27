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
	"sync"
)

var i *Vector = &Vector{1, 0, 0}
var j *Vector = &Vector{0, 1, 0}
var k *Vector = &Vector{0, 0, 1}
var o *Point = &Point{0, 0, 0}

func getClosestIntersection(ray *Ray,
	objects []Object) (*Intersection, Object) {
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

func raytrace(wg *sync.WaitGroup, u, v int, ray *Ray, lights [2]*Light, objects []Object, img *image.RGBA) {
	closestIntersection, closestObject := getClosestIntersection(ray, objects)
	intersection := closestIntersection.Point
	n := closestIntersection.Normal

	if intersection != nil {
		var illumination *Color = Black

		for _, light := range lights {
			illumination = illumination.Add(closestObject.Material().Ambient.
				Multiply(light.Ambient))

			l := light.Position.Subtract(intersection).Normalize()

			rayToLight := &Ray{
				intersection,
				l,
			}

			obstruction, _ := getClosestIntersection(rayToLight, objects)
			if obstruction.Point != nil {
				continue
			}

			r := n.Scale(2 * l.Dot(n)).Subtract(l)
			vv := ray.V.Normalize().Scale(-1)

			diffuseTerm := light.Diffuse.Scale(l.Dot(n)).
				Multiply(closestObject.Material().Diffuse)
			specularTerm := light.Specular.
				Scale(math.Pow(r.Dot(vv), closestObject.Material().Shininess)).
				Multiply(closestObject.Material().Specular)

			if l.Dot(n) > 0 {
				illumination = illumination.Add(diffuseTerm)
			}

			if r.Dot(vv) > 0 {
				illumination = illumination.Add(specularTerm)
			}

		}

		fill := &image.Uniform{color.RGBA{
			illumination.R,
			illumination.G,
			illumination.B,
			illumination.A,
		}}

		draw.Draw(img, image.Rect(u, v, u+1, v+1), fill, image.ZP, draw.Src)
	}

	wg.Done()
}

func main() {
	eye := o
	camera := NewCamera(eye.Add(k.Scale(10)), k.Scale(-1), j)
	screen := &Screen{800, 600, 45}

	// lights
	lights := [...]*Light{
		&Light{
			&Point{0, 4, -4},
			White.Scale(0.1),
			White,
			White,
		},
		&Light{
			&Point{10, 4, 2},
			White.Scale(0.1),
			White,
			White,
		},
	}

	// objects
	objects := []Object{
		NewSphere(&Point{0, 0, -4}, 1, &Material{
			Red,
			Red,
			White,
			20,
		}),
		NewSphere(&Point{-2, 2, -4}, 1, &Material{
			Green,
			Green,
			White,
			20,
		}),
		NewSphere(&Point{2, -4.5, -4}, 3, &Material{
			Blue,
			Blue,
			White,
			20,
		}),
	}

	//hit := &Ray{&Point{0, 0, 0}, &Vector{-0.01, 0.01, -1}}
	//miss := &Ray{&Point{0, 5, 0}, &Vector{0, 0, -4}}
	//_, _, n := sphere.Intersect(hit)
	//intersection, t := sphere.Intersect(miss)
	//fmt.Println(n)

	out, err := os.Create("./output.png")
	if err != nil {
		os.Exit(1)
	}

	imgRect := image.Rect(0, 0, screen.Width, screen.Height)
	img := image.NewRGBA(imgRect)
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.ZP, draw.Src)

	//runtime.GOMAXPROCS(1)
	wg := sync.WaitGroup{}

	for u := 0; u < screen.Width; u++ {
		for v := 0; v < screen.Height; v++ {
			wg.Add(1)
			ray := camera.GetRayTo(screen, u, v)
			go raytrace(&wg, u, v, ray, lights, objects, img)
		}
	}

	wg.Wait()

	err = png.Encode(out, img)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
