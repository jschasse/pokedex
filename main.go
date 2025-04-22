package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		fmt.Printf("Your command was: %s\n", cleanedInput[0])
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

