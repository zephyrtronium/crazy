package crazy

import "testing"

func TestUniform1_2(t *testing.T) {
	d := Uniform1_2{CryptoSeeded(NewMT64(), mt64N)}
	n := 1 << 28
	if testing.Short() {
		n = 1 << 20
	}
	for i := 0; i < n; i++ {
		if x := d.Next(); x < 1 || x >= 2 {
			t.Fail()
		}
	}
}

func TestUniform0_1(t *testing.T) {
	d := Uniform0_1{CryptoSeeded(NewMT64(), mt64N)}
	n := 1 << 28
	if testing.Short() {
		n = 1 << 20
	}
	for i := 0; i < n; i++ {
		if x := d.Next(); x < 0 || x >= 1 {
			t.Fail()
		}
	}
}

func BenchmarkUniform1_2(b *testing.B) {
	d := Uniform1_2{CryptoSeeded(NewMT64(), mt64N)}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = d.Next()
	}
}

func BenchmarkUniform0_1(b *testing.B) {
	d := Uniform0_1{CryptoSeeded(NewMT64(), mt64N)}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = d.Next()
	}
}

func BenchmarkUniform(b *testing.B) {
	d := Uniform{CryptoSeeded(NewMT64(), mt64N), -1, 1}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = d.Next()
	}
}
