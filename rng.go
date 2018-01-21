package crazy

import (
	"encoding/binary"
	"math/big"
)

// RNG adapts a Source to produce integers.
type RNG struct {
	Source
}

// Uint64 generates a random uint64.
func (r RNG) Uint64() uint64 {
	p := [8]byte{}
	r.Read(p[:])
	return binary.LittleEndian.Uint64(p[:])
}

// Uint32 generates a random uint32.
func (r RNG) Uint32() uint32 {
	p := [4]byte{}
	r.Read(p[:])
	return binary.LittleEndian.Uint32(p[:])
}

// Uint16 generates a random uint16.
func (r RNG) Uint16() uint16 {
	p := [2]byte{}
	r.Read(p[:])
	return binary.LittleEndian.Uint16(p[:])
}

// Uintn generates a random uint in the interval [0, max). It panics if
// max == 0.
func (r RNG) Uintn(max uint) uint {
	bad := ^uint(0) - ^uint(0)%max
	x := uint(r.Uint64())
	for x > bad {
		x = uint(r.Uint64())
	}
	return x % max
}

// Intn generates a random int in the interval [0, max). It panics if max <= 0.
func (r RNG) Intn(max int) int {
	if max < 0 {
		panic("maximum below zero")
	}
	return int(r.Uintn(uint(max)))
}

// Big generates a random number with maximum bit length nbits.
func (r RNG) Big(nbits int) *big.Int {
	p := make([]byte, (uint(nbits)+7)>>3)
	r.Read(p)
	p[0] &= byte(0xff >> (8 - uint(nbits)&7))
	return new(big.Int).SetBytes(p)
}

// Bign generates a random number in the range [0, max). It panics if max <= 0.
func (r RNG) Bign(max *big.Int) *big.Int {
	if max.Sign() <= 0 {
		panic("maximum zero or below")
	}
	n := make([]big.Word, len(max.Bits()))
	for i := range n[1:] {
		n[i+1] = ^big.Word(0)
	}
	nbits := max.BitLen() + 1
	m := new(big.Int)
	m.Sub(m.SetBit(m, nbits, 1), big.NewInt(1))
	k := new(big.Int)
	k.Rem(m, max)
	k.Sub(m, k)
	x := r.Big(nbits)
	for x.Cmp(k) > 0 {
		x = r.Big(nbits)
	}
	return x.Rem(x, max)
}
