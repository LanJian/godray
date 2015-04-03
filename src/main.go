package main

import . "./godray"

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	//"runtime"
	"sync"
)

var i *Vector = &Vector{1, 0, 0}
var j *Vector = &Vector{0, 1, 0}
var k *Vector = &Vector{0, 0, 1}
var o *Point = &Point{0, 0, 0}

func main() {
	eye := o
	camera := NewCamera(eye, k.Scale(-1), j)
	screen := &Screen{800, 600, 45}

	// colors
	Red := color.RGBA{255, 0, 0, 255}

	// objects
	sphereMaterial := &Material{
		Red,
		color.White,
		2,
	}
	sphere := NewSphere(&Point{0, 0, -4}, 1, sphereMaterial)

	// lights
	light := &Light{
		&Point{4, 0, -2},
		color.White,
		color.White,
		color.White,
	}

	//hit := &Ray{&Point{0, 0, 0}, &Vector{-0.01, 0.01, -1}}
	//miss := &Ray{&Point{0, 5, 0}, &Vector{0, 0, -4}}
	//_, _, n := sphere.Intersect(hit)
	//intersection, t := sphere.Intersect(miss)
	//fmt.Println(n)

	out, err := os.Create("./output.png")
	if err != nil {
		//fmt.Println(err)
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
			go func(u, v int) {
				ray := camera.GetRayTo(screen, u, v)
				intersection, _, n := sphere.Intersect(ray)

				if intersection != nil {
					l := light.Position.Subtract(intersection)
					r := n.Scale(2 * l.Dot(n)).Subtract(l)
					vv := ray.V.Normalize().Scale(-1)

					// normal map
					//color := color.RGBA{
					//uint8((normal.X + 1) / 2 * 255),
					//uint8((normal.Y + 1) / 2 * 255),
					//uint8((normal.Z + 1) / 2 * 255),
					//255,
					//}
					//fill := &image.Uniform{color}

					draw.Draw(img, image.Rect(u, v, u+1, v+1), fill, image.ZP, draw.Src)
				}
				wg.Done()
			}(u, v)
		}
	}

	wg.Wait()

	err = png.Encode(out, img)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
