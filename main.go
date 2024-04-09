package main

import (
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	router := chi.NewRouter()
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Println("Server is running on port " + portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
