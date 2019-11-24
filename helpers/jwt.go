package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AccessTokenClaim struct {
	UserID         string `json:"user_id"`
	CredentialUUID string `json:"credential_uuid"`
	jwt.StandardClaims
}

func GenerateAccessToken(userID, credentialUuid string) (string, error) {
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &AccessTokenClaim{
		UserID:         userID,
		CredentialUUID: credentialUuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey())
	if err != nil {
		return "", fmt.Errorf("error exchanging jwt token")
	}

	return tokenString, nil
}

func VerifyAccessToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&AccessTokenClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return secretKey(), nil
		})
	if err != nil {
		return "", err
	}

	return token.Claims.(*AccessTokenClaim).UserID, nil
}

func secretKey() []byte {
	return []byte(os.Getenv("SECRET_KEY"))
}
