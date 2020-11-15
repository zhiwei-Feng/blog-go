package config

import (
	"github.com/kataras/iris/v12/middleware/jwt"
	"time"
)

type UserClaims struct {
	jwt.Claims
	Username string
	UserId   uint
}

var J *jwt.JWT

func init() {
	J = jwt.HMAC(1*time.Hour, "secret")
}

func Authenticate(username string, userId uint) (string, error) {
	standardClaims := jwt.Claims{Issuer: "blog-go"}
	// NOTE: if custom claims then the `j.Expiry(claims)` (or jwt.Expiry(duration, claims))
	// MUST be called in order to set the expiration time.
	customClaims := UserClaims{
		Claims:   J.Expiry(standardClaims),
		Username: username,
		UserId:   userId,
	}
	return J.Token(customClaims)
}
