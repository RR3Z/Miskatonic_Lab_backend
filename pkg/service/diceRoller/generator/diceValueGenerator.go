package generator

import (
	"crypto/rand"
	"math/big"
)

func GenerateDiceValue(sides int) (int, error) {
	max := big.NewInt(int64(sides))

	// Range = [0,sides)
	val, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}

	// Result is [1, sides] (because +1)
	return int(val.Int64()) + 1, nil
}
