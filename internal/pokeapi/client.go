package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/MaelBriantin/pokedexcli/internal/pokecache"
)

type Result struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationPokeAPIResponse struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  []Result `json:"results"`
}

type PokemonEncounter struct {
	Pokemon struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokemon"`
}

type LocationDetailsPokeAPIResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

var cache = pokecache.NewCache(10 * time.Minute)
var baseURL = "https://pokeapi.co/api/v2" // Ajouter cette variable

func locationJsonDecode(data []byte, response *LocationPokeAPIResponse) LocationPokeAPIResponse {
	if err := json.Unmarshal(data, response); err != nil {
		return LocationPokeAPIResponse{}
	}
	return *response
}

func encounterJsonDecode(data []byte, response *LocationDetailsPokeAPIResponse) LocationDetailsPokeAPIResponse {
	if err := json.Unmarshal(data, response); err != nil {
		return LocationDetailsPokeAPIResponse{}
	}
	return *response
}

func GetLocationAreas(url string) LocationPokeAPIResponse {
	// Check cache first
	if data, found := cache.Get(url); found {
		response := LocationPokeAPIResponse{}
		return locationJsonDecode(data, &response)
	}
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching locations:", err)
		return LocationPokeAPIResponse{}
	}
	data, err := io.ReadAll(res.Body)
	res.Body.Close()
	// Store in cache
	cache.Add(url, data)
	return locationJsonDecode(data, &LocationPokeAPIResponse{})
}

func GetLocationDetails(location string) LocationDetailsPokeAPIResponse {
	url := baseURL + "/location-area/" + location
	// Check cache first
	if data, found := cache.Get(url); found {
		response := LocationDetailsPokeAPIResponse{}
		return encounterJsonDecode(data, &response)
	}
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching location details:", err)
		return LocationDetailsPokeAPIResponse{}
	}
	data, err := io.ReadAll(res.Body)
	res.Body.Close()
	// Store in cache
	cache.Add(url, data)
	return encounterJsonDecode(data, &LocationDetailsPokeAPIResponse{})
}
