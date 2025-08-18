package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/crisp-coder/go-pokedex/internal/pokeapi"
)

func main() {
	// Load objects
	cfg := pokeapi.Config{
		Next: "https://pokeapi.co/api/v2/location-area/",
		Prev: "",
	}
	pokeClient := pokeapi.NewPokeClient(time.Second * 10)
	registry := pokeapi.MakeCommandRegistry(&cfg, pokeClient)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Prompt for input
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			fmt.Println("Scanner error reading input.")
			return
		}

		// Read command from terminal
		line := scanner.Text()
		cleanedInput := cleanInput(line)
		if len(cleanedInput) == 0 {
			fmt.Println("Input was empty.")
			continue
		}

		// Check if command exists
		command := cleanedInput[0]
		pokedex_command, ok := registry[command]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if pokedex_command.Name == "explore" {
			if len(cleanedInput) == 2 {
				cfg.ExploreTarget = "https://pokeapi.co/api/v2/location-area/" + cleanedInput[1]
			} else {
				fmt.Println("Missing location area id or name.")
				continue
			}
		}

		// Run command
		err := pokedex_command.Callback(&cfg)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
	}
}
