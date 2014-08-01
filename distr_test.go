package crazy

import "testing"

func TestUniform1_2(t *testing.T) {
	d := Uniform1_2{CryptoSeeded(NewMT64(), mt64N)}
	for i := 0; i < 1000000; i++ {
		if x := d.Next(); x < 1 || x >= 2 {
			t.Fail()
		}
	}
}
