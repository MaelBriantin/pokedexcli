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

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Forms          []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	Height                 int    `json:"height"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Name                   string `json:"name"`
	Order                  int    `json:"order"`
	Sprites                struct {
		BackDefault      string `json:"back_default"`
		BackFemale       string `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  string `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      string `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale string `json:"front_shiny_female"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
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

func GetPokemonDetails(pokemon string) Pokemon {
	url := baseURL + "/pokemon/" + pokemon
	// Check cache first
	if data, found := cache.Get(url); found {
		response := Pokemon{}
		if err := json.Unmarshal(data, &response); err != nil {
			return Pokemon{}
		}
		return response
	}
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching pokemon details:", err)
		return Pokemon{}
	}
	data, err := io.ReadAll(res.Body)
	res.Body.Close()
	// Store in cache
	cache.Add(url, data)
	response := Pokemon{}
	if err := json.Unmarshal(data, &response); err != nil {
		return Pokemon{}
	}
	return response
}
