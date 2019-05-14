package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	infoserverctrls "github.com/sebaxtian/truora-rest-api-go/controllers"
)

func main() {
	port := ":3333"
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome"))
	})

	r.Get("/infoserver", infoserverctrls.GetInfoServer())

	fmt.Println("API REST Listen on ", port)
	http.ListenAndServe(port, r)
}
