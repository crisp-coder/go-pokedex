package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {

	cfg := Config{
		Next: "https://pokeapi.co/api/v2/location-area/",
		Prev: "",
	}

	command_registry := make(map[string]cliCommand)
	commandHelp := makeCommandHelp(&cfg, command_registry)
	commandExit := makeCommandExit(&cfg)
	commandMap := makeCommandMap(&cfg)
	commandMapb := makeCommandMapb(&cfg)

	command_registry["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
	command_registry["help"] = cliCommand{
		name:        "help",
		description: "Prints a list of commands.",
		callback:    commandHelp,
	}
	command_registry["map"] = cliCommand{
		name:        "map",
		description: "Returns a list of 20 map areas. Call again to get the next 20",
		callback:    commandMap,
	}
	command_registry["mapb"] = cliCommand{
		name:        "mapb",
		description: "Returns a list of previous 20 map areas.",
		callback:    commandMapb,
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			fmt.Println("Scanner error reading input")
			return
		}

		line := scanner.Text()
		cleanedInput := cleanInput(line)

		if len(cleanedInput) > 0 {
			if val, ok := command_registry[cleanedInput[0]]; !ok {
				fmt.Print("Unknown command\n")
			} else {
				err := val.callback(&cfg)
				if err != nil {
					fmt.Printf("Error: %v", err)
				}
			}
		}
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	Next string
	Prev string
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

		map_response := NamedAPIResourceList{}
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

		map_response := NamedAPIResourceList{}
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

func makeCommandHelp(cfg *Config, registry map[string]cliCommand) func(*Config) error {
	return func(cfg *Config) error {
		fmt.Print("Welcome to the Pokedex!\n")
		fmt.Print("Usage:\n\n")
		for key := range registry {
			fmt.Printf("%v: %v\n", registry[key].name, registry[key].description)
		}
		return nil
	}
}

func cleanInput(text string) []string {
	text_lower := strings.ToLower(text)
	var fields []string
	fields = strings.Fields(text_lower)
	return fields
}
