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
	p := make([]byte, 1<<30)
	b.SetBytes(int64(len(p)))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		xoro.Read(p)
	}
}

func BenchmarkRexoroshiro(b *testing.B) {
	rexo := CryptoSeeded(NewRexoroshiro(), 16).(*Rexoroshiro)
	p := make([]byte, 1<<30)
	b.SetBytes(int64(len(p)))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		rexo.Read(p)
	}
}
