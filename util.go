package crazy

import (
	"crypto/rand"
	"encoding/binary"
	mathrand "math/rand"
)

// XorCompose composes two Sources using XOR.
type XorCompose struct {
	A, B Source
}

// Read produces random bits from the composition of the two sources: A ^ B.
func (x XorCompose) Read(p []byte) (int, error) {
	buf := make([]byte, len(p))
	x.A.Read(p)
	x.B.Read(buf)
	for i, v := range buf {
		p[i] ^= v
	}
	return len(p), nil
}

// SelectCompose composes three sources such that each bit of S chooses which
// of A and B has its corresponding bit chosen.
type SelectCompose struct {
	S, A, B Source
}

// Read produces random bits selected randomly between A and B.
func (x SelectCompose) Read(p []byte) (int, error) {
	s, b := make([]byte, len(p)), make([]byte, len(p))
	x.S.Read(s)
	x.A.Read(p)
	x.B.Read(b)
	for i := range p {
		p[i] ^= (p[i] ^ b[i]) & s[i]
	}
	return len(p), nil
}

type (
	math2crazy struct {
		mathrand.Source
	}

	math642crazy struct {
		mathrand.Source64
	}
)

// AdaptRand turns a math/rand Source into a crazy Seeder. If the argument
// already implements Seeder, it is returned directly. Otherwise, if the
// argument does not implement math/rand Source64, then to preserve
// equidistribution, the Read() method of this Seeder uses two Int63() calls
// per eight bytes of output, using one fewer call when the output is not
// divisible by 8. Either way, the SeedIV() method uses only up to the first
// eight bytes of the argument.
func AdaptRand(src mathrand.Source) Seeder {
	if s, ok := src.(Seeder); ok {
		return s
	}
	if s, ok := src.(mathrand.Source64); ok {
		return math642crazy{s}
	}
	return math2crazy{src}
}

// AdaptRand64 is like AdaptRand but takes a math/rand Source64 instead.
func AdaptRand64(src mathrand.Source64) Seeder {
	if s, ok := src.(Seeder); ok {
		return s
	}
	return math642crazy{src}
}

func (m math2crazy) Read(p []byte) (n int, err error) {
	n = len(p)
	for len(p) >= 8 {
		binary.LittleEndian.PutUint64(p, uint64(m.Int63()^m.Int63()<<1))
		p = p[8:]
	}
	if len(p) > 0 {
		b := [8]byte{}
		binary.LittleEndian.PutUint64(b[:], uint64(m.Int63()))
		copy(p, b[:])
	}
	return n, nil
}

func (m math642crazy) Read(p []byte) (n int, err error) {
	n = len(p)
	for len(p) >= 8 {
		binary.LittleEndian.PutUint64(p, m.Uint64())
		p = p[8:]
	}
	if len(p) > 0 {
		b := [8]byte{}
		binary.LittleEndian.PutUint64(b[:], m.Uint64())
		copy(p, b[:])
	}
	return n, nil
}

func (m math2crazy) SeedIV(iv []byte) {
	b := [8]byte{}
	copy(b[:], iv)
	m.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

func (m math642crazy) SeedIV(iv []byte) {
	b := [8]byte{}
	copy(b[:], iv)
	m.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

type crazy2math struct {
	Seeder
}

// AdaptCrazy turns a crazy Seeder into a math/rand Source. If the argument
// already implements math/rand.Source64, it is returned directly.
func AdaptCrazy(src Seeder) mathrand.Source64 {
	if m, ok := src.(mathrand.Source64); ok {
		return m
	}
	return crazy2math{src}
}

func (c crazy2math) Uint64() uint64 {
	return RNG{c.Seeder}.Uint64()
}

func (c crazy2math) Int63() int64 {
	return int64(RNG{c.Seeder}.Uint64() >> 1)
}

func (c crazy2math) Seed(x int64) {
	SeedInt64(c.Seeder, x)
}

// RandString creates a random string with specified length using only
// characters from the given alphabet. The mechanism for constructing the
// string is choosing random positions from the alphabet, so characters which
// appear more often have higher relative probability.
func RandString(rng RNG, alphabet string, length int) string {
	p := make([]rune, length)
	a := []rune(alphabet)
	for i := range p {
		p[i] = a[int(rng.Uintn(uint(len(a))))]
	}
	return string(p)
}

// Swapper provides an interface for shuffling. This is a strict subset of
// sort.Interface.
type Swapper interface {
	Swap(i, j int)
	Len() int
}

// Shuffle permutes the elements of data into a random order.
func Shuffle(data Swapper, rng RNG) {
	n := data.Len()
	for i := 0; i < n; i++ {
		data.Swap(i, int(rng.Uintn(uint(n))))
	}
}

// Yield sends values generated from the given distribution. It stops and
// returns once a value is received over the quit channel, if that channel is
// not nil. Useful for safely accessing variates from multiple goroutines.
func Yield(from Distribution, over chan<- float64, quit <-chan struct{}) {
	go func() {
		for {
			select {
			case over <- from.Next():
			case <-quit:
				return
			}
		}
	}()
}

// YieldUint64 sends uint64s from the given RNG. It stops and returns once a
// value is received over the quit channel, if that channel is not nil. Useful
// for safely accessing values from multiple goroutines.
func YieldUint64(from RNG, over chan<- uint64, quit <-chan struct{}) {
	go func() {
		for {
			select {
			case over <- from.Uint64():
			case <-quit:
				return
			}
		}
	}()
}

// StopYielding is a value to tell a yielder to stop. You can also close the
// quit channel if you won't need it again.
var StopYielding struct{}

// SeedInt64 seeds a Seeder using an int64. Specifically, it uses the
// little-endian representation of the int64 as the initialization vector.
func SeedInt64(src Seeder, seed int64) {
	b := [8]byte{}
	binary.LittleEndian.PutUint64(b[:], uint64(seed))
	src.SeedIV(b[:])
}

// CryptoSeeded seeds a Seeder with n random bytes from crypto/rand.Reader and
// returns src.
func CryptoSeeded(src Seeder, n int) Seeder {
	iv := make([]byte, n)
	rand.Read(iv)
	src.SeedIV(iv)
	return src
}
