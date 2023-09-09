package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
)

var signingMethod = jwt.SigningMethodHS256

type Client struct {
	JWTSigningKey []byte
}

func New(jwtSigningKey string) *Client {
	return &Client{
		JWTSigningKey: []byte(jwtSigningKey),
	}
}

type AccessToken string

var ErrSigningFailed = errors.New("failed to sign claims")

func (client *Client) CreateToken(userID string) (AccessToken, error) {
	t := jwt.NewWithClaims(signingMethod,
		jwt.MapClaims{
			"user_id": userID,
		})

	s, err := t.SignedString(client.JWTSigningKey)
	if err != nil {
		return "", errors.Join(ErrSigningFailed, err)
	}

	return AccessToken(s), nil
}

var ErrDecodeUnexpectedSignMethod = errors.New("unexpected signing method")
var ErrDecodeClaimsMissing = errors.New("missing claims from token")
var ErrDecodeInvalidToken = errors.New("got invalid token")
var ErrDecodeUserIDNotString = errors.New("user id is not of type string")

func (client *Client) Decode(input AccessToken) (userID string, err error) {
	token, err := jwt.Parse(string(input), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w\n signing method is %v", ErrDecodeUnexpectedSignMethod, token.Header["alg"])
		}

		return client.JWTSigningKey, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", ErrDecodeInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", ErrDecodeUserIDNotString
		}

		return userID, nil
	}

	return "", ErrDecodeClaimsMissing
}
