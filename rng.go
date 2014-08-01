package crazy

import (
	"encoding/binary"
	"math/big"
)

// Adapt a Source to produce integers.
type RNG struct {
	Source
}

// Generate a random uint64. All bits are used.
func (r RNG) Uint64() uint64 {
	p := [8]byte{}
	r.Read(p[:])
	return binary.LittleEndian.Uint64(p[:])
}

// Generate a random uint32. All bits are used.
func (r RNG) Uint32() uint32 {
	p := [4]byte{}
	r.Read(p[:])
	return binary.LittleEndian.Uint32(p[:])
}

// Generate a random uint16. All bits are used.
func (r RNG) Uint16() uint16 {
	p := [2]byte{}
	r.Read(p[:])
	return binary.LittleEndian.Uint16(p[:])
}

// Generate a random uint in the interval [0, max).
func (r RNG) Uintn(max uint) uint {
	bad := ^uint(0) - ^uint(0)%max
	x := uint(r.Uint64())
	for x > bad {
		x = uint(r.Uint64())
	}
	return x % max
}

// Generate a random int in the interval [0, max). Panic if max < 0.
func (r RNG) Intn(max int) int {
	if max < 0 {
		panic("maximum below zero")
	}
	return int(r.Uintn(uint(max)))
}

// Generate a random number with bit length nbits.
func (r RNG) Big(nbits int) *big.Int {
	p := make([]byte, uint(nbits)>>3+uint(nbits|nbits>>1|nbits>>2)&1)
	r.Read(p)
	p[0] &= byte(0xff >> (8 - uint(nbits)&7))
	return new(big.Int).SetBytes(p)
}

// Generate a random number in the range [0, max). Panic if max <= 0.
func (r RNG) Bign(max *big.Int) *big.Int {
	if max.Sign() <= 0 {
		panic("maximum zero or below")
	}
	n := make([]big.Word, len(max.Bits()))
	for i := range n[1:] {
		n[i+1] = ^big.Word(0)
	}
	nbits := max.BitLen() + 1
	// if uint64(^big.Word(0))>>32 != 0 {
	// 	if nbits&63 != 0 {
	// 		n[len(n)-1] = ^big.Word(0) >> (64 - big.Word(nbits)&63)
	// 	} else {
	// 		n[len(n)-1] = ^big.Word(0)
	// 	}
	// } else {
	// 	if nbits&31 != 0 {
	// 		n[len(n)-1] = ^big.Word(0) >> (32 - big.Word(nbits)&31)
	// 	} else {
	// 		n[len(n)-1] = ^big.Word(0)
	// 	}
	// }
	// m := new(big.Int).SetBits(n)
	m := new(big.Int)
	m.Sub(m.SetBit(m, nbits, 1), big.NewInt(1))
	// print("m: ")
	// println(m.String())
	k := new(big.Int)
	// println(max.String())
	k.Rem(m, max)
	// println(k.String())
	k.Sub(m, k)
	// print("bad: ")
	// println(k.String())
	x := r.Big(nbits)
	// print("x: ")
	// println(x.String())
	for x.Cmp(k) > 0 {
		x = r.Big(nbits)
		// print("x: ")
		// println(x.String())
	}
	return x.Rem(x, max)
}
