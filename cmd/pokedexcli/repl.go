package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/MaelBriantin/pokedexcli/internal/pokeapi"
)

type config struct {
	Next     string
	Previous string
	Pokedex  map[string]pokeapi.Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, args []string) error
}

func getCommandRegistry() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    help,
		},
		"map": {
			name:        "map",
			description: "Displays the next list of locations",
			callback:    nextLocationAreas,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous list of locations",
			callback:    previousLocationAreas,
		},
		"explore": {
			name:        "explore",
			description: "List of all the Pokémon located in a specific location",
			callback:    exploreLocation,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokémon (not implemented yet)",
			callback:    catch,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all the Pokémon you've caught",
			callback:    seePokedex,
		},
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func commandExit(cfg *config, _ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func help(cfg *config, _ []string) error {
	commands := getCommandRegistry()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func nextLocationAreas(cfg *config, _ []string) error {
	pokeAPIResponse := pokeapi.GetLocationAreas(cfg.Next)
	cfg.Next = pokeAPIResponse.Next
	cfg.Previous = pokeAPIResponse.Previous
	for _, loc := range pokeAPIResponse.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func previousLocationAreas(cfg *config, _ []string) error {
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

func exploreLocation(_ *config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Please provide a location name")
		return nil
	}
	if len(args) > 1 {
		fmt.Println("Please provide only one location name")
		return nil
	}
	location := args[0]
	pokeAPIResponse := pokeapi.GetLocationDetails(location)
	if len(pokeAPIResponse.PokemonEncounters) == 0 {
		fmt.Println("No Pokemon found in this location")
		return nil
	}
	fmt.Printf("Exploring %s...\n", location)
	fmt.Println("Found Pokemon:")
	for _, encounter := range pokeAPIResponse.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func catch(cfg *config, args []string) error {
	if len(args) > 1 {
		fmt.Println("You can only catch one Pokemon at a time... Don't be greedy!")
		return nil
	}
	if len(args) < 1 {
		fmt.Println("Please provide a Pokemon name to catch")
		return nil
	}
	pokedex := cfg.Pokedex
	pokemon := args[0]
	if _, exists := pokedex[pokemon]; exists {
		fmt.Printf("You already caught a %s!\n", pokemon)
		return nil
	}
	pokeAPIResponse := pokeapi.GetPokemonDetails(pokemon)
	if pokeAPIResponse.Name == "" {
		fmt.Printf("Could not find a Pokemon named %s\n", pokemon)
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokeAPIResponse.Name)
	rndNum := rand.Intn(100)
	successThreshold := pokeAPIResponse.BaseExperience / 2
	if rndNum < successThreshold {
		fmt.Printf("%s escaped!\n", pokeAPIResponse.Name)
		return nil
	}
	fmt.Printf("%s was caught!\n", pokeAPIResponse.Name)
	pokedex[pokeAPIResponse.Name] = pokeAPIResponse
	return nil
}

func seePokedex(cfg *config, _ []string) error {
	pokedex := cfg.Pokedex
	if len(pokedex) == 0 {
		fmt.Println("Your Pokedex is empty. Go catch some Pokemon!")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for name, pokemon := range pokedex {
		fmt.Printf("- %s (Base Experience: %d, Height: %d, Weight: %d)\n", name, pokemon.BaseExperience, pokemon.Height, pokemon.Weight)
	}
	return nil
}
