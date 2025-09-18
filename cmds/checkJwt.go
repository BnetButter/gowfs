package main 
import (
	"github.com/golang-jwt/jwt/v5"
	"fmt"
)

const JWT_SECRET = "e2aeda7be71ff3a9c24048b5bae69ff60f3eb85d80e56ae8617c2c0bd4f24cde"

// ParseSub extracts the "sub" claim from a JWT string as an int.
func ParseSub(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// Since you always store `sub` as an int, just assert directly
	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, fmt.Errorf("sub claim not an int")
	}

	return int(sub), nil
}

func main() {
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjN9.5z5GA65oR__Viieu_1it2Bjr2Ycj-DUIwDISMGfnXIQ"
	id, err := ParseSub(jwt)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	fmt.Println("OK, err", id, err);


}