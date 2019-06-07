package crazy

import "testing"

func BenchmarkYield(b *testing.B) {
	d := Uniform0_1{CryptoSeeded(NewMT64(), mt64N)}
	ch := make(chan float64)
	Yield(d, ch, nil)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = <-ch
		}
	})
}

func BenchmarkYieldUint64(b *testing.B) {
	d := RNG{CryptoSeeded(NewMT64(), mt64N)}
	ch := make(chan uint64)
	YieldUint64(d, ch, nil)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = <-ch
		}
	})
}
