package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type CarInfo struct {
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

func getCarInfoFromApi(regNum string) (CarInfo, error) {
	params := url.Values{}
	params.Add("regNum", regNum)
	apiUrl := os.Getenv("API_URL") + params.Encode()
	response, err := http.Get(apiUrl)
	if err != nil {
		return CarInfo{}, fmt.Errorf("failed api request: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return CarInfo{}, fmt.Errorf("error %s: %s", response.Status, response.Body)
	}
	var car CarInfo
	err = json.NewDecoder(response.Body).Decode(&car)
	if err != nil {
		return CarInfo{}, fmt.Errorf("couldn't decode response: %v", err)
	}
	return car, nil
}
