package main

type config struct {
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config) error
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
			description: "Displays a list of locations",
			callback:    nextLocationAreas,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous list of locations",
			callback:    previousLocationAreas,
		},
	}
}
