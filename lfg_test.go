package crazy

import (
	"bytes"
	"crypto/rand"
	"testing"
)

// LFG.Seed() is trivial so I'm not going to bother.

func TestLFGSeed(t *testing.T) {
	lfg := NewLFG()
	lfg.SeedIV(nil)
	lfg.SeedIV([]byte{lfgK: 0})
	lfg.SeedIV([]byte{8 * lfgK: 0})
	lfg.SeedIV([]byte{9 * lfgK: 0})
}

func TestLFGSeedConsistency(t *testing.T) {
	iv := make([]byte, 128)
	lfg := NewLFG()
	x, y := make([]byte, 8000), make([]byte, 8000)
	for i := 0; i < 2*lfgK; i++ {
		rand.Read(iv)
		lfg.SeedIV(iv)
		lfg.Read(x)
		lfg.SeedIV(iv)
		lfg.Read(y)
		if !bytes.Equal(x, y) {
			t.Fail()
		}
	}
}

func TestLFGSave(t *testing.T) {
	b := bytes.Buffer{}
	lfg := CryptoSeeded(NewLFG(), lfgK).(*LFG)
	x, y := make([]byte, 8000), make([]byte, 8000)
	for i := 0; i < 2*lfgK; i++ {
		lfg.Save(&b)
		lfg.Read(x)
		lfg.Restore(&b)
		lfg.Read(y)
		if !bytes.Equal(x, y) {
			t.Fail()
		}
		b.Reset()
	}
}
