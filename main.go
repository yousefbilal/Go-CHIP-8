package main

import (
	"flag"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/jroimartin/gocui"
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
}
func (g *GraphicsHandler) drawGraphics(chip8 *CPU) {
	for y := 0; y < bufferHeight; y++ {
		for x := 0; x < bufferWidth; x++ {
			if chip8.gfx[y*bufferWidth+x] != 0 {
				mat := g.scale. // Scale the sprite to 10x10 pixels
						Moved(pixel.V(float64(x*g.pixelSize+g.pixelSize/2),
						float64(g.windowHeight-y*g.pixelSize-g.pixelSize/2))) // Move it to the correct position
				g.sprite.Draw(g.win, mat)
			}
		}
	}
}

const pixelSize = 10

func gocuiLoop(g *gocui.Gui) {
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func run() {

	fileName := flag.String("file", "", "file name of the program to run")
	UPS := flag.Uint("ups", 30, "set the updates per second")
	flag.Parse()
	chip8 := NewChip8(*fileName)
	if *fileName == "" {
		panic("file name not specified")
	}
	graphics := NewGraphics(bufferWidth, bufferHeight, pixelSize)
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(func(gui *gocui.Gui) error {
		return layout(gui, chip8)
	})

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	go gocuiLoop(g)

	var interval int64 = int64(1000000 / *UPS)
	var accumulator int64 = 0
	lastTime := time.Now().UnixMicro()
	var delta int64

	for !graphics.win.Closed() {

		delta = time.Now().UnixMicro() - lastTime
		lastTime += delta
		accumulator += delta
		if accumulator >= interval {
			graphics.win.Clear(colornames.Black)
			for accumulator >= interval {
				chip8.EmulationCycle()
				g.Update(func(gui *gocui.Gui) error {
					return updateLayout(gui, chip8)
				})
				graphics.setKeys(chip8)
				accumulator -= interval
			}
			graphics.drawGraphics(chip8)
		}
		graphics.win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
