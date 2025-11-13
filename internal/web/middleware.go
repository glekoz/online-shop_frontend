package web

import (
	"context"
	"net/http"
)

type ctxNameKey struct{}

func (web *Weber) ReadNameFromCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		access, err := r.Cookie("access_token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		name, err := ParseJWTToken(access.Value, web.PublicKeyUser)
		if err != nil {
			w.Write([]byte("Internal"))
		}
		ctx := context.WithValue(r.Context(), ctxNameKey{}, name)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
