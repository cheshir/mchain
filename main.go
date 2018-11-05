package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/cheshir/mchains/markovchain"
)

const (
	startToken markovchain.Token = "^"
	endToken                     = "$"

	maxSentenceLength = 20
)

var sourcePath *string

func init() {
	rand.Seed(time.Now().UnixNano())

	sourcePath = flag.String("source", "", "Path to source file with learning data.")
	flag.Parse()
}

func main() {
	source, err := ioutil.ReadFile(*sourcePath)
	if err != nil {
		log.Fatalf("Failed to read source file: %v", err)
	}

	log.Println("[INFO] Parse source")
	tokens, err := parseSource(string(source))
	if err != nil {
		log.Fatalf("Failed to parse source to tokens: %v", err)
	}

	log.Println("[INFO] Learn model")
	chain := markovchain.New()
	learn(chain, tokens)

	log.Println("[INFO] Generate sentences")
	for i := 0; i < 10; i++ {
		fmt.Println(generateSentence(chain))
	}
}

// Not optimal, but quite fast for prototype.
func parseSource(source string) ([][]markovchain.Token, error) {
	sentences := strings.Split(source, ".")
	result := make([][]markovchain.Token, 0, len(sentences))

	for _, sentence := range sentences {
		tokens := parseSentence(sentence)

		if len(tokens) != 0 {
			result = append(result, tokens)
		}
	}

	return result, nil
}

// Data normalization should be more complex.
// Right now source should be normalized before usage.
func parseSentence(sentence string) []markovchain.Token {
	// Straightforward solution.
	parsed := strings.Fields(sentence)
	if len(parsed) == 0 {
		return nil
	}

	tokens := make([]markovchain.Token, 0, len(parsed)+2)
	tokens = append(tokens, startToken)
	for _, chunk := range parsed {
		tokens = append(tokens, markovchain.Token(chunk))
	}
	tokens = append(tokens, endToken)

	return tokens
}

func learn(chain markovchain.Chain, tokens [][]markovchain.Token) {
	for _, sentence := range tokens {
		for from, to := 0, 1; to < len(sentence); from, to = from+1, to+1 {
			chain.Add(sentence[from], sentence[to])
		}
	}
}

func generateSentence(chain markovchain.Chain) string {
	var tokens []string
	for token, i := chain.Next(startToken), 0; token != endToken && token != markovchain.EOF && i < maxSentenceLength; token, i = chain.Next(token), i+1 {
		tokens = append(tokens, string(token))
	}

	return strings.Join(tokens, " ")
}
