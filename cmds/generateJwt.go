package main;

import (
	"github.com/golang-jwt/jwt/v5"
	"fmt"
)

func main() {
		// Replace with your secret (e.g. from `openssl rand -base64 32`)
	secret := []byte("e2aeda7be71ff3a9c24048b5bae69ff60f3eb85d80e56ae8617c2c0bd4f24cde")

	for i := 0; i < 4; i++ {
		// Create claims with only "sub"
		claims := jwt.MapClaims{
			"sub": i, // subject claim
			// no "exp" → no expiration
		}

		// Create the token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Sign it
		signed, err := token.SignedString(secret)
		if err != nil {
			panic("failed to sign")
		}

		fmt.Printf("JWT %d: %s\n", i, signed)
	}
	
}