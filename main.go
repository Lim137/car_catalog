package main

import (
	"database/sql"
	"fmt"
	_ "github.com/Lim137/car_catalog/docs"
	"github.com/Lim137/car_catalog/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type apiConfig struct {
	DB *database.Queries
}

// @title Car Catalog API
// @version 1.0
// @description This is an API for managing cars in a catalog.
// @BasePath /cars
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	portString := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(fmt.Sprintf("Can't connect to database: %v", err))
	}
	apiCfg := apiConfig{
		DB: database.New(conn),
	}
	log.Println("applying migrations...")
	cmd := exec.Command("goose", "-dir", "sql/schema", "postgres", dbURL, "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Couldn't apply migrations: %v", err)
	}
	log.Println("migrations have been successfully applied")
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	carRouter := chi.NewRouter()
	carRouter.Post("/", apiCfg.handlerCreateCars)
	carRouter.Delete("/", apiCfg.handlerDeleteCarById)
	carRouter.Put("/", apiCfg.handlerUpdateCarById)
	carRouter.Get("/", apiCfg.handlerGetCars)
	serverURL := os.Getenv("SERVER_URL")
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%v:%v/swagger/doc.json", serverURL, portString)),
	))
	router.Mount("/cars", carRouter)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Println("Server is running on port " + portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Couldn't start server: %v", err)
	}
}
