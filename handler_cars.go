package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"net/http"
)

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
