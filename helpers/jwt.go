package helpers

import (
	"errors"
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
	UserType     string `json:"user_type"`
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
		return "", errors.New("error exchanging jwt token")
	}

	return tokenString, nil
}

func GenerateAccessToken(userID, credentialID, userType string) (string, error) {
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &AccessTokenClaim{
		UserID:       userID,
		CredentialID: credentialID,
		UserType:     userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey())
	if err != nil {
		return "", errors.New("error exchanging jwt token")
	}

	return tokenString, nil
}

// returns credentialID, sessionID, error
func VerifySessionToken(tokenString string) (*SessionTokenClaim, error) {
	if len(tokenString) == 0 {
		return nil, errors.New("Invalid session token")
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&SessionTokenClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("There was an error")
			}
			return secretKey(), nil
		})
	if err != nil {
		return nil, err
	}

	return token.Claims.(*SessionTokenClaim), nil
}

// returns userID, error
func VerifyAccessToken(tokenString string) (*AccessTokenClaim, error) {
	if len(tokenString) == 0 {
		return nil, errors.New("Invalid access token")
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&AccessTokenClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("There was an error")
			}
			return secretKey(), nil
		})
	if err != nil {
		return nil, err
	}

	return token.Claims.(*AccessTokenClaim), nil
}

func secretKey() []byte {
	return []byte(envy.Get("SECRET_KEY", "wow-very-discreet-much-secret"))
}
