package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Lim137/car_catalog/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"net/http"
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
