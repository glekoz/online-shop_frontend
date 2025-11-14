package web

import (
	"context"
	"crypto/rsa"
	"errors"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func WriteStatusCode(w http.ResponseWriter, st *status.Status) {
	switch st.Code() {
	case codes.NotFound:
		w.WriteHeader(http.StatusNotFound)
	case codes.AlreadyExists:
		w.WriteHeader(http.StatusConflict)
	case codes.InvalidArgument:
		w.WriteHeader(http.StatusUnprocessableEntity)
	case codes.Internal:
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusTeapot)
	}
}

func ParseIntFromQuery(q string, defaultValue int) (int, error) {
	if q == "" {
		return defaultValue, nil
	}
	num, err := strconv.Atoi(q)
	if err != nil {
		return 0, err
	}
	return num, nil
}
