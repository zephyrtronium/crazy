package crazy

import (
	"bytes"
	"crypto/rand"
	"testing"
)

// MT64.SeedIV() uses .Seed(), so testing IV tests both.

func TestMT64Seed(t *testing.T) {
	mt := NewMT64()
	mt.SeedIV(nil)
	mt.SeedIV([]byte{mt64N: 0})
	mt.SeedIV([]byte{8 * mt64N: 0})
	mt.SeedIV([]byte{9 * mt64N: 0})
}

func TestMT64SeedConsistency(t *testing.T) {
	iv := make([]byte, 128)
	mt := NewMT64()
	x, y := make([]byte, 8000), make([]byte, 8000)
	for i := 0; i < 2*mt64N; i++ {
		rand.Read(iv)
		mt.SeedIV(iv)
		mt.Read(x)
		mt.SeedIV(iv)
		mt.Read(y)
		if !bytes.Equal(x, y) {
			t.Fail()
		}
	}
}

func TestMT64Save(t *testing.T) {
	b := bytes.Buffer{}
	mt := CryptoSeeded(NewMT64(), mt64N).(*MT64)
	x, y := make([]byte, 8000), make([]byte, 8000)
	for i := 0; i < 2*mt64N; i++ {
		mt.Save(&b)
		mt.Read(x)
		mt.Restore(&b)
		mt.Read(y)
		if !bytes.Equal(x, y) {
			t.Fail()
		}
		b.Reset()
	}
}

func BenchmarkMT64(b *testing.B) {
	mt := CryptoSeeded(NewMT64(), mt64N).(*MT64)
	f := func(p []byte) func(b *testing.B) {
		return func(b *testing.B) {
			b.SetBytes(int64(len(p)))
			for n := 0; n < b.N; n++ {
				mt.Read(p)
			}
		}
	}
	b.Run("8", f(make([]byte, 8)))
	b.Run("K", f(make([]byte, 1<<10)))
	b.Run("M", f(make([]byte, 1<<25)))
	b.Run("G", f(make([]byte, 1<<30)))
}
