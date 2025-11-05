package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Serve() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // logs a panic, returns HTTP 500

	r.Get("/",
		func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("hi mom"))
			if err != nil {
				log.Printf("failed to send response")
			}
		})

	err := http.ListenAndServe(":4200", r)
	if err != nil {
		log.Fatalf("yikes, failed to start...")
	}
}
