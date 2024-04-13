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

func (apiCfg *apiConfig) handlerDeleteCarById(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	queryParams := url.Query()
	carIdStr := queryParams.Get("carId")
	carId, err := uuid.Parse(carIdStr)
	if err != nil {
		log.Printf("Error parsing car ID: %v\nURL: %v", err, url)
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse car ID: %v", err))
		return
	}

	err = apiCfg.DB.DeleteCarById(r.Context(), carId)
	if err != nil {
		log.Printf("Error deleting car from DB: %v", err)
		respondWithError(w, 500, fmt.Sprintf("Couldn't delete car from DB: %v", err))
		return
	}
	respondWithJSON(w, 200, MessageResponse{Message: "Car was successfully deleted"})
}

func (apiCfg *apiConfig) handlerCreateCars(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		RegNums []string `json:"regNums"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error parsing JSON: %v", err)
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse JSON: %v", err))
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

func (apiCfg *apiConfig) handlerUpdateCarById(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	queryParams := url.Query()
	carIdStr := queryParams.Get("carId")
	carId, err := uuid.Parse(carIdStr)
	if err != nil {
		log.Printf("Error parsing car ID: %v\nURL: %v", err, url)
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse car ID: %v", err))
		return
	}
	type parameters struct {
		RegNum          string `json:"regNum"`
		Mark            string `json:"mark"`
		Model           string `json:"model"`
		Year            int    `json:"year"`
		OwnerName       string `json:"ownerName"`
		OwnerSurname    string `json:"ownerSurname"`
		OwnerPatronymic string `json:"ownerPatronymic"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{
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
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse JSON: %v", err))
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
		respondWithError(w, 500, fmt.Sprintf("Couldn't update car in DB: %v", err))
		return
	}

	respondWithJSON(w, 200, updatedCarInfo)
}

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
			respondWithError(w, 500, fmt.Sprintf("Couldn't parse year: %v", err))
			return
		}
	}
	if pageStr == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			log.Printf("Error parsing page: %v\npageStr: %v", err, pageStr)
			respondWithError(w, 500, fmt.Sprintf("Couldn't parse page: %v", err))
			return
		}
	}
	if pageSizeStr == "" {
		pageSize = -1
	} else {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			log.Printf("Error parsing page size: %v\npageSizeStr: %v", err, pageSizeStr)
			respondWithError(w, 500, fmt.Sprintf("Couldn't parse page size: %v", err))
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
		respondWithError(w, 500, fmt.Sprintf("Couldn't get cars from DB: %v", err))
		return
	}
	if len(cars) == 0 {
		respondWithJSON(w, 404, MessageResponse{Message: "Cars with such parameters not found"})
		return
	}
	respondWithJSON(w, 200, cars)
}
