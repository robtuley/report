package report

import (
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var entropyPool = sync.Pool{
	New: func() interface{} {
		return rand.New(rand.NewSource(time.Now().UnixNano()))
	},
}

// ULID creates a Universally Unique Lexicographically Sortable Identifier
//
// A ULID is:
// * compatible with UUID/GUID
// * Lexicographically sortable
// * Canonically encoded as a 26 character string, as opposed
//   to the 36 character ID
// * Uses Crockford's base32 for better efficiency and readability
//   (5 bits per character)
// * Case insensitive
// * No special characters (URL safe)
//
func createULID() string {
	entropy := entropyPool.Get().(*rand.Rand)
	defer func() {
		entropyPool.Put(entropy)
	}()

	u, err := ulid.New(ulid.Timestamp(time.Now()), entropy)
	if err != nil {
		return randString(26)
	}
	return u.String()
}

// Crockford's Base32 is used as shown. This alphabet excludes the
// letters I, L, O, and U to avoid confusion and abuse.
const randStringLetters = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = randStringLetters[rand.Intn(len(randStringLetters))]
	}
	return string(b)
}
