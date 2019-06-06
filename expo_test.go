package crazy

import "testing"

func TestExponential(t *testing.T) {
	d := NewExponential(CryptoSeeded(NewMT64(), mt64N), 1)
	n := 1 << 28
	if testing.Short() {
		n = 1 << 20
	}
	for i := 0; i < n; i++ {
		if x := d.Next(); x < 0 {
			t.Fail()
		}
	}
}

func BenchmarkExponential(b *testing.B) {
	d := NewExponential(CryptoSeeded(NewMT64(), mt64N), 1)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = d.Next()
	}
}
