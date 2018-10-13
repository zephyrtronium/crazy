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

Currently implemented PRNGs are LFG(273, 607), MT64-19937, xoroshiro128+
a modification of xoroshiro128+ that rearranges the output bytes, and
xoshiro256*​*. crypto/rand.Reader naturally implements Source.

The only currently implemented distributions are normal and exponential, but
the ziggurat directory contains a Python script to calculate the necessary
parameters for any monotonically decreasing distribution.

## Which PRNG?

As mentioned above, crazy includes a variety of different generators. For most
applications, xoshiro256** is the best generator because it is almost the
fastest and has excellent output stream properties. However, there are some
situations where the others may be better:

- MT64-19937 has the property of 623-equidistribution, meaning that every tuple
	of 623 values appears in the output sequence, except the all-zero tuple.
	This property makes it well-suited to applications requiring uniformity in
	many dimensions, like MCMC over a large graph. It is, however, the slowest
	generator implemented in crazy: about a third the throughput of xoshiro and
	a bit less than half that of LFG.
- LFG(273, 607) is fast with reasonable quality. It is the same generator as
	the one used in the standard library, but it travels in the opposite
	direction and uses a different seeding algorithm. If the standard library
	works for you but xoshiro doesn't, and MT64 is too slow, LFG should work.

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

The package currently contains some testing artifacts related to unusual
slowdown observed in Xoroshiro between Go 1.9.7 and 1.10. In particular, the
package currently may only build on amd64, and the Asmxoro function will go
away once the issues are resolved.
