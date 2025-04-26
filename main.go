package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"github.com/jschasse/pokedex/internal/api"
	"math/rand"
	"time"
)


type cliCommand struct {
	name string
	description string
	callback func(c *config, args []string) error
}

type config struct {
	Next string
	Previous string
}

var commands map[string]cliCommand
var pokedex map[string]api.PokemonInfo
const apiURL = "https://pokeapi.co/api/v2/location-area/"
const apiURLPokemon = "https://pokeapi.co/api/v2/pokemon/"

func init(){
	rand.New(rand.NewSource(time.Now().UnixNano()))
	pokedex = make(map[string]api.PokemonInfo)
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
		"catch": {
			name:        "catch",
			description: "Attempts to catch the pokemon that you input",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Displays the stats of the pokemon that you have caught",
			callback:    commandInspect,
		},
	}
}



func main() {
	scanner := bufio.NewScanner(os.Stdin)

	c := &config {
		Next: apiURL,
		Previous: "",
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		
		input := scanner.Text()
		commandWord, args := cleanInput(input)
        if len(commandWord) == 0 && len(args) == 0 {
            continue
        }

        

        command, exists := commands[commandWord]
        if exists {
            err := command.callback(c, args)
            if err != nil {
                fmt.Println("Error:", err)
            }
        } else {
            fmt.Println("Unknown command")
        }
	}

}

func cleanInput(text string) (string, []string) {
	parts := strings.Fields(text)
	trimmedParts := make([]string, len(parts))

	for i, part := range parts {
		trimmedParts[i] = strings.ToLower(strings.TrimSpace(part))
	}

	command := trimmedParts[0]
	

	return command, trimmedParts[1:]
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
	pokeList, err := api.GetPokeapiList(c.Next)
	if err != nil {
		fmt.Print(err)
	}

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
	pokeList, err := api.GetPokeapiList(c.Previous)
	if err != nil {
		fmt.Print(err)
	}

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

func commandExplore(c *config, args []string) error {
	pokeInfo, err := api.GetPokeAreaInfo(apiURL + args[0])
	if err != nil {
		return err
	}

	fmt.Println("Exploring ", args[0])
	fmt.Println("Found Pokemon:")
	for i := 0; i < len(pokeInfo.Pokemon_Encounters); i++ {
		fmt.Printf("- %s\n", pokeInfo.Pokemon_Encounters[i].Pokemon.Name)
	}

	return nil
}

func commandCatch(c *config, args []string) error {
	if len(args) == 0 {
        return fmt.Errorf("you must provide a pokemon name to catch")
    }
    if len(args) > 1 {
        return fmt.Errorf("please provide only one pokemon name")
    }

	pokemonName := args[0]

	poke, err := api.GetPokemonInfo(apiURLPokemon + pokemonName)
	if err != nil {
		return err
	}

    const maxThreshold = 400
    threshold := maxThreshold - poke.Base_Experience
    if threshold < 20 {
        threshold = 20
    }


    randomNumber := rand.Intn(maxThreshold)

    fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
    time.Sleep(1 * time.Second)


    if randomNumber < threshold {
        fmt.Printf("%s was caught!\n", pokemonName)
        
        pokedex[pokemonName] = poke
    } else {
        fmt.Printf("%s escaped!\n", pokemonName)
    }

    return nil
}

func commandInspect(c *config, args []string) error {
	if len(args) == 0 {
        return fmt.Errorf("you must provide a pokemon name")
    }
    if len(args) > 1 {
        return fmt.Errorf("please provide only one pokemon name")
    }

	pokemonName := args[0]

	pokemon, exists := pokedex[pokemonName]
	if exists == false {
		return fmt.Errorf("%s is not in your pokedex", pokemonName)
	}

	fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\n", pokemon.Name, pokemon.Height, pokemon.Weight)
	fmt.Println("Stats:")
	for i := 0; i < len(pokemon.Stats); i++ {
		fmt.Printf("	-%s: %d\n", pokemon.Stats[i].Stat.Name, pokemon.Stats[i].Base_Stat)
	}
	fmt.Println("Types:")
	for i := 0; i < len(pokemon.Types); i++ {
		fmt.Printf("	- %s\n", pokemon.Types[i].Type.Name)
	}

	return nil

}
