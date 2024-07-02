package main

import (
	"image"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	windowWidth  = 640
	windowHeight = 320
	bufferWidth  = 64
	bufferHeight = 32
	pixelSize    = 10
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "CHIP-8 Display",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight), // Scale up for visibility
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	chip8 := NewChip8("file")

	chip8.gfx[0] = 1
	chip8.gfx[1] = 1
	chip8.gfx[2] = 1
	chip8.gfx[64] = 1
	chip8.gfx[66] = 1
	chip8.gfx[128] = 1
	chip8.gfx[129] = 1
	chip8.gfx[130] = 1

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.White) // Set the single pixel to white

	// Create a picture from the image and then a sprite from the picture
	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())
	scale := pixel.IM.Scaled(pixel.ZV, pixelSize)
	for !win.Closed() {
		win.Clear(colornames.Black)
		for y := 0; y < 32; y++ {
			for x := 0; x < 64; x++ {
				if chip8.gfx[y*64+x] != 0 {
					mat := scale. // Scale the sprite to 10x10 pixels
							Moved(pixel.V(float64(x*pixelSize+pixelSize/2), float64(windowHeight-y*pixelSize-pixelSize/2))) // Move it to the correct position
					sprite.Draw(win, mat)
				}
			}
		}
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
