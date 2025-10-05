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
	VersionDetails []struct {
		EncounterDetails []struct {
			Chance          int   `json:"chance"`
			ConditionValues []any `json:"condition_values"`
			MaxLevel        int   `json:"max_level"`
			Method          struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"method"`
			MinLevel int `json:"min_level"`
		} `json:"encounter_details"`
		MaxChance int `json:"max_chance"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"version_details"`
}

type LocationDetailsPokeAPIResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

var cache = pokecache.NewCache(10 * time.Minute)

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
	url := "https://pokeapi.co/api/v2/location-area/" + location
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
