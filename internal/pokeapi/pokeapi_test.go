package pokeapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pokedexcli/internal/pokecache"
	"testing"
	"time"
)

func fakeAPIResponse() []byte {
	resp := PokeAPIResponse{
		Count: 1,
		Results: []Result{
			{Name: "test-location", URL: "https://pokeapi.co/api/v2/location-area/1/"},
		},
	}
	data, _ := json.Marshal(resp)
	return data
}

func TestGetLocationAreas_HTTP(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(fakeAPIResponse())
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
		w.Write(fakeAPIResponse())
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
