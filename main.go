package main

import (
	"fmt"
	"github.com/samuel/go-pcx/pcx"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"log"
	"math"
	"time"
)

var (
	wnd              *sdlcanvas.Window
	cv               *canvas.Canvas
	tileset          map[string]*canvas.Image
	tilefiles        []string
	offset           int
	board            Board
	selector         Board
	selected_tile    string
	scale            float64
	cscale, nscale   float64
	screenw, screenh = 1280, 720
	rescale          bool
)

func load() {
	tilefiles = []string{"forest", "grass", "marsh", "village", "rocket", "water"}
	tileset = make(map[string]*canvas.Image)

	err := pcx.FormatError("")
	if err != "" {
	}

	for _, tf := range tilefiles {
		file := "u2_" + tf + ".pcx"
		img2, err := canvas.LoadImage(file)
		if err != nil {
			log.Fatalf("failed to load image: %v", err)
		}
		tileset[tf] = img2
	}

	selected_tile = "grass"
	img := tileset[selected_tile]

	offset = img.Width() / 2
	scale = float64(img.Width())
	println("scale = ", scale)
	cscale = 1.0
}

func main() {
	var err error
	wnd, cv, err = sdlcanvas.CreateWindow(screenw, screenh, "Tile Map")
	if err != nil {
		log.Println(err)
		return
	}
	defer wnd.Destroy()

	load()

	/*
		sdlsurface, err := wnd.Window.GetSurface()
		if err != nil {
			log.Println(err)
			return
		}
	*/

	cv.SetFont("SometypeMono-Medium.ttf", 12) /// SometypeMono-Regular Righteous-Regular

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
			cscale += 0.05
			nscale = 1.05
			rescale = true
		}
		if y == -1 {
			cscale -= 0.05
			nscale = 0.95
			rescale = true
		}
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

		//cv.Save()

		w, h := float64(cv.Width()), float64(cv.Height())
		// Clear the screen
		cv.SetFillStyle("#000")
		cv.FillRect(0, 0, w, h)
		if rescale { /// http://jsfiddle.net/mBzVR/4/
			nw := w * nscale
			nh := h * nscale			
			//nw := w * cscale
			//nh := h * cscale
			cv.Translate(-((nw - w) / 2), -((nh - h) / 2))
			cv.Scale(nscale, nscale)
			//cv.Scale(cscale, cscale)
			rescale = false
		}



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
				cv.DrawImage(tileset[t.Type], float64(t.X), float64(t.Y))
			}
		}

		cv.SetFillStyle("#778899")
		cv.FillText(fmt.Sprintf("x:%d  y:%d", int(tx), int(ty)), tx, ty-2.0)

		//cv.Restore()

	})
}

func fit_gridf(mx, my float64) (tx, ty float64) {
	//nxt := offset * 2
	
	nxt := offset * 2 //* int(cscale)
	
	
	//nxt := int(float64(offset*2)/ cscale*cscale)

	//nx := int((float64(cv.Width())/float64(screenw)* cscale) * mx)
	//ny := int((float64(cv.Height())/float64(screenh)* cscale) * my)

	//nx := int((float64(cv.Width())* cscale/float64(screenw)) * mx)
	//ny := int((float64(cv.Height())* cscale/float64(screenh)) * my)

	//nx := int((float64(cv.Width())/float64(screenw) )* mx)
	//ny := int((float64(cv.Height())/float64(screenh) )* my)


	//tx = ( mx*float64(nxt)/cscale*float64(nxt))
	//ty = ( my*float64(nxt)/cscale*float64(nxt))

	nx, ny := int(mx), int(my)
	//tx, ty = float64((nx/nxt)*nxt)*cscale, float64((ny/nxt)*nxt)*cscale
	//nx, ny := int(mx/cscale), int(my/cscale)
	//nxt := float64(offset * 2) 	*cscale
	//nx, ny := mx, my
	
	tx, ty = float64((nx/nxt)*nxt), float64((ny/nxt)*nxt)
	//tx, ty = tx/cscale, ty/cscale
	//nx, ny := int(mx*cscale), int(my*cscale)
	//tx, ty = float64(nx/nxt), float64(ny/nxt)

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

/*

func new_grid(cv *canvas.Canvas) {
	cs := scale * cscale
	penwidth := 1.0
	ix, iy := cs*2, cs*2
	vstep, hstep := cs, cs
	step := 1.0 * cs

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
*/
