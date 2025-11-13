package web

import (
	"context"
	"crypto/rsa"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func ParseJWTToken(tokenString string, publicKey *rsa.PublicKey) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
		return publicKey, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("ну а что поделать")
	}
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return "", errors.New("ну а что поделать")
	}
	data, ok := (*claims)["data"].(map[string]any)
	if !ok {
		return "", errors.New("ну а что поделать")
	}
	return data["Name"].(string), nil
}

func readNameFromCtx(ctx context.Context) (string, bool) {
	name, ok := ctx.Value(ctxNameKey{}).(string)
	if !ok {
		return "", false
	}
	return name, true
}
