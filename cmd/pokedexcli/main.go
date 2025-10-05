package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/MaelBriantin/pokedexcli/internal/pokeapi"
)

func main() {
	commandRegistry := getCommandRegistry()
	scanner := bufio.NewScanner(os.Stdin)
	config := &config{
		Next:    "https://pokeapi.co/api/v2/location-area",
		Pokedex: make(map[string]pokeapi.Pokemon),
	}
	for {
		fmt.Print("Pokedex > ")
		input := scanner.Scan()
		if !input {
			break
		}
		text := scanner.Text()
		cleanText := cleanInput(text)
		command := cleanText[0]
		if len(command) == 0 {
			continue
		}
		if cmd, ok := commandRegistry[command]; ok {
			if err := cmd.callback(config, cleanText[1:]); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
