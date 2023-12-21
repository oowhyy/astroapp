package vector

import (
	"fmt"
	"math"
)

type Vector struct {
	X, Y float64
}

func FromFloats(x, y float64) Vector {
	return Vector{x, y}
}

func FromInts(x, y int) Vector {
	return Vector{float64(x), float64(y)}
}

func Sum(v1, v2 Vector) Vector {
	return Vector{v1.X + v2.X, v1.Y + v2.Y}
}

func Diff(v1, v2 Vector) Vector {
	return Vector{v1.X - v2.X, v1.Y - v2.Y}
}

func Scaled(v Vector, x float64) Vector {
	return Vector{v.X * x, v.Y * x}
}

func (v *Vector) Sub(v2 Vector) {
	v.X -= v2.X
	v.Y -= v2.Y
}

func (v *Vector) Add(v2 Vector) {
	v.X += v2.X
	v.Y += v2.Y
}

func (v *Vector) Scale(x float64) {
	v.X *= x
	v.Y *= x
}

func (v Vector) Len() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vector) Reset() {
	v.X = 0
	v.Y = 0
}

func (v Vector) String() string {
	return fmt.Sprintf("(%f, %f)", v.X, v.Y)
}
