package crazy

import (
	"encoding/binary"
	"io"
)

const (
	mt64N        = 312
	mt64M        = 156
	mt64A uint64 = 0xB5026F5AA96619E9
)

// MT64 is a 64-bit Mersenne twister. The choice of parameters yields a period
// of 2**19937 - 1.
type MT64 struct {
	i int
	s [mt64N]uint64
}

// NewMT64 produces an unseeded 64-bit Mersenne twister. Call mt.Seed[IV]() or
// mt.Restore() prior to use.
func NewMT64() *MT64 {
	return &MT64{}
}

// Seed initializes mt using an int64. This serves to satisfy the rand.Source
// interface.
func (mt *MT64) Seed(seed int64) {
	mt.s[0] = uint64(seed)
	for i := 1; i < mt64N; i++ {
		mt.s[i] = 6364136223846793005*(mt.s[i-1]^mt.s[i-1]>>62) + uint64(i)
	}
	mt.i = mt64N
}

// SeedIV initializes the generator using all bits of iv, which may be of any
// size or nil.
func (mt *MT64) SeedIV(iv []byte) {
	mt.Seed(19650218)
	if len(iv) == 0 {
		return
	}
	if len(iv)&7 != 0 {
		// We need a multiple of 8, but we don't really own iv, so we'll make a
		// copy with the right length.
		t := make([]byte, len(iv)+8-len(iv)&7)
		copy(t, iv)
		iv = t
	}
	k := mt64N
	if mt64N < len(iv) {
		k = len(iv)
	}
	i := 1
	for j := 0; k > 0; k-- {
		mt.s[i] = (mt.s[i]^mt.s[i-1]^mt.s[i-1]>>62)*3935559000370003845 + binary.LittleEndian.Uint64(iv[j*8:]) + uint64(j)
		if i++; i >= mt64N {
			mt.s[0] = mt.s[mt64N-1]
			i = 1
		}
		if j++; j*8 >= len(iv) {
			j = 0
		}
	}
	for k = mt64N - 1; k > 0; k-- {
		mt.s[i] = (mt.s[i]^mt.s[i-1]^mt.s[i-1]>>62)*2862933555777941757 - uint64(i)
		if i++; i >= mt64N {
			mt.s[0] = mt.s[mt64N-1]
			i = 1
		}
	}
}

// Uint64 produces a 64-bit pseudo-random value. This primarily serves to
// satisfy the rand.Source64 interface, but it also provides direct access to
// the algorithm's values, which can simplify usage in some scenarios.
func (mt *MT64) Uint64() uint64 {
	if mt.i >= mt64N {
		i := 0
		for i < mt64N-mt64M {
			x := mt.s[i]&0xffffffff80000000 | mt.s[i+1]&0x000000007fffffff
			x = x>>1 ^ mt64A*(x&1)
			mt.s[i] = mt.s[i+mt64M] ^ x
			i++
		}
		for i < mt64N-1 {
			x := mt.s[i]&0xffffffff80000000 | mt.s[i+1]&0x000000007fffffff
			x = x>>1 ^ mt64A*(x&1)
			mt.s[i] = mt.s[i-(mt64N-mt64M)] ^ x
			i++
		}
		x := mt.s[mt64N-1]&0xffffffff80000000 | mt.s[0]&0x000000007fffffff
		x = x>>1 ^ mt64A*(x&1)
		mt.s[mt64N-1] = mt.s[mt64M-1] ^ x
		mt.i = 0
	}

	x := mt.s[mt.i]
	mt.i++
	x ^= x >> 29 & 0x5555555555555555
	x ^= x << 17 & 0x71D67FFFEDA60000
	x ^= x << 37 & 0xFFF7EEE000000000
	x ^= x >> 43
	return x
}

// Read fills p with random bytes generated 64 bits at a time, discarding
// unused bytes. n will always be len(p) and err will always be nil.
func (mt *MT64) Read(p []byte) (n int, err error) {
	n = len(p)
	for len(p) >= 8 {
		binary.LittleEndian.PutUint64(p, mt.Uint64())
		p = p[8:]
	}
	b := [8]byte{}
	binary.LittleEndian.PutUint64(b[:], mt.Uint64())
	copy(p, b[:])
	return n, nil
}

// Int63 generates an integer in the interval [0, 2**63 - 1]. This serves to
// satisfy the rand.Source interface.
func (mt *MT64) Int63() int64 {
	return int64(mt.Uint64() >> 1)
}

// Save serializes the current state of the Mersenne twister. Values produced
// by an MT that has Restore()d this state are guaranteed to match those
// produced by this exact generator. n should always be 2 + N*8 = 2498 bytes.
func (mt *MT64) Save(into io.Writer) (n int, err error) {
	// Unlike with LFG, we can't* get away without saving mt.i, so we have to
	// spend an extra two bytes. That said, the state size of this is actually
	// much smaller than that of LFG.
	// *It is possible to produce a rotation of the MT state, but much more
	// expensive, and I don't feel like figuring out how to do it anyway.
	p := [2 + mt64N*8]byte{}
	for i, v := range mt.s {
		binary.LittleEndian.PutUint64(p[i<<3:], v)
	}
	// Writing mt.i last allows us to stay aligned while writing mt.s.
	p[len(p)-2] = byte(mt.i >> 8)
	p[len(p)-1] = byte(mt.i)
	return into.Write(p[:])
}

// Restore loads a Save()d MT state. This reads 2 + N*8 = 2498 bytes as the
// feed and state values.
func (mt *MT64) Restore(from io.Reader) (n int, err error) {
	p := [2 + mt64N*8]byte{}
	if n, err = from.Read(p[:]); err != nil {
		return n, err
	}
	for i := range mt.s {
		mt.s[i] = binary.LittleEndian.Uint64(p[i<<3:])
	}
	mt.i = int(p[len(p)-2])<<8 | int(p[len(p)-1])
	return n, nil
}
