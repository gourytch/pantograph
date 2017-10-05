package main

import (
	"image"
	"math"

	"github.com/fogleman/gg"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/image/draw"
	)

var (
	buf  screen.Buffer
	win  screen.Window
	tx   screen.Texture
	Size = image.Pt(800, 800)
	Bounds = image.Rect(0, 0, Size.X, Size.Y)
	coverage *image.RGBA
)

/*
var p = &Pantograph{
	E1: Engine{
			Position: Position{X: 300.0, Y: 400.0},
			R: 100.0,
			A: math.Pi / 2.0,
			Step: math.Pi / 100,
			MinVal: -30,
			MaxVal: +30},
	E2: Engine{
			Position: Position{X: 500.0, Y: 400.0},
			R: 200.0,
			A: math.Pi / 2.0,
			Step: math.Pi / 100,
			MinVal: -30,
			MaxVal: +30,
		},
	L1: 150.0,
	L2: 150.0,
	A1: 0.1,
	A2: -0.2,
}
*/

var p = &Pantograph{
	E1: Engine{
		Position: Position{X: 180.0, Y: 500.0},
		R: 200.0,
		A: math.Pi / 3,
		Step: math.Pi / 100,
		MinVal: -30,
		MaxVal: +30},
	E2: Engine{
		Position: Position{X: 620.0, Y: 500.0},
		R: 200.0,
		A: math.Pi * 2 / 3,
		Step: math.Pi / 100,
		MinVal: -30,
		MaxVal: +30,
	},
	L1: 250.0,
	L2: 250.0,
	A1: 0.1,
	A2: -0.2,
}

//var drawing = []Position{} // набор точек

func render(img *image.RGBA) {
	p.Solve()
	gc := gg.NewContextForRGBA(img)
	// рисуем центр картинки
	gc.DrawLine(380, 400, 420, 400)
	gc.DrawLine(400, 380, 400, 420)
	gc.DrawCircle(400, 400, 10)
	gc.SetRGB(0, 0, 0)
	gc.Stroke()

	// рисуем движки
	gc.DrawCircle(p.E1.X, p.E1.Y, 5.0)
	gc.SetRGB(0, 0, 0)
	gc.Fill()
	gc.DrawCircle(p.E2.X, p.E2.Y, 5.0)
	gc.SetRGB(0, 0, 0)
	gc.Fill()

	if p.Valid {
		// рисуем ведущие тяги
		gc.DrawLine(p.E1.X, p.E1.Y, p.N1.X, p.N1.Y)
		gc.SetRGB(0, 0, 0)
		gc.Stroke()
		gc.DrawLine(p.E2.X, p.E2.Y, p.N2.X, p.N2.Y)
		gc.SetRGB(0, 0, 0)
		gc.Stroke()

		// рисуем ведомые тяги
		gc.DrawLine(p.N1.X, p.N1.Y, p.P1.X, p.P1.Y)
		gc.DrawLine(p.P1.X, p.P1.Y, p.N2.X, p.N2.Y)
		gc.SetRGB(0.5, 0, 0)
		gc.Stroke()

		// рисуем рабочий инструмент
		gc.DrawCircle(p.P1.X, p.P1.Y, 2)
		gc.DrawCircle(p.P1.X, p.P1.Y, 4)
		gc.SetRGB(0, 0, 0)
		gc.Stroke()

		/*
		gc.DrawLine(p.N1.X, p.N1.Y, p.P2.X, p.P2.Y)
		gc.DrawLine(p.P2.X, p.P2.Y, p.N2.X, p.N2.Y)
		gc.SetRGB(0, 0, 0.5)
		gc.Stroke()
		*/
	} else {
		// рисуем ведущие тяги
		gc.DrawLine(p.E1.X, p.E1.Y, p.N1.X, p.N1.Y)
		gc.DrawLine(p.E2.X, p.E2.Y, p.N2.X, p.N2.Y)
		gc.DrawLine(p.N1.X, p.N1.Y, p.P1.X, p.P1.Y)
		gc.DrawLine(p.P2.X, p.P2.Y, p.N2.X, p.N2.Y)
		gc.SetRGBA(0, 0, 0, 0.2)
		gc.Stroke()
	}
}

func publish() {
	win.Publish()
	win.Send(paint.Event{})
}

func main() {
	driver.Main(func(src screen.Screen) {
		win, _ = src.NewWindow(&screen.NewWindowOptions{Size.X, Size.Y, "Pantograph"})
		buf, _ = src.NewBuffer(Size)
		tx, _ = src.NewTexture(Size)
		coverage = image.NewRGBA(Bounds)
		p.Coverage(coverage)
		for {
			switch e := win.NextEvent().(type) {
			case paint.Event:
				// r := image.Rect(200, 200, 400, 400)
				img := buf.RGBA()
				draw.Copy(img, image.Point{}, coverage, Bounds, draw.Over, nil)
				render(img)
				tx.Upload(image.Point{}, buf, buf.Bounds())
				win.Copy(image.Point{}, tx, buf.Bounds(), screen.Over, nil)
				publish()
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
			case key.Event:
				if e.Code == key.CodeEscape {
					return
				}
			case mouse.Event:
				p.A1 = math.Pi * float64(e.X-400) / 400
				p.A2 = math.Pi * float64(e.Y-400) / 400

			}
		}
	})
}
