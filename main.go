package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	godotenv.Load(".env")
	portString := os.Getenv("PORT")

	fmt.Println(portString)
}
