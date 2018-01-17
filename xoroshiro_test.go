// +build go1.9

package crazy

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestXoroSeed(t *testing.T) {
	lfg := NewLFG()
	lfg.SeedIV(nil)
	lfg.SeedIV([]byte{7: 0})
	lfg.SeedIV([]byte{15: 0})
	lfg.SeedIV([]byte{23: 0})
}

func TestXoroSeedConsistency(t *testing.T) {
	iv := make([]byte, 16)
	xoro := NewXoroshiro()
	x, y := make([]byte, 8000), make([]byte, 8000)
	for i := 0; i < 1024; i++ {
		rand.Read(iv)
		xoro.SeedIV(iv)
		xoro.Read(x)
		xoro.SeedIV(iv)
		xoro.Read(y)
		if !bytes.Equal(x, y) {
			t.Fail()
		}
	}
}

func TestXoroSave(t *testing.T) {
	b := bytes.Buffer{}
	xoro := CryptoSeeded(NewXoroshiro(), 16).(*Xoroshiro)
	x, y := make([]byte, 8000), make([]byte, 8000)
	for i := 0; i < 1024; i++ {
		xoro.Save(&b)
		xoro.Read(x)
		xoro.Restore(&b)
		xoro.Read(y)
		if !bytes.Equal(x, y) {
			t.Fail()
		}
		b.Reset()
	}
}

func BenchmarkXoroshiro(b *testing.B) {
	xoro := CryptoSeeded(NewXoroshiro(), 16).(*Xoroshiro)
	f := func(p []byte) func(b *testing.B) {
		return func(b *testing.B) {
			b.SetBytes(int64(len(p)))
			for n := 0; n < b.N; n++ {
				xoro.Read(p)
			}
		}
	}
	b.Run("8", f(make([]byte, 8)))
	b.Run("K", f(make([]byte, 1<<10)))
	b.Run("M", f(make([]byte, 1<<25)))
	b.Run("G", f(make([]byte, 1<<30)))
}

func BenchmarkRexoroshiro(b *testing.B) {
	rexo := CryptoSeeded(NewRexoroshiro(), 16).(*Rexoroshiro)
	f := func(p []byte) func(b *testing.B) {
		return func(b *testing.B) {
			b.SetBytes(int64(len(p)))
			for n := 0; n < b.N; n++ {
				rexo.Read(p)
			}
		}
	}
	b.Run("8", f(make([]byte, 8)))
	b.Run("K", f(make([]byte, 1<<10)))
	b.Run("M", f(make([]byte, 1<<25)))
	b.Run("G", f(make([]byte, 1<<30)))
}
