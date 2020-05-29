package auth

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func issueJwt(id uint64, name, email, picture string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = id
	claims["name"] = name
	claims["email"] = email
	claims["picture"] = picture
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}
