package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/crisp-coder/go-pokedex/internal/pokeapi"
)

type CLICommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	Next string
	Prev string
}

func makeCommandRegistry(cfg *Config) map[string]CLICommand {
	command_registry := make(map[string]CLICommand)
	commandHelp := makeCommandHelp(cfg, command_registry)
	commandExit := makeCommandExit(cfg)
	commandMap := makeCommandMap(cfg)
	commandMapb := makeCommandMapb(cfg)

	command_registry["exit"] = CLICommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
	command_registry["help"] = CLICommand{
		name:        "help",
		description: "Prints a list of commands.",
		callback:    commandHelp,
	}
	command_registry["map"] = CLICommand{
		name:        "map",
		description: "Returns a list of 20 map areas. Call again to get the next 20",
		callback:    commandMap,
	}
	command_registry["mapb"] = CLICommand{
		name:        "mapb",
		description: "Returns a list of previous 20 map areas.",
		callback:    commandMapb,
	}

	return command_registry
}

func makeCommandMap(cfg *Config) func(*Config) error {
	return func(cfg *Config) error {
		if cfg.Next == "" || cfg.Next == "null" {
			fmt.Println("You are on the last page.")
			return nil
		}

		resp, err := http.Get(cfg.Next)
		if err != nil {
			return fmt.Errorf("Error getting map data from pokeapi %w", err)
		}

		byteArray, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error getting map data from pokeapi %w", err)
		}

		map_response := pokeapi.NamedAPIResourceList{}
		err = json.Unmarshal(byteArray, &map_response)
		if err != nil {
			return fmt.Errorf("Error getting map data from pokeapi %w", err)
		}

		if cfg.Next != "" && cfg.Next != "null" {
			cfg.Next = map_response.Next
			cfg.Prev = map_response.Previous
		}

		for _, value := range map_response.Results {
			fmt.Printf("%v\n", value.Name)
		}

		return nil
	}
}

func makeCommandMapb(cfg *Config) func(*Config) error {
	return func(cfg *Config) error {
		if cfg.Prev == "" || cfg.Prev == "null" {
			fmt.Println("You are on the first page.")
			return nil
		}

		resp, err := http.Get(cfg.Prev)
		if err != nil {
			return fmt.Errorf("Error getting mapb data from pokeapi %w", err)
		}

		byteArray, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error getting mapb data from pokeapi %w", err)
		}

		map_response := pokeapi.NamedAPIResourceList{}
		err = json.Unmarshal(byteArray, &map_response)
		if err != nil {
			return fmt.Errorf("Error getting mapb data from pokeapi %w", err)
		}

		if cfg.Prev != "" && cfg.Prev != "null" {
			cfg.Next = map_response.Next
			cfg.Prev = map_response.Previous
		}

		for _, value := range map_response.Results {
			fmt.Printf("%v\n", value.Name)
		}

		return nil
	}
}

func makeCommandExit(cfg *Config) func(*Config) error {
	return func(cfg *Config) error {
		fmt.Print("Closing the Pokedex... Goodbye!\n")
		os.Exit(0)
		return nil
	}
}

func makeCommandHelp(cfg *Config, registry map[string]CLICommand) func(*Config) error {
	return func(cfg *Config) error {
		fmt.Print("Welcome to the Pokedex!\n")
		fmt.Print("Usage:\n\n")
		for key := range registry {
			fmt.Printf("%v: %v\n", registry[key].name, registry[key].description)
		}
		return nil
	}
}
