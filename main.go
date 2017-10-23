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
	buf      screen.Buffer
	win      screen.Window
	tx       screen.Texture
	Size     = image.Pt(800, 800)
	Bounds   = image.Rect(0, 0, Size.X, Size.Y)
	coverage *image.RGBA
	use_IK   = true // use inverse kinematic
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

const NUM_STEPS = 200
const MAX_ROT = (NUM_STEPS / 10) * (NUM_STEPS / 10)
const MAX_DIST = 3.0
const FULL_ANGLE = math.Pi / 2
const MIN_STEP_VAL = -NUM_STEPS / 2
const MAX_STEP_VAL = +NUM_STEPS / 2

var p = &Pantograph{
	E1: Engine{
		Position: Position{X: 200.0, Y: 500.0},
		R:        200.0,
		A:        math.Pi / 2,
		Step:     FULL_ANGLE / NUM_STEPS,
		MinVal:   MIN_STEP_VAL,
		MaxVal:   MAX_STEP_VAL},
	E2: Engine{
		Position: Position{X: 600.0, Y: 500.0},
		R:        200.0,
		A:        math.Pi / 2,
		Step:     FULL_ANGLE / NUM_STEPS,
		MinVal:   MIN_STEP_VAL,
		MaxVal:   MAX_STEP_VAL},
	L1: 250.0,
	L2: 250.0,
	A1: 0.1,
	A2: -0.2,
}

var histCap = 100

type History []int

var hist1 History = make(History, 0, histCap)
var hist2 History = make(History, 0, histCap)

func (h History) push(a int) {
}

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
	// рисуем историю доворотов
	xOffs := p.E1.X - float64(histCap/2)
	yOffs := p.E1.Y + 50

	for i, v := range hist1 {
		d := 200.0 * 2 * float64(v) / NUM_STEPS;
		gc.DrawLine(xOffs+float64(i), yOffs, xOffs+float64(i), yOffs+d)
		if 0 < d {
			gc.SetRGBA(1, 0, 0, 1)
		} else {
			gc.SetRGBA(0, 0, 1, 1)
		}
		gc.Stroke()
	}

	xOffs = p.E2.X - float64(histCap/2)
	yOffs = p.E2.Y + 50
	for i, v := range hist2 {
		d := 200.0 * 2 * float64(v) / NUM_STEPS;
		gc.DrawLine(xOffs+float64(i), yOffs, xOffs+float64(i), yOffs+d)
		if 0 < d {
			gc.SetRGBA(1, 0, 0, 1)
		} else {
			gc.SetRGBA(0, 0, 1, 1)
		}
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
		p.Evaluate()
		var matches []Match = nil

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
				if use_IK {
					pt := Position{float64(e.X), float64(e.Y)}
					cur := Match{
						Projection{
							int(p.A1 / p.E1.Step),
							int(p.A2 / p.E2.Step),
							p.P1},
						-1,
						0.0}
					matches = p.MatchMove(cur, pt, MAX_DIST, MAX_ROT)
					if 0 < len(matches) {
						m := matches[0]
						p.A1 = float64(m.Step1) * p.E1.Step
						p.A2 = float64(m.Step2) * p.E2.Step
						ds1 := m.Step1 - cur.Step1
						ds2 := m.Step2 - cur.Step2
						if len(hist1) < histCap {
							hist1 = append(hist1, ds1)
							hist2 = append(hist2, ds2)
						} else {
							hist1 = append(hist1[len(hist1)-histCap:], ds1)
							hist2 = append(hist2[len(hist2)-histCap:], ds2)
						}
					}

				} else {
					p.A1 = math.Pi * float64(e.X-400) / 400
					p.A2 = math.Pi * float64(e.Y-400) / 400
				}
			}
		}
	})
}
