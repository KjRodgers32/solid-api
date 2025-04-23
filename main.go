package main

import (
	"fmt"
	"github.com/lpernett/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}

	port := os.Getenv("SERVER_PORT")
	fmt.Println(port)

}
