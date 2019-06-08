package crazy

// Ziggurat implements a generalized ziggurat algorithm for producing random
// values according to any monotone decreasing density, optionally symmetric
// about zero. See http://www.jstatsoft.org/v05/i08/paper.
type Ziggurat struct {
	// PDF is the probability density function for the desired distribution.
	PDF func(x float64) float64
	// Used to generate a value distributed along the tail of the desired
	// distribution. The value returned should be positive; if it is used and
	// the ziggurat is mirrored, there is a 50% chance of the value being
	// negated before being returned.
	Tail func(src Source) float64
	// If true, one random bit is used to determine the sign of the values
	// produced by the ziggurat.
	Mirrored bool
	// K[i] = floor(2**53 * (x[i-1]/x[i]))
	// K[0] = floor(2**53 * r * PDF(r) / v)
	K [1024]uint64
	// W[i] = x[i] / 2**53
	// W[0] = v*PDF(r) / 2**53
	W [1024]float64
	// F[i] = PDF(x[i])
	F [1024]float64
}

// GenNext generates a value distributed according to the ziggurat.
func (z *Ziggurat) GenNext(src Source) float64 {
	for {
		j := int64(RNG{src}.Uint64())
		i := uint16(j >> 54 & 0x3ff)
		if z.Mirrored {
			j = j << 10 >> 10
		} else {
			j = j & 0x001fffffffffffff
		}
		if uint64(j+j>>63^j>>63) < z.K[i] {
			return float64(j) * z.W[i]
		}
		if i != 0 {
			x := float64(j) * z.W[i]
			if z.F[i]+(Uniform0_1{src}.Next())*(z.F[i-1]-z.F[i]) < z.PDF(x) {
				return x
			}
		} else {
			x := z.Tail(src)
			if j < 0 {
				return -x
			}
			return x
		}
	}
}
