package crazy

import (
	"encoding/binary"
	"math"
)

// Adapt a source to produce floating-point numbers.
type Distribution interface {
	Next() float64
}

// Produce numbers in the interval [1, 2).
type Uniform1_2 struct {
	Source
}

func (u Uniform1_2) Next() float64 {
	b := [8]byte{}
	u.Read(b[:7])
	return math.Float64frombits(binary.LittleEndian.Uint64(b[:]) | 0x3fff000000000000)
}

// Produce numbers in the interval [0, 1).
type Uniform struct {
	Source
	Low, High float64
}

func (u Uniform) Next() float64 {
	return u.High*(Uniform1_2{u.Source}.Next()-1) + u.Low
}
