package tokens

import (
	cryptorand "crypto/rand"
	"encoding/hex"
)

func Invite() (string, error) {
	return random(16)
}

func Secret() (string, error) {
	return random(24)
}

func random(n int) (string, error) {
	b := make([]byte, n)
	if _, err := cryptorand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
