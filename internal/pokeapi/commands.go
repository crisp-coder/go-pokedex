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
	Pokedex       TempPokedex
	CommandParams []string
}

type TempPokedex struct {
	KnownPokemon map[string]Pokemon
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
	commandCatch := makeCommandCatch(cfg, client)
	commandInspect := makeCommandInspect(cfg)
	commandPokedex := makeCommandPokedex(cfg)

	command_registry["exit"] = CLICommand{
		Name:        "exit",
		Description: "Exit the Pokedex",
		Callback:    commandExit,
	}
	command_registry["help"] = CLICommand{
		Name:        "help",
		Description: "Lists commands.",
		Callback:    commandHelp,
	}
	command_registry["map"] = CLICommand{
		Name:        "map",
		Description: "Lists next 20 map areas.",
		Callback:    commandMap,
	}
	command_registry["mapb"] = CLICommand{
		Name:        "mapb",
		Description: "Lists previous 20 map areas.",
		Callback:    commandMapb,
	}
	command_registry["explore"] = CLICommand{
		Name:        "explore",
		Description: "Explore an area. Lists pokemon in area. Accepts a Name.",
		Callback:    commandExplore,
	}
	command_registry["catch"] = CLICommand{
		Name:        "catch",
		Description: "Attempt to catch a pokemon. Accepts a Name.",
		Callback:    commandCatch,
	}
	command_registry["inspect"] = CLICommand{
		Name:        "inspect",
		Description: "List height, weight, stats, types of pokemon. Accepts a Name.",
		Callback:    commandInspect,
	}
	command_registry["pokedex"] = CLICommand{
		Name:        "pokedex",
		Description: "Lists all pokemon in user's pokedex.",
		Callback:    commandPokedex,
	}

	return command_registry
}

func makeCommandPokedex(cfg *Config) func() error {
	return func() error {
		fmt.Printf("Your Pokedex:\n")
		pokedex := cfg.Pokedex
		for _, pokemon := range pokedex.KnownPokemon {
			fmt.Printf(" - %v\n", pokemon.Name)
		}
		return nil
	}
}

func makeCommandInspect(cfg *Config) func() error {
	return func() error {
		if pokemon, ok := cfg.Pokedex.KnownPokemon[cfg.CommandParams[0]]; ok {
			fmt.Printf("Name: %s\n", pokemon.Name)
			fmt.Printf("Height: %v\n", pokemon.Height)
			fmt.Printf("Weight: %v\n", pokemon.Weight)
			fmt.Printf("Stats:\n")
			for _, stat := range pokemon.Stats {
				fmt.Printf(" - %v: %v\n", stat.Stat.Name, stat.Base_stat)
			}
			fmt.Printf("Types:\n")
			for _, poke_type := range pokemon.Types {
				fmt.Printf(" - %v\n", poke_type.Type.Name)
			}

		} else {
			fmt.Println("You have not caught that pokemon.")
		}

		return nil
	}
}

func makeCommandCatch(cfg *Config, client *PokeClient) func() error {
	return func() error {
		bytes, err := client.Get(cfg.API + "pokemon/" + cfg.CommandParams[0])
		if err != nil {
			return fmt.Errorf("error getting pokemon data: %w", err)
		}

		catch_response := Pokemon{}
		err = json.Unmarshal(bytes, &catch_response)
		if err != nil {
			return fmt.Errorf("error getting pokemon data: %w", err)
		}

		fmt.Printf("Throwing a Pokeball at %s...\n", catch_response.Name)
		fmt.Printf("Caught %s!\n", catch_response.Name)
		cfg.Pokedex.KnownPokemon[catch_response.Name] = catch_response

		return nil
	}
}

func makeCommandExplore(cfg *Config, client *PokeClient) func() error {
	return func() error {
		bytes, err := client.Get(cfg.API + "location-area/" + cfg.CommandParams[0])
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
