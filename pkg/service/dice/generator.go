package dice

import (
	"crypto/rand"
	"math/big"
)

func GenerateDiceValue(sides int) (int, error) {
	max := big.NewInt(int64(sides))

	val, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}

	return int(val.Int64()) + 1, nil
}
