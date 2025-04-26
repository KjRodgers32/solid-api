package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lpernett/godotenv"
)

type API struct {
	addr string
}

func loggingMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		f(w, r)
	}
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

	mux.HandleFunc("/api/v1/users", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Fprint(w, "returning all users in db")
		case http.MethodPost:
			http.Error(w, "adding users not implemented yet", http.StatusNotImplemented)
		case http.MethodPut, http.MethodPatch:
			http.Error(w, "updating users not implemented yet", http.StatusNotImplemented)
		case http.MethodDelete:
			http.Error(w, "deleting users not implemented yet", http.StatusNotImplemented)
		default:
			http.Error(w, "bad request on users", http.StatusBadRequest)
		}
	}))

	mux.HandleFunc("/api/v1/users/{id}", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id := r.PathValue("id")
			fmt.Fprintf(w, "user id: %s", id)
		case http.MethodPost:
			id := r.PathValue("id")
			fmt.Fprintf(w, "added user %s to db", id)
		default:
			http.Error(w, "not implemented yet", http.StatusNotImplemented)
		}
	}))

	mux.HandleFunc("/api/v1", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello from root")
	}))

	log.Println("Serving on port", srv.Addr[1:])
	log.Fatal(srv.ListenAndServe())
}
