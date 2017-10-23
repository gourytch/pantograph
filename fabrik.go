package main

import (
	"sort"
)

// Вычисление всех возможных валидных состояний и позиций
func (p *Pantograph) Evaluate() {
	N1 := p.E1.MaxVal - p.E1.MinVal + 1
	N2 := p.E2.MaxVal - p.E2.MinVal + 1
	p.Cloud = make([]Projection, 0, N1*N2)
	for a1 := 0; a1 < N1; a1++ {
		for a2 := 0; a2 < N2; a2++ {
			p.A1 = float64(p.E1.MinVal+a1) * p.E1.Step
			p.A2 = float64(p.E2.MinVal+a2) * p.E2.Step
			p.Solve()
			if p.Valid {
				p.Cloud = append(p.Cloud, Projection{
					Step1: a1 + p.E1.MinVal,
					Step2: a2 + p.E2.MinVal,
					Pos:   p.P1,
				})
			}
		}
	}
}

type ByError []Match

func (a ByError) Len() int           { return len(a) }
func (a ByError) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByError) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

// построение списка индексов облака, примерно подходящих под координату
func (p *Pantograph) MatchPosition(pos Position, maxDist float64) (matches []Match) {
	matches = make([]Match, 0, 10)
	// сперва собираем
	for ix, prj := range p.Cloud {
		dist := Distance(prj.Pos, pos)
		if maxDist < dist {
			continue // позиция за границей
		}
		matches = append(matches, Match{prj, ix, dist})
	}
	// потом сортируем по возрастанию (меньший - вперёд)
	sort.Reverse(ByError(matches))
	// наконец, возвращаем
	return
}

// построение списка индексов облака, примерно подходящих под перемещение из позиции
// предполагается, что перемещений через запрещенные состояния не происходит
// сортируется по минимизации необходимых телодвижений для перехода.
func (p *Pantograph) MatchMove(from Match, to Position, maxDist float64, maxRot2 int) (matches []Match) {
	matches = make([]Match, 0, 10)
	// сперва собираем
	for ix, prj := range p.Cloud {
		err := Distance(prj.Pos, to)
		if maxDist < err {
			continue // позиция за границей
		}
		// ключом сортировки будет сумма квадратов количества доворотов двигателей
		ds1 := prj.Step1 - from.Step1
		ds2 := prj.Step2 - from.Step2
		ds := ds1*ds1 + ds2*ds2
		if maxRot2 < ds {
			continue // слишком много вертеть
		}
		matches = append(matches, Match{prj, ix, float64(ds)})
	}
	// потом сортируем по возрастанию (меньший - вперёд)
	sort.Reverse(ByError(matches))
	// наконец, возвращаем
	return
}
