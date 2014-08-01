package crazy

import (
	"math/big"
	"testing"
)

func TestUintn(t *testing.T) {
	r := RNG{CryptoSeeded(NewMT64(), mt64N)}
	for i := uint(1); i <= 1000000; i++ {
		if r.Uintn(i) >= i {
			t.Fail()
		}
	}
}

func TestBign(t *testing.T) {
	r := RNG{CryptoSeeded(NewMT64(), mt64N)}
	for i := 1; i <= 30000; i++ {
		v := new(big.Int)
		v.Sub(v.SetBit(v, i, 1), big.NewInt(int64(i)))
		if r.Bign(v).Cmp(v) >= 0 {
			t.Fail()
		}
	}
}
