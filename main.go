package main

import (
	"image"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type GraphicsHandler struct {
	sprite       *pixel.Sprite
	scale        *pixel.Matrix
	win          *pixelgl.Window
	windowWidth  int
	windowHeight int
	pixelSize    int
}

func NewGraphics(bufferWidth, bufferHeight, pixelSize int) *GraphicsHandler {
	windowWidth := bufferWidth * pixelSize
	windowHeight := bufferHeight * pixelSize
	cfg := pixelgl.WindowConfig{
		Title:  "CHIP-8 Display",
		Bounds: pixel.R(0, 0, float64(windowWidth), float64(windowHeight)),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.White)
	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())
	scale := pixel.IM.Scaled(pixel.ZV, float64(pixelSize))

	return &GraphicsHandler{
		sprite:       sprite,
		scale:        &scale,
		win:          win,
		windowWidth:  windowWidth,
		windowHeight: windowHeight,
		pixelSize:    pixelSize,
	}
}

const (
	bufferWidth  = 64
	bufferHeight = 32
	pixelSize    = 10
)

func (g *GraphicsHandler) drawGraphics(chip8 *CPU) {
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if chip8.gfx[y*64+x] != 0 {
				mat := g.scale. // Scale the sprite to 10x10 pixels
						Moved(pixel.V(float64(x*g.pixelSize+g.pixelSize/2),
						float64(g.windowHeight-y*g.pixelSize-g.pixelSize/2))) // Move it to the correct position
				g.sprite.Draw(g.win, mat)
			}
		}
	}
}

func run() {

	chip8 := NewChip8("file")

	g := NewGraphics(bufferWidth, bufferHeight, pixelSize)

	g.win.Clear(colornames.Black)

	for !g.win.Closed() {
		g.drawGraphics(chip8)
		g.win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
