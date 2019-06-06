# crazy

Package crazy provides interfaces and implementations for sources of randomness
and pseudo-randomness.

Crazy rejects certain basic assumptions taken by math/rand, including:

- Randomness does not necessarily produce numbers. Rather than returning
  (u)int64s, crazy's Source is identical to io.Reader. (Crazy's sources
  also implement rand.Source64, though.)
- Not all randomness sources can be seeded. In particular, crypto/rand is a
  Source but not a Seeder.
- Sometimes people want to save and restore exact PRNG states. A Saver has
  this capability.

Currently implemented PRNGs are LFG(273, 607), MT64-19937, xoroshiro128+
a modification of xoroshiro128+ that rearranges the output bytes, and
xoshiro256*​*. crypto/rand.Reader naturally implements Source.

The only currently implemented distributions are normal and exponential, but
the ziggurat directory contains a Python script to calculate the necessary
parameters for any monotonically decreasing distribution.

## Which PRNG?

As mentioned above, crazy includes a variety of different generators. For most
applications, xoshiro256** is the best generator because it is very fast, very
small, and has excellent output stream properties. However, there are some
situations where the others may be better:

- MT64-19937 has the property of 623-equidistribution, meaning that every tuple
	of 623 values appears in the output sequence, except the all-zero tuple.
	This property makes it well-suited to applications requiring uniformity in
	many dimensions, like random walks over a highly connected graph. It is,
	however, the slowest generator implemented in crazy.
- LFG(273, 607) is fast with reasonable quality. It is the same generator as
	the one used in the standard library, but it travels in the opposite
	direction and uses a different seeding algorithm. LFG may be suitable for
	moderately dimensional MCMC, where MT64's slow speed doesn't outweigh its
	superior distribution and huge period, but xoshiro's 4-dimensional
	equidistribution is still insufficient.
- xoroshiro128+ has the lowest dimension of PRNGs in crazy, but it is also the
	fastest. For tasks where speed is the only significant factor, and the low
	linear complexity in the low bits of its output stream is acceptable,
	xoroshiro is a good fit.

## Benchmarks

Crazy includes benchmarks for each generator to fill blocks of various sizes.
These benchmarks are named following the convention of BenchmarkGenerator/S,
where Generator is LFG, MT64, Xoroshiro, or Rexoroshiro; and S is 8, K, M, or G
to benchmark filling blocks of size 8 B, 1 kB, 32 MB, or 1 GB, respectively.
Generally, the G tests give the best indication of average performance, M tests
are for consideration of those who don't want to lose a gigabyte of memory, and
K tests give an indication of performance when paging is mitigated. (The state
size of LFG(273, 607) is itself larger than the typical page size for x86/64,
so it may appear slower than usual relative to the others in K tests.) 8 tests
indicate the speed of generating individual values at a time, but this is not
recommended if avoidable.

An example of running benchmarks might look like:

```
> go test -bench /[KG] -benchtime 10s -timeout 1h
goos: windows
goarch: amd64
pkg: github.com/zephyrtronium/crazy
BenchmarkLFG/K-8                50000000               242 ns/op        4228.04 MB/s
BenchmarkLFG/G-8                      50         247917030 ns/op        4331.05 MB/s
BenchmarkMT64/K-8               30000000               586 ns/op        1745.66 MB/s
BenchmarkMT64/G-8                     20         603785375 ns/op        1778.35 MB/s
BenchmarkXoroshiro/K-8          100000000              199 ns/op        5128.92 MB/s
BenchmarkXoroshiro/G-8               100         203835315 ns/op        5267.69 MB/s
BenchmarkRexoroshiro/K-8        100000000              216 ns/op        4736.30 MB/s
BenchmarkRexoroshiro/G-8             100         219902115 ns/op        4882.82 MB/s
BenchmarkXoshiro/K-8            50000000               256 ns/op        3997.27 MB/s
BenchmarkXoshiro/G-8                  50         282325022 ns/op        3803.21 MB/s
PASS
ok      github.com/zephyrtronium/crazy  176.454s
```

In these results, the ns/op measures the time spent to fill an entire 1K or 1G
block, not just to generate a single value). The MB/s throughput is generally a
better indicator of performance. When benchmarking for yourself, use the
`-benchtime` argument to `go test` in order to measure generation of more values.
