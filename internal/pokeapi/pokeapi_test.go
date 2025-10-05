package pokeapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MaelBriantin/pokedexcli/internal/pokecache"
)

func fakeLocationAPIResponse() []byte {
	resp := LocationPokeAPIResponse{
		Count:    1,
		Next:     "",
		Previous: "",
		Results: []Result{
			{Name: "test-location", URL: "https://pokeapi.co/api/v2/location-area/1/"},
		},
	}
	data, _ := json.Marshal(resp)
	return data
}

func fakeLocationDetailsAPIResponse() []byte {
	resp := LocationDetailsPokeAPIResponse{
		PokemonEncounters: []PokemonEncounter{
			{Pokemon: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "test-pokemon", URL: "https://pokeapi.co/api/v2/pokemon/1/"},
			},
		},
	}
	data, _ := json.Marshal(resp)
	return data
}

func TestGetLocationAreas_HTTP(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(fakeLocationAPIResponse())
	}))
	defer ts.Close()

	cache = pokecache.NewCache(10 * time.Minute)

	resp := GetLocationAreas(ts.URL)
	if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "test-location" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestGetLocationAreas_Cache(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(fakeLocationAPIResponse())
	}))
	defer ts.Close()

	cache = pokecache.NewCache(10 * time.Minute)

	resp1 := GetLocationAreas(ts.URL)
	if resp1.Count != 1 || len(resp1.Results) != 1 || resp1.Results[0].Name != "test-location" {
		t.Errorf("unexpected first response: %+v", resp1)
	}

	ts.Close()

	resp2 := GetLocationAreas(ts.URL)
	if resp2.Count != 1 || len(resp2.Results) != 1 || resp2.Results[0].Name != "test-location" {
		t.Errorf("unexpected response from cache: %+v", resp2)
	}
}

func TestGetLocationDetails_HTTP(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(fakeLocationDetailsAPIResponse())
	}))
	defer ts.Close()

	originalBaseURL := baseURL
	baseURL = ts.URL
	defer func() { baseURL = originalBaseURL }()

	cache = pokecache.NewCache(10 * time.Minute)

	resp := GetLocationDetails("test-location")

	if len(resp.PokemonEncounters) != 1 {
		t.Errorf("expected 1 pokemon encounter, got %d", len(resp.PokemonEncounters))
	}

	if len(resp.PokemonEncounters) > 0 && resp.PokemonEncounters[0].Pokemon.Name != "test-pokemon" {
		t.Errorf("expected pokemon name 'test-pokemon', got '%s'", resp.PokemonEncounters[0].Pokemon.Name)
	}
}

func TestGetLocationDetails_Cache(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(fakeLocationDetailsAPIResponse())
	}))
	defer ts.Close()

	originalBaseURL := baseURL
	baseURL = ts.URL
	defer func() { baseURL = originalBaseURL }()

	cache = pokecache.NewCache(10 * time.Minute)

	resp1 := GetLocationDetails("test-location")
	if len(resp1.PokemonEncounters) != 1 || resp1.PokemonEncounters[0].Pokemon.Name != "test-pokemon" {
		t.Errorf("unexpected first response: %+v", resp1)
	}

	ts.Close()

	resp2 := GetLocationDetails("test-location")
	if len(resp2.PokemonEncounters) != 1 || resp2.PokemonEncounters[0].Pokemon.Name != "test-pokemon" {
		t.Errorf("unexpected response from cache: %+v", resp2)
	}
}
