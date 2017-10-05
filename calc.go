package main

import (
	"math"
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


func area(A, B, C Position) float64 {
	return (B.X - A.X) * (C.Y - A.Y) - (B.Y - A.Y) * (C.X - A.X)
}

func intersect1(A, B, C, D float64) bool {
	if B < A { A, B = B, A }
	if D < C { C, D = D, C }
	if C < A { C = A } // C = max(A,C)
	if D < B { B = D } // B = min(B,D)
	return (C < B)
}

// проверяем отрезки AB и CD на наличие пересечения
func HasIntersect(A, B, C, D Position) bool {
	return intersect1(A.X, B.X, C.X, D.X) &&
		intersect1(A.Y, B.Y, C.Y, D.Y) &&
		area(A,B,C) * area(A,B,D) <= 0 &&
		area(C,D,A) * area(C,D,B) <= 0
}

// ищем пересечение двух окружностей
func CircleCross(N1 Position, R1 float64, N2 Position, R2 float64) (P1, P2 Position, Ok bool) {
	Ok = false
	Dx := N2.X - N1.X
	Dy := N2.Y - N1.Y

	D2 := Dx*Dx + Dy*Dy
	D := math.Sqrt(D2)

	if D == 0 {
		return // совпадение осей. вычислить точку не получится
	}

	// сделаем заглушку позиции в случае разрыва
	P1.X = N1.X + Dx*R1/D
	P1.Y = N1.Y + Dy*R1/D
	P2.X = N2.X - Dx*R2/D
	P2.Y = N2.Y - Dy*R2/D

	// теперь вычислим координаты инструментального узла
	// как точки пересечения двух окружностей

	R12 := R1 * R1
	R22 := R2 * R2
	A := (R12 - R22 + D2) / (2 * D)
	A2 := A * A
	H2 := R12 - A2
	if H2 < 0 {
		return
	}
	H := math.Sqrt(H2)
	P2X := N1.X + A*Dx/D
	P2Y := N1.Y + A*Dy/D
	P1.X = P2X + (H*Dy)/D
	P1.Y = P2Y - (H*Dx)/D
	P2.X = P2X - (H*Dy)/D
	P2.Y = P2Y + (H*Dx)/D
	Ok = true
	return
}

func (p *Pantograph) Solve() {
	p.Valid = false // пока у нас не доказано, что всё хорошо - всё плохо.

	// вычислим координаты концов ведущих тяг
	Sa1, Ca1 := math.Sincos(p.E1.A + p.A1)
	Sa2, Ca2 := math.Sincos(p.E2.A + p.A2)
	p.N1.X = p.E1.X + p.L1*Ca1
	p.N1.Y = p.E1.Y - p.L1*Sa1
	p.N2.X = p.E2.X + p.L2*Ca2
	p.N2.Y = p.E2.Y - p.L2*Sa2

	p.P1, p.P2, p.Valid = CircleCross(p.N1, p.L1, p.N2, p.L2)
	p.Valid = p.Valid && !(
		HasIntersect(p.E1.Position, p.N1, p.E2.Position, p.N2) ||
		HasIntersect(p.E1.Position, p.N1, p.N2, p.P1) ||
		HasIntersect(p.E2.Position, p.N2, p.N1, p.P1))
}
