package main

import (
	"fmt"
	"os"
	"strings"

	"pokedexcli/internal/pokeapi"
)

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func help(cfg *config) error {
	commands := getCommandRegistry()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func nextLocationAreas(cfg *config) error {
	pokeAPIResponse := pokeapi.GetLocationAreas(cfg.Next)
	cfg.Next = pokeAPIResponse.Next
	cfg.Previous = pokeAPIResponse.Previous
	for _, loc := range pokeAPIResponse.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func previousLocationAreas(cfg *config) error {
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	pokeAPIResponse := pokeapi.GetLocationAreas(cfg.Previous)
	cfg.Next = pokeAPIResponse.Next
	cfg.Previous = pokeAPIResponse.Previous
	for _, loc := range pokeAPIResponse.Results {
		fmt.Println(loc.Name)
	}
	return nil
}
