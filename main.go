package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	command_registry := make(map[string]cliCommand)
	commandHelp := makeCommandHelp(command_registry)
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
				val.callback()
			}
		}
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandMap() error {
	return nil
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func makeCommandHelp(registry map[string]cliCommand) func() error {
	return func() error {
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
