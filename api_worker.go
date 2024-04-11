package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Car struct {
	RegNum string `json:"regNum"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Year   int    `json:"year,omitempty"`
	Owner  People `json:"owner"`
}

type People struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

func getCarInfoFromApi(regNum string) Car {
	params := url.Values{}
	params.Add("regNum", regNum)
	apiUrl := os.Getenv("API_URL") + params.Encode()
	response, err := http.Get(apiUrl)
	if err != nil {
		log.Println("Failed api request:", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Println("Error", response.Status+":", response.Body)
		return Car{}
	}
	var car Car
	err = json.NewDecoder(response.Body).Decode(&car)
	if err != nil {
		log.Println("Couldn't decode response:", err)
		return Car{}
	}
	fmt.Println(car)
	return car
}
