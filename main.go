package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	cfg := Config{
		Next: "https://pokeapi.co/api/v2/location-area/",
		Prev: "",
	}

	command_registry := makeCommandRegistry(&cfg)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			fmt.Println("Scanner error reading input.")
			return
		}

		line := scanner.Text()
		cleanedInput := cleanInput(line)
		if len(cleanedInput) == 0 {
			fmt.Println("Input was empty.")
			continue
		}
		command := cleanedInput[0]
		val, ok := command_registry[command]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := val.callback(&cfg)
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
	}
}
