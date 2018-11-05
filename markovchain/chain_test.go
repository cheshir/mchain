package markovchain

import (
	"reflect"
	"testing"
)

func TestChain_Add(t *testing.T) {
	tt := []struct {
		Description   string
		From          Token
		To            Token
		InitialChain  Chain
		ExpectedChain Chain
	}{
		{
			Description:  "Adding new token",
			From:         "hello",
			To:           "world",
			InitialChain: New(),
			ExpectedChain: Chain(map[Token]*bucket{
				"hello": {
					items: []bucketItem{
						{
							token:  "world",
							weight: 1,
						},
					},
					weight: 1,
				},
			}),
		},
		{
			Description: "Extending token",
			From:        "hello",
			To:          "puppy",
			InitialChain: Chain(map[Token]*bucket{
				"hello": {
					items: []bucketItem{
						{
							token:  "world",
							weight: 1,
						},
					},
					weight: 1,
				},
			}),
			ExpectedChain: Chain(map[Token]*bucket{
				"hello": {
					items: []bucketItem{
						{
							token:  "world",
							weight: 1,
						},
						{
							token:  "puppy",
							weight: 1,
						},
					},
					weight: 2,
				},
			}),
		},
		{
			Description: "Adding the same token",
			From:        "hello",
			To:          "world",
			InitialChain: Chain(map[Token]*bucket{
				"hello": {
					items: []bucketItem{
						{
							token:  "world",
							weight: 1,
						},
					},
					weight: 1,
				},
			}),
			ExpectedChain: Chain(map[Token]*bucket{
				"hello": {
					items: []bucketItem{
						{
							token:  "world",
							weight: 2,
						},
					},
					weight: 2,
				},
			}),
		},
		{
			Description: "Adding EOF token",
			From:        "hello",
			To:          EOF,
			InitialChain: Chain(map[Token]*bucket{
				"hello": {
					items: []bucketItem{
						{
							token:  "world",
							weight: 1,
						},
					},
					weight: 1,
				},
			}),
			ExpectedChain: Chain(map[Token]*bucket{
				"hello": {
					items: []bucketItem{
						{
							token:  "world",
							weight: 1,
						},
					},
					weight: 1,
				},
			}),
		},
	}

	for _, tc := range tt {
		t.Run(tc.Description, func(t *testing.T) {
			actualChain := tc.InitialChain
			actualChain.Add(tc.From, tc.To)
			if !chainsAreEqual(tc.ExpectedChain, actualChain) {
				t.Logf("actual: %v", actualChain["hello"])
				t.Errorf("chains are not equal\nExpected:\t%#v\nActual:\t\t%#v", tc.ExpectedChain, actualChain)
			}
		})
	}
}

func chainsAreEqual(expected, actual Chain) bool {
	if len(expected) != len(actual) {
		return false
	}

	for tokenFrom, expectedBucket := range expected {
		actualBucket, ok := actual[tokenFrom]
		if !ok {
			return false
		}

		if len(expectedBucket.items) != len(actualBucket.items) || expectedBucket.weight != actualBucket.weight {
			return false
		}

		if !reflect.DeepEqual(expectedBucket.items, actualBucket.items) {
			return false
		}
	}

	return true
}

func TestChain_Next(t *testing.T) {
	tt := []struct {
		Description   string
		FromToken     Token
		InitialChain  Chain
		ExpectedToken Token
	}{
		{
			Description: "Getting token",
			FromToken:   "hello",
			InitialChain: Chain(map[Token]*bucket{
				"hello": {
					items: []bucketItem{
						{
							token:  "world",
							weight: 1,
						},
					},
					weight: 1,
				},
			}),
			ExpectedToken: "world",
		},
		{
			Description: "Getting non existent token",
			FromToken:   "hey",
			InitialChain: Chain(map[Token]*bucket{
				"hello": {
					items: []bucketItem{
						{
							token:  "world",
							weight: 1,
						},
					},
					weight: 1,
				},
			}),
			ExpectedToken: EOF,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Description, func(t *testing.T) {
			actualToken := tc.InitialChain.Next(tc.FromToken)
			if tc.ExpectedToken != actualToken {
				t.Errorf("tokens are not equal\nExpected:\t%#v\nActual:\t\t%#v", tc.ExpectedToken, actualToken)
			}
		})
	}
}
