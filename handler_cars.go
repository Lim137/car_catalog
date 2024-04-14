package main

import (
	"encoding/json"
	"fmt"
	"github.com/Lim137/car_catalog/internal/database"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
	"time"
)

type CreateSuccessfully struct {
	RegNum string    `json:"regNum"`
	Status string    `json:"status"`
	ID     uuid.UUID `json:"id"`
}
type CreateError struct {
	RegNum string `json:"regNum"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type errRespond struct {
	Error string `json:"error"`
}

type carParameters struct {
	RegNum          string `json:"regNum"`
	Mark            string `json:"mark"`
	Model           string `json:"model"`
	Year            int    `json:"year"`
	OwnerName       string `json:"ownerName"`
	OwnerSurname    string `json:"ownerSurname"`
	OwnerPatronymic string `json:"ownerPatronymic"`
}

//@summary Delete a car by ID
//@description This endpoint deletes a car from the database by its ID in database.
//@tags cars
//@produce json
//@param carId query string true "CarID"
//@success 200 {object} MessageResponse "Car was successfully deleted"
//@failure 400 {object} errRespond "Error parsing request"
//@failure 500 {object} errRespond "Error deleting car from DB"
//@Router / [delete]
func (apiCfg *apiConfig) handlerDeleteCarById(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	queryParams := url.Query()
	carIdStr := queryParams.Get("carId")
	carId, err := uuid.Parse(carIdStr)
	if err != nil {
		log.Printf("Error parsing car ID: %v\nURL: %v", err, url)
		respondWithJSON(w, 400, errRespond{Error: fmt.Sprintf("Couldn't parse car ID: %v", err)})
		return
	}

	err = apiCfg.DB.DeleteCarById(r.Context(), carId)
	if err != nil {
		log.Printf("Error deleting car from DB: %v", err)
		respondWithJSON(w, 500, errRespond{Error: fmt.Sprintf("Couldn't delete car from DB: %v", err)})
		return
	}
	respondWithJSON(w, 200, MessageResponse{Message: "Car was successfully deleted"})
}

// @summary Create a new car
// @description This endpoint creates a new car in the database. It takes an array of car registration numbers, makes API requests to fetch data about each car, and then adds them to the database.
// @tags cars
// @accept json
// @produce json
// @param request body []string true "Array of car registration numbers"
// @success 200 {array} CreateSuccessfully "An array containing information about each successfully added car"
// @failure 500 {array} CreateError "An array containing errors for cars that couldn't be added to the database"
// @failure 400 {object} errRespond "Error parsing request"
// @Router / [post]
func (apiCfg *apiConfig) handlerCreateCars(w http.ResponseWriter, r *http.Request) {
	type parametersForCreateCars struct {
		RegNums []string `json:"regNums"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parametersForCreateCars{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error parsing JSON: %v", err)
		respondWithJSON(w, 400, errRespond{Error: fmt.Sprintf("Couldn't parse JSON: %v", err)})
		return
	}
	var result []interface{}
	for _, value := range params.RegNums {
		carInfoFromApi, err := getCarInfoFromApi(value)
		if err != nil {
			log.Printf("Error getting car info from API: %v", err)
			result = append(result, CreateError{
				RegNum: value,
				Status: "failed",
				Error:  fmt.Sprintf("Couldn't get car info from API: %v", err),
			})
			continue
		}

		carIdInDB, err := apiCfg.DB.CreateCar(r.Context(), database.CreateCarParams{
			ID:              uuid.New(),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			RegNum:          carInfoFromApi.RegNum,
			Mark:            carInfoFromApi.Mark,
			Model:           carInfoFromApi.Model,
			Year:            int32(carInfoFromApi.Year),
			OwnerName:       carInfoFromApi.Owner.Name,
			OwnerSurname:    carInfoFromApi.Owner.Surname,
			OwnerPatronymic: carInfoFromApi.Owner.Patronymic,
		})
		if err != nil {
			log.Printf("Error creating car in DB: %v", err)
			result = append(result, CreateError{
				RegNum: value,
				Status: "failed",
				Error:  fmt.Sprintf("Couldn't create car in DB: %v", err),
			})
			continue
		}
		result = append(result, CreateSuccessfully{
			RegNum: carInfoFromApi.RegNum,
			Status: "success",
			ID:     carIdInDB,
		})
	}
	hasErrors := false
	for _, res := range result {
		if _, ok := res.(CreateError); ok {
			hasErrors = true
			break
		}
	}
	if hasErrors {
		respondWithJSON(w, 500, result)
	} else {
		respondWithJSON(w, 200, result)
	}
}

// @summary Update a car by ID
// @description This endpoint updates a car in the database by its ID.
// @tags cars
// @accept json
// @produce json
// @param carId query string true "Car ID"
// @param request body carParameters true "Car parameters that need to be updated"
// @success 200 {object} database.Car "Updated car information"
// @failure 400 {object} errRespond "Error parsing car ID"
// @failure 400 {object} errRespond "Error parsing JSON"
// @failure 500 {object} errRespond "Error updating car in DB"
// @Router / [put]
func (apiCfg *apiConfig) handlerUpdateCarById(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	queryParams := url.Query()
	carIdStr := queryParams.Get("carId")
	carId, err := uuid.Parse(carIdStr)
	if err != nil {
		log.Printf("Error parsing car ID: %v\nURL: %v", err, url)
		respondWithJSON(w, 400, errRespond{Error: fmt.Sprintf("Couldn't parse car ID: %v", err)})
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := carParameters{
		RegNum:          "",
		Mark:            "",
		Model:           "",
		Year:            -1,
		OwnerName:       "",
		OwnerSurname:    "",
		OwnerPatronymic: "",
	}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error parsing JSON: %v", err)
		respondWithJSON(w, 400, errRespond{Error: fmt.Sprintf("Couldn't parse JSON: %v", err)})
		return
	}
	updatedCarInfo, err := apiCfg.DB.UpdateCarById(r.Context(), database.UpdateCarByIdParams{
		ID:      carId,
		Column2: params.RegNum,
		Column3: params.Mark,
		Column4: params.Model,
		Column5: int32(params.Year),
		Column6: params.OwnerName,
		Column7: params.OwnerSurname,
		Column8: params.OwnerPatronymic,
	})
	if err != nil {
		log.Printf("Error updating car in DB: %v", err)
		respondWithJSON(w, 500, errRespond{Error: fmt.Sprintf("Couldn't update car in DB: %v", err)})
		return
	}

	respondWithJSON(w, 200, updatedCarInfo)
}

// @summary Get cars
// @description This endpoint retrieves cars from the catalog based on specified parameters.
// @tags cars
// @produce json
// @param regNum query string false "Car registration number"
// @param mark query string false "Car mark"
// @param model query string false "Car model"
// @param year query string false "Car year (It is expected that it will be possible to convert to integer)"
// @param ownerName query string false "Owner's name"
// @param ownerSurname query string false "Owner's surname"
// @param ownerPatronymic query string false "Owner's patronymic"
// @param pageSize query string false "Page size"
// @param page query string false "Page number"
// @success 200 {array} database.Car "List of cars"
// @success 404 {object} MessageResponse "Cars with such parameters not found"
// @failure 500 {object} errRespond "Error parsing year"
// @failure 500 {object} errRespond "Error parsing page"
// @failure 500 {object} errRespond "Error parsing page size"
// @failure 500 {object} errRespond "Error getting cars from DB"
// @Router / [get]
func (apiCfg *apiConfig) handlerGetCars(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	queryParams := url.Query()
	regNum := queryParams.Get("regNum")
	mark := queryParams.Get("mark")
	model := queryParams.Get("model")
	yearStr := queryParams.Get("year")
	ownerName := queryParams.Get("ownerName")
	ownerSurname := queryParams.Get("ownerSurname")
	ownerPatronymic := queryParams.Get("ownerPatronymic")
	pageSizeStr := queryParams.Get("pageSize")
	pageStr := queryParams.Get("page")

	var year, page, pageSize int
	var err error
	if yearStr == "" {
		year = 0
	} else {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			log.Printf("Error parsing year: %v\nyearStr: %v", err, yearStr)
			respondWithJSON(w, 500, errRespond{Error: fmt.Sprintf("Couldn't parse year: %v", err)})
			return
		}
	}
	if pageStr == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			log.Printf("Error parsing page: %v\npageStr: %v", err, pageStr)
			respondWithJSON(w, 500, errRespond{Error: fmt.Sprintf("Couldn't parse page: %v", err)})
			return
		}
	}
	if pageSizeStr == "" {
		pageSize = -1
	} else {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			log.Printf("Error parsing page size: %v\npageSizeStr: %v", err, pageSizeStr)
			respondWithJSON(w, 500, errRespond{Error: fmt.Sprintf("Couldn't parse page size: %v", err)})
			return
		}
	}
	cars, err := apiCfg.DB.GetCars(r.Context(), database.GetCarsParams{
		RegNum:          regNum,
		Mark:            mark,
		Model:           model,
		Year:            int32(year),
		OwnerName:       ownerName,
		OwnerSurname:    ownerSurname,
		OwnerPatronymic: ownerPatronymic,
		Column8:         int32(pageSize),
		Column9:         int32(page),
	})
	if err != nil {
		log.Printf("Error getting cars from DB: %v", err)
		respondWithJSON(w, 500, errRespond{Error: fmt.Sprintf("Couldn't get cars from DB: %v", err)})
		return
	}
	if len(cars) == 0 {
		respondWithJSON(w, 404, MessageResponse{Message: "Cars with such parameters not found"})
		return
	}
	respondWithJSON(w, 200, cars)
}
