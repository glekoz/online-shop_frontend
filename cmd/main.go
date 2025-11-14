package main

import (
	"log"
	"net/http"

	"github.com/glekoz/online-shop_frontend/internal/web"
)

func main() {
	web := web.New()

	srv := &http.Server{
		Handler: web.Routes(),
		Addr:    ":8088",
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
