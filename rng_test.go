package crazy

import (
	"math/big"
	"testing"
)

func TestUintn(t *testing.T) {
	r := RNG{CryptoSeeded(NewMT64(), mt64N)}
	n := uint(1 << 28)
	if testing.Short() {
		n = 1 << 20
	}
	for i := uint(1); i <= n; i++ {
		if r.Uintn(i) >= i {
			t.Fail()
		}
	}
}

func TestBign(t *testing.T) {
	r := RNG{CryptoSeeded(NewMT64(), mt64N)}
	n := 1 << 18
	if testing.Short() {
		n = 1 << 12
	}
	for i := 1; i <= n; i++ {
		v := new(big.Int)
		v.Sub(v.SetBit(v, i, 1), big.NewInt(int64(i)))
		if r.Bign(v).Cmp(v) >= 0 {
			t.Fail()
		}
	}
}
