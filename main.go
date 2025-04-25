package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"github.com/jschasse/pokedex/internal/api"
)


type cliCommand struct {
	name string
	description string
	callback func(c *config) error
}

type config struct {
	Next string
	Previous string
}

var commands map[string]cliCommand


func init(){
	commands = map[string]cliCommand {
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the name of the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the name of the previous 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays the name of the pokemon avaliable at location",
			callback:    commandExplore,
		},
	}
}



func main() {
	scanner := bufio.NewScanner(os.Stdin)
	url := "https://pokeapi.co/api/v2/location-area/"

	c := &config {
		Next: url,
		Previous: "",
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		
		input := scanner.Text()
		cleanedInput := cleanInput(input)
        if len(cleanedInput) == 0 {
            continue
        }

        commandName := cleanedInput[0]

        command, exists := commands[commandName]
        if exists {
            err := command.callback(c)
            if err != nil {
                fmt.Println("Error:", err)
            }
        } else {
            fmt.Println("Unknown command")
        }
	}

}

func cleanInput(text string) []string {
	parts := strings.Fields(text)
	trimmedParts := make([]string, len(parts))

	for i, part := range parts {
		trimmedParts[i] = strings.ToLower(strings.TrimSpace(part))
	}
	

	return trimmedParts
}

func commandExit(c *config, args []string) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	for key, _ := range commands {
		fmt.Printf("%s: %s\n", commands[key].name, commands[key].description)
	}
	fmt.Println()
	return nil
}

func commandMap(c *config, args []string) error {
	pokeList, _ := api.GetPokeapiList(c.Next)

	for i := 0; i < len(pokeList.Results); i++ {
		fmt.Printf("%s\n", pokeList.Results[i].Name)
	}


	if pokeList.Previous != nil {
        c.Previous = *pokeList.Previous
    }
    if pokeList.Next != nil {
        c.Next = *pokeList.Next
    }

	return nil
}

func commandMapb(c *config, args []string) error {
	pokeList, _ := api.GetPokeapiList(c.Previous)

	for i := 0; i < len(pokeList.Results); i++ {
		fmt.Printf("%s\n", pokeList.Results[i].Name)
	}

	if pokeList.Previous != nil {
        c.Previous = *pokeList.Previous
    }
    if pokeList.Next != nil {
        c.Next = *pokeList.Next
    }

	return nil
}

func commandExplore(args []string) error {

}

