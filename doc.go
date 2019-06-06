/*
Package crazy provides interfaces and implementations for sources of randomness
and pseudo-randomness.

Crazy does not assume that randomness is in the form of numbers. crazy.Source
is a synonym for io.Reader; random bytes are produced. This also means that any
Reader can be used to produce values, including crypto/rand.Reader. Turning
a Source into numbers is the job of RNGs (for integers) and Distributions (for
floats).

Not all randomness sources, especially including crypto/rand, can be seeded.
Crazy supports this concept by splitting seeding into a separate interface.
Since (most) PRNGs have much more entropy than can be encoded in 64 bits,
seeding is done with []byte instead of int64.

Saving and restoring exact PRNG states can be useful. math/rand seems to assume
that seeding is enough, but a crazy Saver allows a PRNG to be saved into any
io.Writer and restored from any io.Reader.

Currently implemented PRNGs are LFG(273, 607), MT64-19937, xoroshiro128+,
and xoshiro256**. io.Reader and, in particular, crypto/rand.Reader naturally
implement Source.

The only currently implemented distributions are normal and exponential, but
the ziggurat directory contains a Python script to calculate the necessary
parameters for any monotonically decreasing distribution.
*/
package crazy
