package pokeAPI

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Result struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokeAPIResponse struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  []Result `json:"results"`
}

func GetLocationAreas(url string) PokeAPIResponse {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching locations:", err)
		return PokeAPIResponse{}
	}
	data, err := io.ReadAll(res.Body)
	res.Body.Close()
	response := PokeAPIResponse{}
	error := json.Unmarshal(data, &response)
	if error != nil {
		fmt.Println("Error decoding JSON:", err)
		return PokeAPIResponse{}
	}
	return response
}
