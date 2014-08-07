package crazy

import (
	"crypto/rand"
	"encoding/binary"
	mathrand "math/rand"
)

// Compose two Sources using XOR.
type XorCompose struct {
	A, B Source
	buf  []byte
}

// A ^ B
func (x XorCompose) Read(p []byte) (int, error) {
	if cap(x.buf) < len(p) {
		x.buf = make([]byte, len(p))
	} else {
		x.buf = x.buf[:len(p)]
	}
	x.A.Read(p)
	x.B.Read(x.buf)
	for i, v := range x.buf {
		p[i] ^= v
	}
	return len(p), nil
}

type math2crazy struct {
	mathrand.Source
}

// Adapt a math/rand Source into a crazy Seeder. If the argument already
// implements Seeder, it is returned directly. Otherwise, to preserve
// equidistribution, the Read() method of this Seeder uses two Int63() calls
// per eight bytes of output, using one fewer call when the output is not
// divisible by 8. Additionally, the SeedIV() method uses only up to the first
// eight bytes of the argument.
func AdaptRand(src mathrand.Source) Seeder {
	if s, ok := src.(Seeder); ok {
		return s
	}
	return math2crazy{src}
}

func (m math2crazy) Read(p []byte) (n int, err error) {
	n = len(p)
	for len(p) >= 8 {
		binary.LittleEndian.PutUint64(p, uint64(m.Int63()^m.Int63()<<1))
		p = p[8:]
	}
	if len(p) == 8 {
		binary.LittleEndian.PutUint64(p, uint64(m.Int63()^m.Int63()<<1))
	} else {
		b := [8]byte{}
		binary.LittleEndian.PutUint64(b[:], uint64(m.Int63()))
		copy(p, b[:])
	}
	return n, nil
}

func (m math2crazy) SeedIV(iv []byte) {
	b := [8]byte{}
	copy(b[:], iv)
	m.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

type crazy2math struct {
	Seeder
}

// Adapt a crazy Seeder into a math/rand Source. If the argument already
// implements math/rand.Source, it is returned directly.
func AdaptCrazy(src Seeder) mathrand.Source {
	if m, ok := src.(mathrand.Source); ok {
		return m
	}
	return crazy2math{src}
}

func (c crazy2math) Int63() int64 {
	return int64(RNG{c.Seeder}.Uint64() >> 1)
}

func (c crazy2math) Seed(x int64) {
	SeedInt64(c.Seeder, x)
}

// Create a random string with specified length using only characters from the
// given alphabet. The mechanism for constructing the string is choosing random
// positions from the alphabet, so characters which appear more often have
// higher relative probability.
func RandString(rng RNG, alphabet string, length int) string {
	p := make([]rune, length)
	a := []rune(alphabet)
	for i := range p {
		p[i] = a[int(rng.Uintn(uint(len(a))))]
	}
	return string(p)
}

// Swapping interface for shuffling.
type Swapper interface {
	Swap(i, j int)
	Len() int
}

// Permute the elements of data into a random order.
func Shuffle(data Swapper, rng RNG) {
	n := data.Len()
	for i := 0; i < n; i++ {
		data.Swap(i, int(rng.Uintn(uint(n))))
	}
}

// Yield values generated from the given distribution. Stop once a value is
// received over the quit channel, if that channel is not nil. Useful for
// synchronizing a distribution between multiple goroutines.
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

// Yield uint64s from the given RNG. Stop once a value is received over the
// quit channel, if that channel is not nil. Useful for synchronizing a
// distribution between multiple goroutines.
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

// Tell a yielder to stop.
var StopYielding struct{}

// Seed a Seeder using an int64.
func SeedInt64(src Seeder, seed int64) {
	b := [8]byte{}
	binary.LittleEndian.PutUint64(b[:], uint64(seed))
	src.SeedIV(b[:])
}

// Seed a Seeder with n random bytes from crypto/rand.Reader. Returns src.
func CryptoSeeded(src Seeder, n int) Seeder {
	iv := make([]byte, n)
	rand.Read(iv)
	src.SeedIV(iv)
	return src
}
