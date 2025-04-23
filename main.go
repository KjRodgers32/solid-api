package main

import (
	"fmt"
	"github.com/lpernett/godotenv"
	"log"
	"net/http"
	"os"
)

type API struct {
	addr string
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}

	port := os.Getenv("SERVER_PORT")

	api := &API{
		addr: fmt.Sprintf(":%s", port),
	}

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    api.addr,
		Handler: mux,
	}
	log.Println("Hello from kj on the sever")
	log.Fatal(srv.ListenAndServe())

}
