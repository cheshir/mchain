// Very naive implementation of markov chains.
// You should setup random generator by yourself.
package markovchain

import "math/rand"

// End of chain but EOF is more familiar abbreviation for such purposes.
const EOF Token = ""

// Obviously is not thread safe.
type Chain map[Token]*bucket

type Token string

func New() Chain {
	return make(Chain)
}

func (c Chain) Add(from, to Token) {
	if to == EOF {
		return
	}

	if _, ok := c[from]; !ok {
		c[from] = newBucket()
	}

	c[from].Add(to)
}

// Returns token related to the given one.
// Tokens are chosen randomly with respect to their weights.
// Imagine that we place all available tokens in sequence on x axis started from zero.
// Then we generate random number in range from the first token's start point to the last token's end point.
// Returns token located in generated point.
func (c Chain) Next(from Token) Token {
	bucket := c[from]
	if bucket == nil || bucket.weight == 0 {
		return EOF
	}

	point := rand.Intn(bucket.weight + 1)
	axisXCursor := 0

	for _, item := range bucket.items {
		axisXCursor += item.weight
		if point <= axisXCursor {
			return item.token
		}
	}

	return EOF
}
