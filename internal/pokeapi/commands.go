package pokeapi

import (
	"encoding/json"
	"fmt"
	"os"
)

type CLICommand struct {
	name        string
	description string
	Callback    func(*Config) error
}

func MakeCommandRegistry(cfg *Config, client *PokeClient) map[string]CLICommand {
	command_registry := make(map[string]CLICommand)
	commandHelp := makeCommandHelp(cfg, command_registry)
	commandExit := makeCommandExit(cfg)
	commandMap := makeCommandMap(cfg, client)
	commandMapb := makeCommandMapb(cfg, client)

	command_registry["exit"] = CLICommand{
		name:        "exit",
		description: "Exit the Pokedex",
		Callback:    commandExit,
	}
	command_registry["help"] = CLICommand{
		name:        "help",
		description: "Prints a list of commands",
		Callback:    commandHelp,
	}
	command_registry["map"] = CLICommand{
		name:        "map",
		description: "Returns a list of next 20 map areas",
		Callback:    commandMap,
	}
	command_registry["mapb"] = CLICommand{
		name:        "mapb",
		description: "Returns a list of previous 20 map areas",
		Callback:    commandMapb,
	}

	return command_registry
}

func makeCommandMap(cfg *Config, client *PokeClient) func(*Config) error {
	return func(cfg *Config) error {
		if cfg.Next == "" || cfg.Next == "null" {
			fmt.Println("You are on the last page.")
			return nil
		}

		bytes, err := client.Get(cfg.Next)
		if err != nil {
			return fmt.Errorf("Error getting map data from pokeapi %w", err)
		}

		map_response := NamedAPIResourceList{}
		err = json.Unmarshal(bytes, &map_response)
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

func makeCommandMapb(cfg *Config, client *PokeClient) func(*Config) error {
	return func(cfg *Config) error {
		if cfg.Prev == "" || cfg.Prev == "null" {
			fmt.Println("You are on the first page.")
			return nil
		}

		bytes, err := client.Get(cfg.Prev)
		if err != nil {
			return fmt.Errorf("Error getting mapb data from pokeapi %w", err)
		}

		map_response := NamedAPIResourceList{}
		err = json.Unmarshal(bytes, &map_response)
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
