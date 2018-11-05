package markovchain

const (
	defaultWeight         = 1
	defaultBucketCapacity = 1 // Bucket will store at least one element.
)

// Chain links.
// Is not efficient for insert operations, optimized for reading.
// FIXME Make sense to split learning and usage steps to use optimized data structures.
type bucket struct {
	items  []bucketItem
	weight int // Cum weight.
}

type bucketItem struct {
	token  Token
	weight int
}

func newBucket() *bucket {
	return &bucket{
		items: make([]bucketItem, 0, defaultBucketCapacity),
	}
}

// Upsert items collection with given token.
func (b *bucket) Add(token Token) {
	upserted := false

	for i, item := range b.items {
		if item.token == token {
			item.weight++
			b.items[i] = item
			upserted = true
		}
	}

	if !upserted {
		b.items = append(b.items, bucketItem{token, defaultWeight})
	}

	b.weight++
}
