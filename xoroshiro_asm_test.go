package crazy

import (
	"bytes"
	"testing"
)

func TestAsm(t *testing.T) {
	xoro1 := CryptoSeeded(NewXoroshiro(), 16).(*Xoroshiro)
	xoro2 := NewXoroshiro()
	b := bytes.Buffer{}
	xoro1.Save(&b)
	xoro2.Restore(&b)
	for i := 0; i < 1e6; i++ {
		if xoro1.Uint64() != xoro2.Uint64() {
			t.Fail()
		}
	}
}

func BenchmarkXoroshiroUint64(b *testing.B) {
	xoro := CryptoSeeded(NewXoroshiro(), 16).(*Xoroshiro)
	for n := 0; n < b.N; n++ {
		xoro.Uint64()
	}
}

func BenchmarkAsmxoroUint64(b *testing.B) {
	xoro := CryptoSeeded(NewXoroshiro(), 16).(*Xoroshiro)
	for n := 0; n < b.N; n++ {
		Asmxoro(xoro)
	}
}

func BenchmarkLFGUint64(b *testing.B) {
	lfg := CryptoSeeded(NewLFG(), lfgK).(*LFG)
	for n := 0; n < b.N; n++ {
		lfg.Uint64()
	}
}
