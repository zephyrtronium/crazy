# crazy

Package crazy provides interfaces and implementations for sources of randomness
and pseudo-randomness.

Crazy rejects certain basic assumptions taken by math/rand, including:

- Randomness does not necessarily produce numbers. Rather than returning an
  int64, crazy's Source is identical to io.Reader.
- Not all randomness sources can be seeded. In particular, crypto/rand is a
  Source but not a Seeder.
- Sometimes people want to save and restore exact PRNG states. A Saver has
  this capability.

Currently implemented PRNGs are LFG(273, 607), MT64-19937, xoroshiro128+, and
a modification of xoroshiro128+ that increases the entropy of low-order bits.
crypto/rand.Reader naturally implements Source.

The only currently implemented distributions are normal and exponential, but
the ziggurat directory contains a Python script to calculate the necessary
parameters for any monotonically decreasing distribution.

## Which PRNG?

As mentioned above, crazy includes a variety of different generators. For most
applications, xoroshiro128+ is the best generator because it is the fastest.
However, there are some situations where the others may be better:

- rexoroshiro may be better for situations where the PRNG is most frequently
	used to choose between two to four options, especially true/false, because
	plain xoroshiro128+ is known to have relatively low entropy in its lowest
	bits. Most applications won't notice the difference, but it may sometimes
	be important.
- MT64-19937 has the property of 623-dimensional k-distribution, which means in
	practice that it is the best option for things like MCMC where high
	uniformity in many dimensions is required. It is, however, the slowest
	generator implemented in crazy: about a third the throughput of
	xoroshiro128+ and a bit less than half that of LFG. Eventually, I'd like to
	implement one of the WELL generators for even better dimensionality at
	faster speeds.
- LFG(273, 607) is the same PRNG implemented in math/rand. Therefore, it can be
	used where interoperability with the same is required; for example, to
	transition from the standard library to crazy in order to be able to save
	and restore state while not risking changes to test results.
- Cryptographically secure pseudo-random number generators are outside the
	scope of crazy. There will never be a CSPRNG in crazy. Avoid using any
	crazy generator for any cryptographic applications.

## Benchmarks

Crazy includes benchmarks for each generator to fill blocks of various sizes.
These benchmarks are named following the convention of BenchmarkGenerator/S,
where Generator is LFG, MT64, Xoroshiro, or Rexoroshiro; and S is 8, K, M, or G
to benchmark filling blocks of size 8 B, 1 kB, 32 MB, or 1 GB, respectively.
Generally, the G tests give the best indication of average performance, M tests
are for consideration of those who don't want to lose a gigabyte of memory, and
K tests give an indication of performance when paging is mitigated. (The state
size of LFG(273, 607) is itself larger than the typical page size for x86/64,
so it may appear slower than usual relative to the others in K tests.)

An example of running benchmarks might look like:

```
> go test -bench /[KG] -benchtime 1m -timeout 24h
goos: windows
goarch: amd64
pkg: github.com/zephyrtronium/crazy
BenchmarkLFG/K-8                300000000              246 ns/op        4157.82 MB/s
BenchmarkLFG/G-8                     300         249721152 ns/op        4299.76 MB/s
BenchmarkMT64/K-8               100000000              602 ns/op        1700.07 MB/s
BenchmarkMT64/G-8                    100         621177901 ns/op        1728.56 MB/s
BenchmarkXoroshiro/K-8          500000000              200 ns/op        5102.19 MB/s
BenchmarkXoroshiro/G-8               500         207571328 ns/op        5172.88 MB/s
BenchmarkRexoroshiro/K-8        500000000              206 ns/op        4957.86 MB/s
BenchmarkRexoroshiro/G-8             500         213313608 ns/op        5033.63 MB/s
PASS
ok      github.com/zephyrtronium/crazy  821.086s
```

In these results, the ns/op measures the time spent to fill an entire 1K or 1G
block (_not_ just to generate a single value). The MB/s throughput is generally
a better indicator of performance. It is encouraged to use the `-benchtime`
argument to `go test` in order to measure generation of more values.
