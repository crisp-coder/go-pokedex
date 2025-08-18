package pokeapi

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	API           string
	NextArea      string
	PrevArea      string
	ExploreTarget string
	CaptureTarget string
}

type CLICommand struct {
	Name        string
	Description string
	Callback    func() error
}

func MakeCommandRegistry(cfg *Config, client *PokeClient) map[string]CLICommand {
	command_registry := make(map[string]CLICommand)
	commandHelp := makeCommandHelp(cfg, command_registry)
	commandExit := makeCommandExit(cfg)
	commandMap := makeCommandMap(cfg, client)
	commandMapb := makeCommandMapb(cfg, client)
	commandExplore := makeCommandExplore(cfg, client)

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

func makeCommandExplore(cfg *Config, client *PokeClient) func() error {
	return func() error {
		bytes, err := client.Get(cfg.API + "location-area/" + cfg.ExploreTarget)
		if err != nil {
			return fmt.Errorf("error getting explore area data: %w", err)
		}

		exploreResponse := LocationArea{}
		err = json.Unmarshal(bytes, &exploreResponse)
		if err != nil {
			return fmt.Errorf("error getting explore area data: %w", err)
		}

		if len(exploreResponse.Pokemon_encounters) > 0 {
			for _, encounter := range exploreResponse.Pokemon_encounters {
				fmt.Printf(" - %v\n", encounter.Pokemon.Name)
			}

		} else {
			fmt.Println("No pokemon found in this area.")
		}

		return nil
	}
}

func makeCommandMap(cfg *Config, client *PokeClient) func() error {
	return func() error {
		if cfg.NextArea == "" || cfg.NextArea == "null" {
			cfg.NextArea = "?offset=0&limit=20"
		}

		bytes, err := client.Get(cfg.API + "location-area/" + cfg.NextArea)
		if err != nil {
			return fmt.Errorf("error getting map data from pokeapi: %w", err)
		}

		map_response := NamedAPIResourceList{}
		err = json.Unmarshal(bytes, &map_response)
		if err != nil {
			return fmt.Errorf("error getting map data from pokeapi: %w", err)
		}

		if map_response.Next != "" && map_response.Next != "null" {
			cfg.NextArea = map_response.Next[40:]
		} else {
			cfg.NextArea = ""
		}

		if map_response.Previous != "" && map_response.Previous != "null" {
			cfg.PrevArea = map_response.Previous[40:]
		} else {
			cfg.PrevArea = ""
		}

		for _, value := range map_response.Results {
			fmt.Printf("%v\n", value.Name)
		}

		return nil
	}
}

func makeCommandMapb(cfg *Config, client *PokeClient) func() error {
	return func() error {
		if cfg.PrevArea == "" || cfg.PrevArea == "null" {
			cfg.PrevArea = "?offset=1063&limit=20"
		}

		bytes, err := client.Get(cfg.API + "location-area/" + cfg.PrevArea)
		if err != nil {
			return fmt.Errorf("error getting mapb data from pokeapi: %w", err)
		}

		map_response := NamedAPIResourceList{}
		err = json.Unmarshal(bytes, &map_response)
		if err != nil {
			return fmt.Errorf("error getting mapb data from pokeapi: %w", err)
		}

		if map_response.Next != "" && map_response.Next != "null" {
			cfg.NextArea = map_response.Next[40:]
		} else {
			cfg.NextArea = ""
		}

		if map_response.Previous != "" && map_response.Previous != "null" {
			cfg.PrevArea = map_response.Previous[40:]
		} else {
			cfg.PrevArea = ""
		}

		for _, value := range map_response.Results {
			fmt.Printf("%v\n", value.Name)
		}

		return nil
	}
}

func makeCommandExit(cfg *Config) func() error {
	return func() error {
		fmt.Print("Closing the Pokedex... Goodbye!\n")
		os.Exit(0)
		return nil
	}
}

func makeCommandHelp(cfg *Config, registry map[string]CLICommand) func() error {
	return func() error {
		fmt.Print("Welcome to the Pokedex!\n")
		fmt.Print("Usage:\n\n")
		for key := range registry {
			fmt.Printf("%v: %v\n", registry[key].Name, registry[key].Description)
		}
		return nil
	}
}
