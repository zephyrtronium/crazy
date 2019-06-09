package crazy

import "io"

// A Source is a source of (pseudo) randomness.
type Source interface {
	// Read fills p with random bytes. The actual number of bytes written and
	// any error are returned.
	Read(p []byte) (n int, err error)
}

// A Seeder is a PRNG that can be seeded. After seeding with a particular
// value, all generators of the same type must always produce the same values.
type Seeder interface {
	Source
	// SeedIV seeds using the given initialization vector. iv may be longer
	// than the state, the same size, shorter, or nil.
	SeedIV(iv []byte)
}

// A Saver is a PRNG that can save and restore its state. A generator that
// saves its state must produce the same output stream thenceforth as another
// generator of the same type which restores that state.
type Saver interface {
	Seeder
	// Save writes a representation that can be restored into an equivalent
	// PRNG later. Returns the number of bytes written and any error that
	// occurred.
	Save(into io.Writer) (n int, err error)
	// Restore loads a serialized state of this type of PRNG. Returns the
	// number of bytes read and any error that occurred.
	Restore(from io.Reader) (n int, err error)
}

// A Copier is a PRNG that can produce a copy of itself that will generate the
// same sequence of values. This interface is only necessary when using PRNGs
// through interfaces; PRNG values can be copied directly by dereferencing.
type Copier interface {
	Seeder
	// Copy creates a copy of the PRNG.
	Copy() Copier
}

// A Jumper is a PRNG that can efficiently "jump" to a new state, such that the
// output sequences of the pre- and post-jump states will not overlap for a
// long time. This facilitates parallel scaling by ensuring that many separate
// random processes each have unique contributions to a single problem.
type Jumper interface {
	Seeder
	// Jump advances the state of the PRNG by 2**(64+n) iterations.
	Jump(n uint)
}
