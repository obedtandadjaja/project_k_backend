package helpers

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/envy"
)

type SessionTokenClaim struct {
	CredentialID string `json:"credential_id"`
	SessionID    string `json:"session_id"`
	jwt.StandardClaims
}

type AccessTokenClaim struct {
	UserID       string `json:"user_id"`
	CredentialID string `json:"credential_id"`
	jwt.StandardClaims
}

func GenerateSessionToken(credentialID, sessionID string) (string, error) {
	expirationTime := time.Now().Add(10 * 24 * time.Hour)
	claims := &SessionTokenClaim{
		CredentialID: credentialID,
		SessionID:    sessionID,
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

func GenerateAccessToken(userID, credentialID string) (string, error) {
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &AccessTokenClaim{
		UserID:       userID,
		CredentialID: credentialID,
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

func VerifySessionToken(tokenString string) (string, string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&SessionTokenClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return secretKey(), nil
		})
	if err != nil {
		return "", "", err
	}

	return token.Claims.(*SessionTokenClaim).CredentialID,
		token.Claims.(*SessionTokenClaim).SessionID,
		nil
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
	return []byte(envy.Get("SECRET_KEY", "wow-very-discreet-much-secret"))
}
