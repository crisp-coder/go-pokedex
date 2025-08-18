package pokeapi

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Next          string
	Prev          string
	ExploreTarget string
}

type CLICommand struct {
	Name        string
	Description string
	Callback    func(*Config) error
}

func MakeCommandRegistry(cfg *Config, client *PokeClient) map[string]CLICommand {
	command_registry := make(map[string]CLICommand)
	commandHelp := makeCommandHelp(command_registry)
	commandExit := makeCommandExit()
	commandMap := makeCommandMap(client)
	commandMapb := makeCommandMapb(client)
	commandExplore := makeCommandExplore(client)

	command_registry["exit"] = CLICommand{
		Name:        "exit",
		Description: "Exit the Pokedex",
		Callback:    commandExit,
	}
	command_registry["help"] = CLICommand{
		Name:        "help",
		Description: "Prints a list of commands",
		Callback:    commandHelp,
	}
	command_registry["map"] = CLICommand{
		Name:        "map",
		Description: "Returns a list of next 20 map areas",
		Callback:    commandMap,
	}
	command_registry["mapb"] = CLICommand{
		Name:        "mapb",
		Description: "Returns a list of previous 20 map areas",
		Callback:    commandMapb,
	}
	command_registry["explore"] = CLICommand{
		Name:        "explore",
		Description: "Asks for area to explore, then queries pokemon in area.",
		Callback:    commandExplore,
	}

	return command_registry
}

func makeCommandExplore(client *PokeClient) func(*Config) error {
	return func(cfg *Config) error {
		bytes, err := client.Get(cfg.ExploreTarget)
		if err != nil {
			return fmt.Errorf("Error getting explore area data: %w", err)
		}

		explore_response := LocationArea{}
		err = json.Unmarshal(bytes, &explore_response)
		if err != nil {
			return fmt.Errorf("Error getting explore area data: %w", err)
		}

		for _, pokemon_encounter := range explore_response.Pokemon_encounters {
			fmt.Printf("%v\n", pokemon_encounter.Pokemon.Name)
		}

		return nil
	}
}

func makeCommandMap(client *PokeClient) func(*Config) error {
	return func(cfg *Config) error {
		if cfg.Next == "" || cfg.Next == "null" {
			fmt.Println("You are on the last page.")
			return nil
		}

		bytes, err := client.Get(cfg.Next)
		if err != nil {
			return fmt.Errorf("Error getting map data from pokeapi: %w", err)
		}

		map_response := NamedAPIResourceList{}
		err = json.Unmarshal(bytes, &map_response)
		if err != nil {
			return fmt.Errorf("Error getting map data from pokeapi: %w", err)
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

func makeCommandMapb(client *PokeClient) func(*Config) error {
	return func(cfg *Config) error {
		if cfg.Prev == "" || cfg.Prev == "null" {
			fmt.Println("You are on the first page.")
			return nil
		}

		bytes, err := client.Get(cfg.Prev)
		if err != nil {
			return fmt.Errorf("Error getting mapb data from pokeapi: %w", err)
		}

		map_response := NamedAPIResourceList{}
		err = json.Unmarshal(bytes, &map_response)
		if err != nil {
			return fmt.Errorf("Error getting mapb data from pokeapi: %w", err)
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

func makeCommandExit() func(*Config) error {
	return func(cfg *Config) error {
		fmt.Print("Closing the Pokedex... Goodbye!\n")
		os.Exit(0)
		return nil
	}
}

func makeCommandHelp(registry map[string]CLICommand) func(*Config) error {
	return func(cfg *Config) error {
		fmt.Print("Welcome to the Pokedex!\n")
		fmt.Print("Usage:\n\n")
		for key := range registry {
			fmt.Printf("%v: %v\n", registry[key].Name, registry[key].Description)
		}
		return nil
	}
}
