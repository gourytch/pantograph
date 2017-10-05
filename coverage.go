package main

import (
	"math"
	"image"
)

type CoordMap
// отрисуем полные карты покрытия
func (p *Pantograph) Coverage() *image.RGBA {
	N1 := p.E1.MaxVal - p.E1.MinVal + 1
	N2 := p.E2.MaxVal - p.E2.MinVal + 1
	pix_avail := image.NewRGBA(image.Rect(0,0, N1, N2))
	for a1 := 0; a1 < N1; a1++ {
		for a2 := 0; a2 < N2; a2++ {
			p.A1 = float64(p.E1.MinVal + a1) * p.E1.Step
			p.A2 = float64(p.E2.MinVal + a1) * p.E2.Step
			p.Solve()
			if p.Valid {

			}
		}
	}

}
