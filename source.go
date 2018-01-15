package crazy

import "io"

// A Source is a source of (pseudo) randomness.
type Source interface {
	// Fill p with random bytes. The actual number of bytes written and any
	// error are returned.
	Read(p []byte) (n int, err error)
}

// A Seeder is a PRNG that can be seeded. After seeding with a particular
// value, all generators of the same type must always produce the same values.
type Seeder interface {
	Source
	// Seed using the given initialization vector. If possible, this should be
	// prepared to handle iv being longer than the state, the same size,
	// shorter, or nil.
	SeedIV(iv []byte)
}

// A Saver is a PRNG that can save and restore its state. A type implementing
// this interface must guarantee that the values produced after saving and
// restoring state are identical to those that would have been produced had
// there been no such actions.
type Saver interface {
	// It makes no sense to be able to save and restore state without being
	// able to seed. Right?
	Seeder
	// Write a representation that can be converted into an equivalent PRNG
	// later. Returns the number of bytes written and any error that occurred.
	Save(into io.Writer) (n int, err error)
	// Load a serialized state of this type of PRNG. Returns the number of
	// bytes read and any error that occurred.
	Restore(from io.Reader) (n int, err error)
}
