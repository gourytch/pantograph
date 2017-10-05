package main

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

// отрисуем карты покрытия разрешенными состояниями
func (p *Pantograph) Coverage(img *image.RGBA)  {
	N1 := p.E1.MaxVal - p.E1.MinVal + 1
	N2 := p.E2.MaxVal - p.E2.MinVal + 1
	//var R1, R2 image.Rectangle
	gc := gg.NewContextForRGBA(img)
	// first points
	p.A1 = float64(p.E1.MinVal) * p.E1.Step
	p.A2 = float64(p.E2.MinVal) * p.E2.Step
	p.Solve()
	Q1 := p.P1

	cp1 := color.RGBA{255,255,0, 255}
	c0 := color.RGBA{0, 0,0,1}
	c1 := color.RGBA{0, 0,0,10}

	gc.SetRGB(1.0, 1.0, 0.8)
	gc.Clear()

	for a1 := 0; a1 < N1; a1++ {
		for a2 := 0; a2 < N2; a2++ {
			p.A1 = float64(p.E1.MinVal + a1) * p.E1.Step
			p.A2 = float64(p.E2.MinVal + a2) * p.E2.Step
			p.Solve()
			//if p.Valid {
				//R1.Add(image.Point{int(p.P1.X), int(p.P1.Y)})
				//R2.Add(image.Point{int(p.P2.X), int(p.P2.Y)})
			//}

			if p.Valid {
				gc.DrawLine(Q1.X, Q1.Y, p.P1.X, p.P1.Y)
				gc.SetColor(cp1)
				gc.Stroke()
				gc.DrawLine(p.E1.X, p.E1.Y, p.N1.X, p.N1.Y)
				gc.DrawLine(p.E2.X, p.E2.Y, p.N2.X, p.N2.Y)
				gc.SetColor(c1)
				gc.Stroke()
			} else {
				gc.DrawLine(p.E1.X, p.E1.Y, p.N1.X, p.N1.Y)
				gc.DrawLine(p.E2.X, p.E2.Y, p.N2.X, p.N2.Y)
				gc.SetColor(c0)
				gc.Stroke()
			}

			Q1 = p.P1
		}
	}
}
