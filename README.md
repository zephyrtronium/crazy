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

Currently implemented PRNGs are LFG(273, 607) and MT64-19937.
crypto/rand.Reader naturally implements Source.