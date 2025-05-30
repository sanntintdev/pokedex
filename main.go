package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sanntintdev/pokedex/internal/cli"
)

func cleanInput(text string) []string {
	return strings.Fields(text)
}

func main() {
	cfg := &cli.Config{}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input := scanner.Text()
			cleanedInput := cleanInput(input)

			if len(cleanedInput) == 0 {
				fmt.Println("Please enter a command !")
				continue
			}
			commandName := cleanedInput[0]
			args := cleanedInput[1:]
			command, ok := cli.GetCommands()[commandName]
			if ok {
				err := command.Callback(cfg, args)
				if err != nil {
					fmt.Print(err)
					os.Exit(0)
				}
			} else {
				fmt.Println("Unknown command.")
			}
		}
	}
}
