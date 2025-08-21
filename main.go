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
		API: "https://pokeapi.co/api/v2/",
		Pokedex: pokeapi.TempPokedex{
			KnownPokemon: make(map[string]pokeapi.Pokemon),
		},
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

		if len(cleanedInput) == 2 {
			cfg.CommandParams = cleanedInput[1:]
		}

		// Run command
		err := pokedex_command.Callback()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
	}
}
