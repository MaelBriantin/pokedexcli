package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
	"time"
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

var cache = pokecache.NewCache(10 * time.Minute)

func jsonDecode(data []byte, response *PokeAPIResponse) PokeAPIResponse {
	if err := json.Unmarshal(data, response); err != nil {
		return PokeAPIResponse{}
	}
	return *response
}

func GetLocationAreas(url string) PokeAPIResponse {
	// Check cache first
	if data, found := cache.Get(url); found {
		response := PokeAPIResponse{}
		return jsonDecode(data, &response)
	}

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching locations:", err)
		return PokeAPIResponse{}
	}
	data, err := io.ReadAll(res.Body)
	res.Body.Close()

	// Store in cache
	cache.Add(url, data)

	return jsonDecode(data, &PokeAPIResponse{})
}
