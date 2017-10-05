package main

import (
	"math"
	"github.com/driusan/de/demodel"
)

type Position struct {
	X, Y float64 // координата в пространстве
}
type Engine struct {
	Position         // положение оси двигателя/ведущей тяги в пространстве
	R        float64 // длина ведущей тяги
	A        float64 // базовый CCW-поворот тяги от абциссы в радианах
	Step     float64 // шаг поворота
	MinVal   int     // минимальное значение шага
	MaxVal   int     // максимальное значение шага
}

type Link struct {
	X1, Y1, X2, Y2 float64
}

type Pantograph struct {
	E1, E2 Engine   // движки
	L1, L2 float64  // длины ведомых тяг
	A1, A2 float64  // отклонения ведущих тяг от "нейтрали"
	N1, N2 Position // вычисленные координаты концов ведущих тяг
	P1     Position // вычисленная координата инструментального узла
	P2     Position // вычисленная координата инструментального узла (альтернативная точка)
	Valid  bool     // True если пантограф находится в разрешенном состоянии
}

func (p *Pantograph) Solve() {
	p.Valid = false // пока у нас не доказано, что всё хорошо

	// вычислим координаты концов ведущих тяг
	Sa1, Ca1 := math.Sincos(p.E1.A + p.A1)
	Sa2, Ca2 := math.Sincos(p.E2.A + p.A2)
	p.N1.X = p.E1.X + p.L1*Ca1
	p.N1.Y = p.E1.Y - p.L1*Sa1
	p.N2.X = p.E2.X + p.L2*Ca2
	p.N2.Y = p.E2.Y - p.L2*Sa2

	Dx := p.N2.X - p.N1.X
	Dy := p.N2.Y - p.N1.Y

	D2 := Dx*Dx + Dy*Dy
	D := math.Sqrt(D2)

	if D == 0 {
		return // совпадение осей. вычислить точку не получится
	}

	// сделаем заглушку позиции в случае разрыва
	p.P1.X = p.N1.X + Dx*p.L1/D
	p.P1.Y = p.N1.Y + Dy*p.L1/D
	p.P2.X = p.N2.X - Dx*p.L2/D
	p.P2.Y = p.N2.Y - Dy*p.L2/D

	// теперь вычислим координаты инструментального узла
	// как точки пересечения двух окружностей

	R12 := p.L1 * p.L1
	R22 := p.L2 * p.L2
	A := (R12 - R22 + D2) / (2 * D)
	A2 := A * A
	H2 := R12 - A2
	if H2 < 0 {
		return
	}
	H := math.Sqrt(H2)
	P2X := p.N1.X + A*Dx/D
	P2Y := p.N1.Y + A*Dy/D
	p.P1.X = P2X + (H*Dy)/D
	p.P1.Y = P2Y - (H*Dx)/D
	p.P2.X = P2X - (H*Dy)/D
	p.P2.Y = P2Y + (H*Dx)/D
	p.Valid = true
}
