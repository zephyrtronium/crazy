package crazy

import "testing"

func BenchmarkNormal(b *testing.B) {
	d := NewNormal(CryptoSeeded(NewMT64(), mt64N), 0, 1)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = d.Next()
	}
}
