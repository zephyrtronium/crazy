// +build go1.9

package crazy

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestXoshiSeed(t *testing.T) {
	xoshi := NewXoshiro()
	xoshi.SeedIV(nil)
	xoshi.SeedIV([]byte{7: 0})
	xoshi.SeedIV([]byte{15: 0})
	xoshi.SeedIV([]byte{23: 0})
	xoshi.SeedIV([]byte{99: 0})
}

func TestXoshiSeedConsistency(t *testing.T) {
	iv := make([]byte, 32)
	xoshi := NewXoshiro()
	x, y := make([]byte, 8000), make([]byte, 8000)
	for i := 0; i < 1024; i++ {
		rand.Read(iv)
		xoshi.SeedIV(iv)
		xoshi.Read(x)
		xoshi.SeedIV(iv)
		xoshi.Read(y)
		if !bytes.Equal(x, y) {
			t.Fail()
		}
	}
}

func TestXoshiSave(t *testing.T) {
	b := bytes.Buffer{}
	xoshi := CryptoSeeded(NewXoshiro(), 32).(*Xoshiro)
	x, y := make([]byte, 8000), make([]byte, 8000)
	for i := 0; i < 1024; i++ {
		xoshi.Save(&b)
		xoshi.Read(x)
		xoshi.Restore(&b)
		xoshi.Read(y)
		if !bytes.Equal(x, y) {
			t.Fail()
		}
		b.Reset()
	}
}

func BenchmarkXoshiro(b *testing.B) {
	xoshi := CryptoSeeded(NewXoshiro(), 32).(*Xoshiro)
	f := func(p []byte) func(b *testing.B) {
		return func(b *testing.B) {
			b.SetBytes(int64(len(p)))
			for n := 0; n < b.N; n++ {
				xoshi.Read(p)
			}
		}
	}
	b.Run("8", f(make([]byte, 8)))
	b.Run("K", f(make([]byte, 1<<10)))
	b.Run("M", f(make([]byte, 1<<25)))
	b.Run("G", f(make([]byte, 1<<30)))
}
