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
	godotenv.Load(".env")
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

	// @summary Create a new car
	// @description This endpoint creates a new car in the database. It takes an array of car registration numbers, makes API requests to fetch data about each car, and then adds them to the database.
	// @tags cars
	// @accept json
	// @produce json
	// @param request body []string true "Array of car registration numbers"
	// @success 200 {array} interface{} "An array containing information about each successfully added car"
	// @failure 500 {array} CreateError "An array containing errors for cars that couldn't be added to the database"
	carRouter.Post("/", apiCfg.handlerCreateCars)

	// @summary Delete a car by ID
	// @description This endpoint deletes a car from the database by its ID in database.
	// @tags cars
	// @produce json
	// @param carId query string true "CarID"
	// @success 200 {object} MessageResponse "Car was successfully deleted"
	// @failure 400 {object} errRespond "Error parsing request or car not found"
	// @failure 500 {object} errRespond "Error deleting car from DB"
	carRouter.Delete("/", apiCfg.handlerDeleteCarById)

	// @summary Update a car by ID
	// @description This endpoint updates a car in the database by its ID.
	// @tags cars
	// @accept json
	// @produce json
	// @param carId query string true "Car ID"
	// @param request body parameters true "Car parameters that need to be updated"
	// @success 200 {object} database.Car "Updated car information"
	// @failure 400 {object} errRespond "Error parsing car ID"
	// @failure 400 {object} errRespond "Error parsing JSON"
	// @failure 500 {object} errRespond "Error updating car in DB"
	carRouter.Put("/", apiCfg.handlerUpdateCarById)

	//	@summary Get cars
	//	@description This endpoint retrieves cars from the catalog based on specified parameters.
	//	@tags cars
	//	@produce json
	//	@param regNum query string false "Car registration number"
	//	@param mark query string false "Car mark"
	//	@param model query string false "Car model"
	//	@param year query string false "Car year (It is expected that it will be possible to convert to integer)"
	// 	@param ownerName query string false "Owner's name"
	// 	@param ownerSurname query string false "Owner's surname"
	// 	@param ownerPatronymic query string false "Owner's patronymic"
	// 	@param pageSize query string false "Page size"
	// 	@param page query string false "Page number"
	// 	@success 200 {array} database.Car "List of cars"
	// 	@success 404 {object} MessageResponse "Cars with such parameters not found"
	// 	@failure 500 {object} errRespond "Error parsing year"
	// 	@failure 500 {object} errRespond "Error parsing page"
	// 	@failure 500 {object} errRespond "Error parsing page size"
	// 	@failure 500 {object} errRespond "Error getting cars from DB"
	carRouter.Get("/", apiCfg.handlerGetCars)
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
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
