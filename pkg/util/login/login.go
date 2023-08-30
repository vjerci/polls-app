package login

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
)

var ErrSigningFailed = errors.New("failed to sign claims")
var ErrDecodeUnexpectedSignMethod = errors.New("unexpected signing method")
var ErrDecodeClaimsMissing = errors.New("missing claims from token")
var ErrDecodeInvalidToken = errors.New("got invalid token")

var SigningMethod = jwt.SigningMethodHS256

type Client struct {
	JWTSigningKey []byte
}

func New(JWTSigningKey string) *Client {
	return &Client{
		JWTSigningKey: []byte(JWTSigningKey),
	}
}

type AccessToken string

func (client *Client) CreateToken(userID string) (AccessToken, error) {
	t := jwt.NewWithClaims(SigningMethod,
		jwt.MapClaims{
			"user_id": userID,
		})
	s, err := t.SignedString(client.JWTSigningKey)
	if err != nil {
		return "", errors.Join(ErrSigningFailed, err)
	}

	return AccessToken(s), nil
}

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
		userID = claims["user_id"].(string)

		return userID, nil
	}

	return "", ErrDecodeClaimsMissing
}
