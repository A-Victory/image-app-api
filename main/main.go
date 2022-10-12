package main

import (
	"log"
	"net/http"
	"os"

	"app/routes"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading environment variables: ", err)
	}
	r := httprouter.New()

	port := os.Getenv("PORT")

	routes.Routes(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Fatal(srv.ListenAndServe())
}
