package main

import (
	"math"
)

// работа с координатами //////////////////////////////////////

func (pt Position) Length2() float64 {
	return pt.X*pt.X + pt.Y*pt.Y
}

func (pt Position) Length() float64 {
	return math.Sqrt(pt.Length2())
}

func Delta(A, B Position) Position {
	return Position{X: B.X - A.X, Y: B.Y - A.Y}
}

func Distance(A, B Position) float64 {
	return Delta(A, B).Length()
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
	D := Delta(N1, N2)
	dist2 := D.Length2()
	dist := math.Sqrt(dist2)

	if dist == 0 {
		return // совпадение осей. вычислить точку не получится
	}

	// сделаем заглушку позиции в случае разрыва
	P1.X = N1.X + D.X*R1/dist
	P1.Y = N1.Y + D.Y*R1/dist
	P2.X = N2.X - D.X*R2/dist
	P2.Y = N2.Y - D.Y*R2/dist

	// теперь вычислим координаты инструментального узла
	// как точки пересечения двух окружностей

	R12 := R1 * R1
	R22 := R2 * R2
	A := (R12 - R22 + dist2) / (2 * dist)
	A2 := A * A
	H2 := R12 - A2
	if H2 < 0 {
		return
	}
	H := math.Sqrt(H2)
	P2X := N1.X + A*D.X/dist
	P2Y := N1.Y + A*D.Y/dist
	P1.X = P2X + (H*D.Y)/dist
	P1.Y = P2Y - (H*D.X)/dist
	P2.X = P2X - (H*D.Y)/dist
	P2.Y = P2Y + (H*D.X)/dist
	Ok = true
	return
}

// решение пантографа /////////////////////////////////////////

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
