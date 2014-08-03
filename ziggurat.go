package crazy

// Generalized ziggurat algorithm for producing random values according to any
// decreasing density. See http://www.jstatsoft.org/v05/i08/paper [PDF].
type Ziggurat struct {
	// Probability density function for the desired distribution.
	PDF func(x float64) float64
	// Used to generate a value distributed along the tail of the desired
	// distribution. The value returned should be positive; if it is used and
	// the ziggurat is mirrored, there is a 50% chance of the value being
	// negated before being returned.
	Tail func(src Source) float64
	// If true, one random bit is used to determine the sign of the values
	// produced by the ziggurat.
	Mirrored bool
	// K[i] = floor(2**32 * (x[i-1]/x[i]))
	// K[0] = floor(2**32 * r * PDF(r) / v)
	K [128]uint32
	// W[i] = x[i] / 2**32
	// W[0] = v*PDF(r) / 2**32
	W [128]float32
	// F[i] = PDF(x[i])
	F [128]float32
}

// Generate a value distributed according to the ziggurat.
func (z *Ziggurat) GenNext(src Source) float64 {
	for {
		j := int32(RNG{src}.Uint32())
		if !z.Mirrored {
			j &= 0x7fffffff
		}
		i := j & 127
		if uint32(j+j>>31^j>>31) < z.K[i] {
			return float64(j) * float64(z.W[i])
		}
		if i != 0 {
			x := float64(j) * float64(z.W[i])
			if z.F[i]+float32(Uniform1_2{src}.Next()-1)*(z.F[i-1]-z.F[i]) < float32(z.PDF(x)) {
				return x
			}
		} else {
			x := z.Tail(src)
			if z.Mirrored && j < 0 {
				return -x
			}
			return x
		}
	}
}
