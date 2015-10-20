package math

import (
	"math"
)

type Vec2 struct {
	X float32
	Y float32
}

type Mat3 [9]float32


func (v *Vec2) Add(b Vec2) Vec2 {
	return Vec2{v.X + b.X, v.Y + b.Y}
}

func (v *Vec2) Subtract(b Vec2) Vec2 {
	return Vec2{v.X - b.X, v.Y - b.Y}
}

func (v *Vec2) Length() float32 {
	return float32(math.Sqrt(float64(v.X * v.X + v.Y * v.Y)))
}

func (v *Vec2) Mul(m Mat3) Vec2 {
	return Vec2{m[0] * v.X + m[1] * v.Y + m[2], 
							m[3] * v.X + m[4] * v.Y + m[5]}
}

func NewMat3(a,b,c,d,e,f,g,h,i float32) Mat3 {
	return Mat3{a,b,c,d,e,f,g,h,i}
}

func (m *Mat3) Mul(n Mat3) Mat3 {
	return Mat3{m[0]*n[0], m[1]*n[3], m[2]*n[6],
							m[3]*n[1], m[4]*n[4], m[5]*n[7],
							m[6]*n[2], m[7]*n[5], m[8]*n[8]}
}

func Identity() Mat3 {
	return Mat3{1,0,0,
		          0,1,0,
		          0,0,1}
}

func Translate(x,y float32) Mat3 {
	return Mat3{1,0,x,
		          0,1,y,
		          0,0,1}
} 

func TranslateFromVec2(v Vec2) Mat3 {
	return Translate(v.X, v.Y)
}

func Scale(x,y float32) Mat3 {
	return Mat3{x,0,0,
							0,y,0,
						  0,0,1}
}

func ScaleFromVec2(v Vec2) Mat3 {
	return Scale(v.X, v.Y)
}

func Rotate(x float32) Mat3 {
	s := float32(math.Sin(float64(x)))
	c := float32(math.Cos(float64(x)))
	return Mat3{c,-s, 0,
							s, c, 0,
						  0, 0, 1}
}

type Size struct {
	W float32
	H float32
}

func (s *Size) ToVec2() Vec2 {
	return Vec2{s.W, s.H}
}

func Distance(a Vec2, b Vec2) float32 {
	v := a.Subtract(b)
	return v.Length()
}