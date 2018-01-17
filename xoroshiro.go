// +build go1.9

package crazy

import (
	"encoding/binary"
	"io"
	"math/bits"
)

// Xoroshiro implements the Xoroshiro128+ PRNG created by David Blackman and
// Sebastiano Vigna. It has period 2**128-1 with 128 state bits. Xoroshiro is
// the recommended PRNG for most non-scientific applications, as it is
// extremely fast and has a very small state while still having very good
// statistical randomness.
type Xoroshiro [2]uint64

// NewXoroshiro produces an unseeded Xoroshiro. Call either xoro.Seed[IV]() or
// xoro.Restore() prior to use.
func NewXoroshiro() *Xoroshiro {
	return &Xoroshiro{}
}

// SeedIV initializes the generator using all bits of iv, which may be of any
// size or nil.
func (xoro *Xoroshiro) SeedIV(iv []byte) {
	// Seed using Sebastiano Vigna's SplitMix64, as recommended by the
	// xoroshiro128+ source. However, that is the recommendation "if you have
	// a 64-bit seed," so we need a strategy to use more bits than that in a
	// way that will actually increase entropy up to the state size. What we
	// shall do is add (GF2) each 64-bit integer from the iv into the SplitMix
	// state between each advancement of it until we have no more than 128 bits
	// remaining in the iv, then follow the same procedure for the values of
	// the xoroshiro state. If we're given an empty (nil) iv, however, just use
	// SplitMix with 0 seed.
	var sm uint64
	if len(iv) == 0 {
		sm += 0x9e3779b97f4a7c15
		z := sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		(*xoro)[0] = z
		sm += 0x9e3779b97f4a7c15
		z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		(*xoro)[1] = z
		return
	}
	for len(iv) > 8 {
		p := []byte{7: 0}
		iv = iv[copy(p, iv):]
		sm += 0x9e3779b97f4a7c15 ^ binary.LittleEndian.Uint64(p)
	}
	if len(iv) > 0 {
		p := []byte{7: 0}
		iv = iv[copy(p, iv):]
		sm += 0x9e3779b97f4a7c15
		z := sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		(*xoro)[1] = z
		sm ^= binary.LittleEndian.Uint64(p)
	} else {
		sm += 0x9e3779b97f4a7c15
		z := sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		(*xoro)[1] = z
	}
	p := []byte{7: 0}
	copy(p, iv)
	sm += 0x9e3779b97f4a7c15
	z := sm
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	z ^= z >> 31
	(*xoro)[0] = z
}

func (xoro *Xoroshiro) next() uint64 {
	s0, s1 := (*xoro)[0], (*xoro)[1]
	x := s0 + s1
	s1 ^= s0
	(*xoro)[0] = bits.RotateLeft64(s0, 55) ^ s1 ^ s1<<14
	(*xoro)[1] = bits.RotateLeft64(s1, 36)
	return x
}

// Read fills p with random bytes generated 64 bits at a time, discarding
// unused bytes. n will always be len(p) and err will always be nil.
func (xoro *Xoroshiro) Read(p []byte) (n int, err error) {
	n = len(p)
	for len(p) >= 8 {
		binary.LittleEndian.PutUint64(p, xoro.next())
		p = p[8:]
	}
	if len(p) > 8 {
		b := [8]byte{}
		binary.LittleEndian.PutUint64(b[:], xoro.next())
		copy(p, b[:])
	}
	return n, nil
}

// Save serializes the current state of the xoroshiro128+ generator. Values
// produced by such a generator that has Restore()d this state are guaranteed
// to match those produced by this exact generator. n should always be 16
// bytes.
func (xoro *Xoroshiro) Save(into io.Writer) (n int, err error) {
	p := []byte{15: 0}
	binary.LittleEndian.PutUint64(p, (*xoro)[0])
	binary.LittleEndian.PutUint64(p[8:], (*xoro)[1])
	return into.Write(p)
}

// Restore loads a Save()d xoroshiro128+ state.
func (xoro *Xoroshiro) Restore(from io.Reader) (n int, err error) {
	p := []byte{15: 0}
	if n, err = from.Read(p); n < len(p) {
		return n, err
	}
	(*xoro)[0] = binary.LittleEndian.Uint64(p)
	(*xoro)[1] = binary.LittleEndian.Uint64(p[8:])
	return n, nil
}

// Seed is a proxy to SeedInt64. This serves to satisfy the rand.Source
// interface.
func (xoro *Xoroshiro) Seed(x int64) {
	SeedInt64(xoro, x)
}

// Int63 generates an integer in the interval [0, 2**63 - 1]. This serves to
// satisfy the rand.Source interface.
func (xoro *Xoroshiro) Int63() int64 {
	return int64(xoro.next() >> 1)
}

// Rexoroshiro is the same as Xoroshiro but yields values that are bytewise
// reversed. This sacrifices some speed to place the bits with greater apparent
// randomness in the low positions, making the generator stronger when
// producing values interpreted in small ranges; e.g., RNG{rexo}.Intn(2) is
// typically "better" than RNG{xoro}.Intn(2).
type Rexoroshiro Xoroshiro

// NewRexoroshiro creates an uninitialized Rexoroshiro. This is equivalent to
// casting an uninitialized Xoroshiro to Rexoroshiro and has the same caveats.
func NewRexoroshiro() *Rexoroshiro {
	return (*Rexoroshiro)(NewXoroshiro())
}

// SeedIV initializes the generator as if it were a Xoroshiro.
func (rexo *Rexoroshiro) SeedIV(iv []byte) {
	(*Xoroshiro).SeedIV((*Xoroshiro)(rexo), iv)
}

func (rexo *Rexoroshiro) next() uint64 {
	return bits.ReverseBytes64((*Xoroshiro).next((*Xoroshiro)(rexo)))
}

// Read fills p with random bytes that are byte-reversed Xoroshiro values.
func (rexo *Rexoroshiro) Read(p []byte) (n int, err error) {
	n = len(p)
	for len(p) >= 8 {
		binary.LittleEndian.PutUint64(p, rexo.next())
		p = p[8:]
	}
	if len(p) > 0 {
		b := [8]byte{}
		binary.LittleEndian.PutUint64(b[:], rexo.next())
		copy(p, b[:])
	}
	return n, nil
}

// Save serializes the Rexoroshiro state in the same way as the equivalent
// Xoroshiro state.
func (rexo *Rexoroshiro) Save(into io.Writer) (n int, err error) {
	return (*Xoroshiro).Save((*Xoroshiro)(rexo), into)
}

// Restore loads the Rexoroshiro state in the same way as the equivalent
// Xoroshiro state.
func (rexo *Rexoroshiro) Restore(from io.Reader) (n int, err error) {
	return (*Xoroshiro).Restore((*Xoroshiro)(rexo), from)
}

// Seed is a proxy to Xoroshiro's Seed. This serves to satisfy the rand.Source
// interface.
func (rexo *Rexoroshiro) Seed(x int64) {
	(*Xoroshiro).Seed((*Xoroshiro)(rexo), x)
}

// Int63 generates an integer in the interval [0, 2**63 - 1]. This serves to
// satisfy the rand.Source interface.
func (rexo *Rexoroshiro) Int63() int64 {
	// Since the low bits have the highest quality, it makes sense to mask
	// instead of shifting.
	return int64(rexo.next() & 0x7fffffffffffffff)
}
