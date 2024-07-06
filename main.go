package main

import (
	"flag"
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

/*
Keypad                   Keyboard
+-+-+-+-+                +-+-+-+-+
|1|2|3|C|                |1|2|3|4|
+-+-+-+-+                +-+-+-+-+
|4|5|6|D|                |Q|W|E|R|
+-+-+-+-+       =>       +-+-+-+-+
|7|8|9|E|                |A|S|D|F|
+-+-+-+-+                +-+-+-+-+
|A|0|B|F|                |Z|X|C|V|
+-+-+-+-+                +-+-+-+-+
*/
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

func (g *GraphicsHandler) setKeys(chip8 *CPU) {
	buttons := []pixelgl.Button{
		pixelgl.KeyX, pixelgl.Key1, pixelgl.Key2, pixelgl.Key3,
		pixelgl.KeyQ, pixelgl.KeyW, pixelgl.KeyE, pixelgl.KeyA,
		pixelgl.KeyS, pixelgl.KeyD, pixelgl.KeyZ, pixelgl.KeyC,
		pixelgl.Key4, pixelgl.KeyR, pixelgl.KeyF, pixelgl.KeyV,
	}
	for i, v := range buttons {
		chip8.keys[i] = g.win.Pressed(v)
	}

	// fmt.Printf(
	// 	`
	// 	+-+-+-+-+
	// 	|%v|%v|%v|%v|
	// 	+-+-+-+-+
	// 	|%v|%v|%v|%v|
	// 	+-+-+-+-+
	// 	|%v|%v|%v|%v|
	// 	+-+-+-+-+
	// 	|%v|%v|%v|%v|
	// 	+-+-+-+-+
	// 	`,
	// 	chip8.keys[1], chip8.keys[2], chip8.keys[3], chip8.keys[0xC],
	// 	chip8.keys[4], chip8.keys[5], chip8.keys[6], chip8.keys[0xD],
	// 	chip8.keys[7], chip8.keys[8], chip8.keys[9], chip8.keys[0xE],
	// 	chip8.keys[0xA], chip8.keys[0], chip8.keys[0xB], chip8.keys[0xF],
	// )
}
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

const pixelSize = 10

func run() {

	fileName := flag.String("file", "", "file name of the program to run")
	flag.Parse()
	chip8 := NewChip8(*fileName)
	if *fileName == "" {
		panic("file name not specified")
	}

	g := NewGraphics(bufferWidth, bufferHeight, pixelSize)

	g.win.Clear(colornames.Black)

	for !g.win.Closed() {
		g.win.Clear(colornames.Black)
		chip8.EmulationCycle()
		g.drawGraphics(chip8)
		g.win.Update()
		g.setKeys(chip8)
		// time.Sleep(1 * time.Millisecond)
	}
}

func main() {
	pixelgl.Run(run)
}
