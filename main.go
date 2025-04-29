package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"

	"github.com/lpernett/godotenv"

	"github.com/KjRodgers32/solid-api/db"
)

type API struct {
	addr string
}

type DBConfig struct {
	Name             string
	User             string
	Pass             string
	Port             string
	Host             string
	ConnectionString string
}

type UserJSON struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UsersData struct {
	Status string    `json:"status"`
	Data   []db.User `json:"data"`
}

func loadDBConfig() DBConfig {
	return DBConfig{
		Name: os.Getenv("DB"),
		User: os.Getenv("DB_USER"),
		Pass: os.Getenv("DB_PASS"),
		Port: os.Getenv("DB_PORT"),
		Host: os.Getenv("DB_HOST"),
	}
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

	dbConfig := loadDBConfig()
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbConfig.User, dbConfig.Pass, dbConfig.Host, dbConfig.Port, dbConfig.Name)
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		log.Fatal("error connecting to database:", err)
	}

	defer conn.Close(ctx)

	queries := db.New(conn)

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
			users, err := queries.ListUsers(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			responseData := UsersData{
				Status: "OK",
				Data:   users,
			}

			if err = json.NewEncoder(w).Encode(responseData); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case http.MethodPost:
			var user db.User
			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			addedUser, err := queries.CreateUser(ctx, user.Name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			responseData := UserJSON{
				ID:   addedUser.ID,
				Name: addedUser.Name,
			}

			if err = json.NewEncoder(w).Encode(responseData); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		default:
			http.Error(w, "bad request on users", http.StatusBadRequest)
		}
	}))

	mux.HandleFunc("/api/v1/users/{id}", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			user, err := queries.GetUser(ctx, int64(id))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			data := UserJSON{
				ID:   user.ID,
				Name: user.Name,
			}

			if err = json.NewEncoder(w).Encode(data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		case http.MethodPut, http.MethodPatch:
			id, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			user, err := queries.GetUser(ctx, int64(id))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			defer r.Body.Close()
			var respData UserJSON
			if err = json.NewDecoder(r.Body).Decode(&respData); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			updateData := db.UpdateUserParams{
				ID:   user.ID,
				Name: respData.Name,
			}

			data := db.UpdateUserParams(updateData)
			if err = queries.UpdateUser(ctx, data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			w.WriteHeader(http.StatusCreated)
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
