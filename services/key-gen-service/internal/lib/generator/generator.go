package generator

import (
	"crypto/rand"
	"math/big"
	mrand "math/rand"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type Generator struct {
	length  int
	charset []byte
}

func NewGenerator(length int) *Generator {
	shuffled := []byte(charset)
	mrand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return &Generator{
		length:  length,
		charset: shuffled,
	}
}

func (g *Generator) Generate() string {
	result := make([]byte, g.length)
	max := big.NewInt(int64(len(g.charset)))
	for i := range result {
		randInt, err := rand.Int(rand.Reader, max)
		if err != nil {
			continue
		}
		result[i] = g.charset[randInt.Int64()]
	}

	return string(result)
}
