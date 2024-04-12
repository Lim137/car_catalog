package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Lim137/car_catalog/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strconv"
	"time"
)

type CreateResponse struct {
	RegNum string    `json:"regNum"`
	ID     uuid.UUID `json:"id"`
}

func (apiCfg *apiConfig) handlerDeleteCarById(w http.ResponseWriter, r *http.Request) {
	carIdStr := chi.URLParam(r, "carId")
	carId, err := uuid.Parse(carIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse car ID: %v", err))
		return
	}
	err = apiCfg.DB.DeleteCarById(r.Context(), carId)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete car from DB: %v", err))
		return
	}
	respondWithJSON(w, 200, "Car was successfully deleted")
}

func (apiCfg *apiConfig) handlerCreateCars(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		RegNums []string `json:"regNums"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	var addedCarsIds []CreateResponse
	for _, value := range params.RegNums {
		carInfoFromApi, err := getCarInfoFromApi(value)
		if err != nil {
			respondWithError(w, 500, fmt.Sprintf("Error getting car info from API: %v", err))
			continue
		}
		ownerPatronymic := sql.NullString{}
		if carInfoFromApi.Owner.Patronymic != "" {
			ownerPatronymic.String = carInfoFromApi.Owner.Patronymic
			ownerPatronymic.Valid = true

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
			OwnerPatronymic: ownerPatronymic,
		})
		if err != nil {
			respondWithError(w, 500, fmt.Sprintf("Error creating car in DB: %v", err))
			continue
		}
		addedCarsIds = append(addedCarsIds, CreateResponse{
			RegNum: carInfoFromApi.RegNum,
			ID:     carIdInDB,
		})
	}
	respondWithJSON(w, 200, addedCarsIds)
}

func (apiCfg *apiConfig) handlerUpdateCarById(w http.ResponseWriter, r *http.Request) {
	carIdStr := chi.URLParam(r, "carId")
	carId, err := uuid.Parse(carIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse car ID: %v", err))
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't read request body: %v", err))
		return
	}
	var requestBody map[string]interface{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't unmarshal request body: %v", err))
		return
	}
	for key, value := range requestBody {
		if key == "regNum" {
			err := apiCfg.DB.UpdateRegNumById(r.Context(), database.UpdateRegNumByIdParams{ID: carId, RegNum: value.(string)})
			if err != nil {
				respondWithError(w, 500, fmt.Sprintf("Couldn't update car regNum in DB: %v", err))
			}
			continue
		}
		if key == "mark" {
			err := apiCfg.DB.UpdateMarkById(r.Context(), database.UpdateMarkByIdParams{ID: carId, Mark: value.(string)})
			if err != nil {
				respondWithError(w, 500, fmt.Sprintf("Couldn't update car mark in DB: %v", err))
			}
			continue
		}
		if key == "model" {
			err := apiCfg.DB.UpdateModelById(r.Context(), database.UpdateModelByIdParams{ID: carId, Model: value.(string)})
			if err != nil {
				respondWithError(w, 500, fmt.Sprintf("Couldn't update car model in DB: %v", err))
			}
			continue
		}
		if key == "year" {
			err := apiCfg.DB.UpdateYearById(r.Context(), database.UpdateYearByIdParams{ID: carId, Year: int32(value.(float64))})
			if err != nil {
				respondWithError(w, 500, fmt.Sprintf("Couldn't update the year of the car in DB: %v", err))
			}
			continue
		}
		if key == "ownerName" {
			err := apiCfg.DB.UpdateOwnerNameById(r.Context(), database.UpdateOwnerNameByIdParams{ID: carId, OwnerName: value.(string)})
			if err != nil {
				respondWithError(w, 500, fmt.Sprintf("Couldn't update car owner name in DB: %v", err))
			}
			continue
		}
		if key == "ownerSurname" {
			err := apiCfg.DB.UpdateOwnerSurnameById(r.Context(), database.UpdateOwnerSurnameByIdParams{ID: carId, OwnerSurname: value.(string)})
			if err != nil {
				respondWithError(w, 500, fmt.Sprintf("Couldn't update car owner surname in DB: %v", err))
			}
			continue
		}
		if key == "ownerPatronymic" {
			ownerPatronymic := sql.NullString{}
			if value != "" {
				ownerPatronymic.String = value.(string)
				ownerPatronymic.Valid = true

			}
			err := apiCfg.DB.UpdateOwnerPatronymicById(r.Context(), database.UpdateOwnerPatronymicByIdParams{ID: carId, OwnerPatronymic: ownerPatronymic})
			if err != nil {
				respondWithError(w, 500, fmt.Sprintf("Couldn't update car owner patronymic in DB: %v", err))
			}
			continue
		}
	}
	respondWithJSON(w, 200, "Car was successfully updated")
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
	ownerPatronymicStr := queryParams.Get("ownerPatronymic")
	pageSizeStr := queryParams.Get("pageSize")
	pageStr := queryParams.Get("page")
	ownerPatronymic := sql.NullString{}
	if ownerPatronymicStr != "" {
		ownerPatronymic.String = ownerPatronymicStr
		ownerPatronymic.Valid = true

	}
	var year, page, pageSize int
	var err error
	if yearStr == "" {
		year = 0
	} else {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			respondWithError(w, 500, fmt.Sprintf("Couldn't parse year: %v", err))
			return
		}
	}
	if pageStr == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			respondWithError(w, 500, fmt.Sprintf("Couldn't parse page: %v", err))
			return
		}
	}
	if pageSizeStr == "" {
		pageSize = -1
	} else {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			respondWithError(w, 500, fmt.Sprintf("Couldn't parse page: %v", err))
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
		respondWithError(w, 500, fmt.Sprintf("Couldn't get cars from DB: %v", err))
		return
	}
	if len(cars) == 0 {
		respondWithJSON(w, 200, "Cars with such parameters not found")
		return
	}
	respondWithJSON(w, 200, cars)
}
