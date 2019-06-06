// +build go1.9

package crazy

import (
	"encoding/binary"
	"io"
	"math/bits"
)

// Xoshiro implements the xoshiro256** PRNG created by David Blackman and
// Sebastiano Vigna. It has period 2**256-1 with 256 state bits and exhibits
// 4-equidistribution (all tuples of four consecutive 64-bit values except the
// all-zero tuple appear once in the sequence), with irrelevantly large linear
// complexity in all bits.
//
// Compared to MT64-19937, xoshiro256** is much faster and smaller, but its
// period is also much smaller, and it has a much lower dimensional
// distribution. xoshiro recovers from "poor" states quickly.
type Xoshiro struct {
	w, x, y, z uint64
}

// NewXoshiro produces an unseeded Xoshiro. Call Seed[IV]() or Restore() prior
// to use.
func NewXoshiro() *Xoshiro {
	return &Xoshiro{}
}

// SeedIV initializes the generator using all bits of iv, which may be of any
// size or nil.
func (xoshi *Xoshiro) SeedIV(iv []byte) {
	// Seed using Sebastiano Vigna's SplitMix64, as recommended by the
	// xoroshiro128+ source. However, that is the recommendation "if you have
	// a 64-bit seed," so we need a strategy to use more bits than that in a
	// way that will actually increase entropy up to the state size. What we
	// shall do is initialize the generator to an SM64-generated state, then
	// for each 256 bits of the iv, produce four variates from the generator,
	// add those 256 bits, then add four more iterations of the SM64.
	var sm uint64
	sm += 0x9e3779b97f4a7c15
	z := sm
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	z ^= z >> 31
	xoshi.w = z
	sm += 0x9e3779b97f4a7c15
	z = sm
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	z ^= z >> 31
	xoshi.x = z
	sm += 0x9e3779b97f4a7c15
	z = sm
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	z ^= z >> 31
	xoshi.y = z
	sm += 0x9e3779b97f4a7c15
	z = sm
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	z ^= z >> 31
	xoshi.z = z
	for len(iv) >= 32 {
		xoshi.Uint64()
		xoshi.Uint64()
		xoshi.Uint64()
		xoshi.Uint64()
		xoshi.w ^= binary.LittleEndian.Uint64(iv)
		iv = iv[8:]
		xoshi.x ^= binary.LittleEndian.Uint64(iv)
		iv = iv[8:]
		xoshi.y ^= binary.LittleEndian.Uint64(iv)
		iv = iv[8:]
		xoshi.z ^= binary.LittleEndian.Uint64(iv)
		iv = iv[8:]
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.w ^= z
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.x ^= z
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.y ^= z
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.z ^= z
	}
	if len(iv) > 24 {
		xoshi.Uint64()
		xoshi.Uint64()
		xoshi.Uint64()
		xoshi.Uint64()
		xoshi.w ^= binary.LittleEndian.Uint64(iv)
		iv = iv[8:]
		xoshi.x ^= binary.LittleEndian.Uint64(iv)
		iv = iv[8:]
		xoshi.y ^= binary.LittleEndian.Uint64(iv)
		iv = iv[8:]
		p := []byte{7: 0}
		copy(p, iv)
		xoshi.z ^= binary.LittleEndian.Uint64(p)
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.w ^= z
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.x ^= z
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.y ^= z
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.z ^= z
	} else if len(iv) > 16 {
		xoshi.Uint64()
		xoshi.Uint64()
		xoshi.Uint64()
		xoshi.w ^= binary.LittleEndian.Uint64(iv)
		iv = iv[8:]
		xoshi.x ^= binary.LittleEndian.Uint64(iv)
		iv = iv[8:]
		p := []byte{7: 0}
		copy(p, iv)
		xoshi.y ^= binary.LittleEndian.Uint64(p)
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.w ^= z
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.x ^= z
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.y ^= z
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.z ^= z
	} else if len(iv) > 8 {
		xoshi.Uint64()
		xoshi.Uint64()
		xoshi.w ^= binary.LittleEndian.Uint64(iv)
		iv = iv[8:]
		p := []byte{7: 0}
		copy(p, iv)
		xoshi.x ^= binary.LittleEndian.Uint64(p)
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.w ^= z
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.x ^= z
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.y ^= z
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.z ^= z
	} else if len(iv) > 0 {
		xoshi.Uint64()
		p := []byte{7: 0}
		copy(p, iv)
		xoshi.w ^= binary.LittleEndian.Uint64(p)
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.w ^= z
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.x ^= z
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.y ^= z
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		xoshi.z ^= z
	}
}

// Uint64 produces a 64-bit pseudo-random value.
func (xoshi *Xoshiro) Uint64() uint64 {
	r := bits.RotateLeft64(xoshi.x*5, 7) * 9
	t := xoshi.x << 17
	xoshi.y ^= xoshi.w
	xoshi.z ^= xoshi.x
	xoshi.x ^= xoshi.y
	xoshi.w ^= xoshi.z
	xoshi.y ^= t
	xoshi.z = bits.RotateLeft64(xoshi.z, 45)
	return r
}

// Read fills p with random bytes generated 64 bits at a time, discarding
// unused bytes. n will always be len(p) and err will always be nil.
func (xoshi *Xoshiro) Read(p []byte) (n int, err error) {
	n = len(p)
	for len(p) > 8 {
		binary.LittleEndian.PutUint64(p, xoshi.Uint64())
		p = p[8:]
	}
	b := [8]byte{}
	binary.LittleEndian.PutUint64(b[:], xoshi.Uint64())
	copy(p, b[:])
	return n, nil
}

// Save serializes the current state of the xoshiro256** generator. Values
// produced by such a generator that has Restore()d this state are guaranteed
// to match those produced by this exact generator. n should always be 32
// bytes.
func (xoshi *Xoshiro) Save(into io.Writer) (n int, err error) {
	p := []byte{31: 0}
	binary.LittleEndian.PutUint64(p, xoshi.w)
	binary.LittleEndian.PutUint64(p[8:], xoshi.x)
	binary.LittleEndian.PutUint64(p[16:], xoshi.y)
	binary.LittleEndian.PutUint64(p[24:], xoshi.z)
	return into.Write(p)
}

// Restore loads a Save()d xoshiro256** state.
func (xoshi *Xoshiro) Restore(from io.Reader) (n int, err error) {
	p := []byte{31: 0}
	if n, err = from.Read(p); n < len(p) {
		return n, err
	}
	xoshi.w = binary.LittleEndian.Uint64(p)
	xoshi.x = binary.LittleEndian.Uint64(p[8:])
	xoshi.y = binary.LittleEndian.Uint64(p[16:])
	xoshi.z = binary.LittleEndian.Uint64(p[24:])
	return n, nil
}

// Seed is a proxy to SeedInt64. This exists to satisfy the rand.Source
// interface.
func (xoshi *Xoshiro) Seed(x int64) {
	SeedInt64(xoshi, x)
}
