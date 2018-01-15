package crazy

import (
	"encoding/binary"
	"math"
)

// A Distribution adapts a source to produce floating-point numbers.
type Distribution interface {
	Next() float64
}

// Uniform1_2 produces numbers in the interval [1, 2). This interval is chosen
// for speed.
type Uniform1_2 struct {
	Source
}

// Next produces a uniform variate in the interval [1, 2).
func (u Uniform1_2) Next() float64 {
	b := [8]byte{}
	u.Read(b[:7])
	return math.Float64frombits(binary.LittleEndian.Uint64(b[:]) | 0x3ff0000000000000)
}

// Uniform produces numbers in the interval [Low, High).
type Uniform struct {
	Source
	Low, High float64
}

// Next produces a uniform variate.
func (u Uniform) Next() float64 {
	return u.High*(Uniform1_2{u.Source}.Next()-1) + u.Low
}
