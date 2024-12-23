package token

import (
	"be-stepup/config"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type CustomClaims struct {
	UserID string `json:"userid"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
	Iss    string `json:"iss"`
	Sub    string `json:"sub"`
}

func (c CustomClaims) Valid() error {
	if time.Unix(c.Exp, 0).Before(time.Now()) {
		return fmt.Errorf("token is expired")
	}
	return nil
}

const TokenExpiryDuration = time.Hour * 12

func GenerateJWT(userid string) (string, error) {
	claims := CustomClaims{
		UserID: userid,
		Exp:    time.Now().Add(TokenExpiryDuration).Unix(),
		Iat:    time.Now().Unix(),
		Iss:    "stepup",
		Sub:    userid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.AuthSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
