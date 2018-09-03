package main

import (
	"github.com/samuel/go-pcx/pcx"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"image"
	"log"
	"math"
	"os"
	"time"
)

var (
	img    image.Image
	offset int
	scale  float64
)

type tile struct {
	x, y float64
}

func init() {
	file, err := os.Open("u2_village.pcx")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	img, err = pcx.Decode(file)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	
	offset = img.Bounds().Dx()/2
	scale = float64(img.Bounds().Dx())
}

func main() {
	wnd, cv, err := sdlcanvas.CreateWindow(1280, 720, "Tile Map")
	if err != nil {
		log.Println(err)
		return
	}
	defer wnd.Destroy()

	var mx, my, action float64
	tiles := make([]tile, 0, 100)

	wnd.MouseMove = func(x, y int) {
		mx, my = float64(x), float64(y)
	}

	wnd.MouseDown = func(button, x, y int) {
		action = 1
		nxt := offset*2
		nx, ny := int(mx), int(my)
		tx, ty := float64((nx / nxt) * nxt), float64((ny / nxt) * nxt) // round to nearest fit
		tiles = append(tiles, tile{x: tx, y: ty})
	}

	wnd.KeyDown = func(scancode int, rn rune, name string) {
		switch name {
		case "Escape":
			wnd.Close()
		case "Space":
			action = 1
			tiles = append(tiles, tile{x: mx, y: my})
		case "Enter":
			action = 1
			tiles = append(tiles, tile{x: mx, y: my})
		}
	}
	wnd.SizeChange = func(w, h int) {
		cv.SetBounds(0, 0, w, h)
	}

	lastTime := time.Now()

	wnd.MainLoop(func() {
		now := time.Now()
		diff := now.Sub(lastTime)
		lastTime = now
		action -= diff.Seconds() * 3
		action = math.Max(0, action)

		w, h := float64(cv.Width()), float64(cv.Height())

		// Clear the screen
		cv.SetFillStyle("#000")
		cv.FillRect(0, 0, w, h)

		NewGrid(cv)

		// Draw a circle around the cursor
		cv.SetStrokeStyle("#F00")
		cv.SetLineWidth(2)
		cv.BeginPath()
		cv.Arc(mx, my, 24+action*24, 0, math.Pi*2, false)
		cv.Stroke()

		// Draw tiles where the user has clicked
		for _, t := range tiles {
			cv.PutImageData(img.(*image.RGBA), int(t.x), int(t.y))
		}
	})
}

func NewGrid(cv *canvas.Canvas) {
	penwidth := 1.0
	ix, iy := scale*2,scale*2
	vstep, hstep := scale, scale
	step := 1.0 * scale
	
	for x := ix; x <= hstep*step; x += step {
		cv.SetStrokeStyle(int(25), 255, 255)
		cv.SetLineWidth(penwidth)
		cv.BeginPath()
		cv.MoveTo(x, 0)
		cv.LineTo(x, vstep*step)
		cv.Stroke()
	}

	for y := iy; y <= vstep*step; y += step {
		cv.SetStrokeStyle(int(25), 255, 255)
		cv.SetLineWidth(penwidth)
		cv.BeginPath()
		cv.MoveTo(0, y)
		cv.LineTo(hstep*step, y)
		cv.Stroke()
	}
}
