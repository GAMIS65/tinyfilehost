package util

import (
	"crypto/rand"
	"errors"
	"math/big"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateId(length int) (string, error) {
	if length < 0 {
		return "", errors.New("length must be non-negative")
	}

	id := make([]byte, length)

	for i := range id {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		id[i] = alphabet[n.Int64()]
	}
	return string(id), nil
}