package web

import "net/http"

func (web *Weber) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /", web.ReadNameFromCookie(http.HandlerFunc(web.HomeHandler)))

	mux.HandleFunc("GET /register", web.GetRegisterHandler)
	mux.HandleFunc("POST /register", web.PostRegisterHandler)

	mux.HandleFunc("GET /login", web.GetLoginHandler)
	mux.HandleFunc("POST /login", web.PostLoginHandler)

	mux.HandleFunc("GET /logout", web.LogoutHandler)

	return mux
}
