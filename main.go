package main

import (
	"fmt"
	//"github.com/veandco/go-sdl2/gfx"
	"github.com/samuel/go-pcx/pcx"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"log"
	"math"
	"os"
	"time"
)

var (
	tileset             map[string]*canvas.Image
	tilefiles           []string
	offset              int
	scale               float64
	board               Board
	selector            Board
	selected_tile       string
	cscale, cmultiplier float64
)

func init() {
	tilefiles = []string{"forest", "grass", "marsh", "village", "rocket", "water"}
	tileset = make(map[string]*canvas.Image)

	for _, tf := range tilefiles {
		file, err := os.Open("u2_" + tf + ".pcx")
		if err != nil {
			log.Fatalf("failed to open image: %v", err)
		}

		img, err := pcx.Decode(file)
		if err != nil {
			log.Fatalf("failed to open image: %v", err)
		}
		
		img2, err := canvas.LoadImage(img)//.(*image.RGBA))
		if err != nil {
			log.Fatalf("failed to load image: %v", err)
		}		
		tileset[tf] = img2
	}

	selected_tile = "grass"
	img := tileset[selected_tile]

	offset = img.Width() / 2
	scale = float64(img.Width())

	cscale = 1.0
	cmultiplier = 0.8

}

func main() {
	wnd, cv, err := sdlcanvas.CreateWindow(1280, 720, "Tile Map")
	if err != nil {
		log.Println(err)
		return
	}
	defer wnd.Destroy()

/*
	sdlsurface, err := wnd.Window.GetSurface()
	if err != nil {
		log.Println(err)
		return
	}	
*/

	cv.SetFont("Righteous-Regular.ttf", 12)

	ofx, ofy := 2, 2
	rows, columns := 20, 10
	board = Board{
		OffsetX:   ofx,
		OffsetY:   ofy,
		Rows:      rows,
		Columns:   columns,
		Positions: make([]Position, rows*columns),
	}

	ofx, ofy = columns+ofx*2, 2
	rows, columns = len(tilefiles), 1
	selector = Board{
		OffsetX:   ofx,
		OffsetY:   ofy,
		Rows:      rows,
		Columns:   columns,
		Positions: make([]Position, rows*columns),
	}

	selector.AddSelectors()

	var mx, my, action float64
	wnd.MouseMove = func(x, y int) {
		mx, my = float64(x), float64(y)
	}

	wnd.MouseDown = func(button, x, y int) {
		if button == 1 { /// mouse left == 1, mouse right == 3
			action = 1
			board.AddTile(x, y)
			selector.SelectTile(x, y)
		}
		if button == 3 {
			action = 1
			board.DeleteTile(x, y)
		}
	}

	wnd.MouseWheel = func(x, y int) {
		action = 1
		if y == 1 {
			cscale -= 1.0 //cmultiplier
		}
		if y == -1 {
			cscale += 1.0 //cmultiplier
		}
		cv.Scale(cscale, cscale)
		//gfx.ZoomSurface(sdlsurface, cscale, cscale, 1)
		println("y:  ", y)
	}

	wnd.KeyDown = func(scancode int, rn rune, name string) {
		switch name {
		case "Escape":
			wnd.Close()
		case "Space":
			action = 1

		case "Enter":
			action = 1

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

		new_grid(cv)

		// Draw a circle around the cursor
		cv.SetStrokeStyle("#778899")
		cv.SetLineWidth(2)
		cv.BeginPath()

		tx, ty := fit_gridf(mx, my)
		open_tl, open_br := action*12, action*24
		cv.Rect(tx-open_tl, ty-open_tl, scale+open_br, scale+open_br)
		cv.Stroke()

		// Draw tiles where the user has clicked
		for _, p := range board.Positions {
			t := p.PTile
			if t != nil {
				cv.DrawImage(tileset[t.Type], float64(t.X), float64(t.Y))
			}
		}

		for _, p := range selector.Positions {
			t := p.PTile
			if t != nil {
				cv.DrawImage(tileset[t.Type], float64(t.X), float64(t.Y)) //.(*image.RGBA
			}
		}

		cv.SetFillStyle("#778899")
		cv.FillText(fmt.Sprintf("x:%d  y:%d", int(tx), int(ty)), tx, ty-2.0)

	})
}

func fit_gridf(mx, my float64) (tx, ty float64) {
	nxt := offset * 2
	nx, ny := int(mx), int(my)
	tx, ty = float64((nx/nxt)*nxt), float64((ny/nxt)*nxt)
	return
}

func new_grid(cv *canvas.Canvas) {
	penwidth := 1.0
	ix, iy := scale*2, scale*2
	vstep, hstep := scale, scale
	step := 1.0 * scale

	for x := ix; x <= hstep*step; x += step {
		cv.SetStrokeStyle("#1e90ff")
		cv.SetLineWidth(penwidth)
		cv.BeginPath()
		cv.MoveTo(x, 0)
		cv.LineTo(x, vstep*step)
		cv.Stroke()
	}

	for y := iy; y <= vstep*step; y += step {
		cv.SetStrokeStyle("#1e90ff")
		cv.SetLineWidth(penwidth)
		cv.BeginPath()
		cv.MoveTo(0, y)
		cv.LineTo(hstep*step, y)
		cv.Stroke()
	}
}
