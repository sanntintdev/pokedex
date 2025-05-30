package cli

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"

	"github.com/sanntintdev/pokedex/internal/api"
	"github.com/sanntintdev/pokedex/internal/models"
)

var caughtPokedex = make(map[string]models.Pokemon)

type CliCommand struct {
	Name        string
	Description string
	Callback    func(c *Config, args []string) error
}

func GetCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"help": {
			Name:        "help",
			Description: "Display help information",
			Callback:    commandHelp,
		},
		"map": {
			Name:        "map",
			Description: "Display the locations of the Pokemon World",
			Callback:    commandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Display the previous locations of the Pokemon World",
			Callback:    commandMapPrev,
		},
		"explore": {
			Name:        "explore",
			Description: "Explore a location area",
			Callback:    commandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "Catch a Pokemon",
			Callback:    commandCatch,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspect a Pokemon",
			Callback:    commandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "Display the caught Pokemon",
			Callback:    commandPokedex,
		},
	}
}

func commandExit(c *Config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *Config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for Name, command := range GetCommands() {
		fmt.Printf("%s: %s\n", Name, command.Description)
	}
	return nil
}

func commandMap(c *Config, args []string) error {
	url := "https://pokeapi.co/api/v2/location-area"

	// after first api call , if there is next api , this makes the next api call when we call commandMap again
	if c.Next != "" {
		url = c.Next
	}

	locations, err := api.FetchLocations(url)

	if err != nil {
		return err
	}

	if len(locations.Results) == 0 {
		fmt.Println("No locations found.")
		return nil
	}

	c.Next = locations.Next
	c.Prev = locations.Previous

	for _, location := range locations.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}

func commandMapPrev(c *Config, args []string) error {
	if c.Prev == "" {
		fmt.Println("You're on the first page. No previous locations to show.")
		return nil
	}

	locations, err := api.FetchLocations(c.Prev)

	if err != nil {
		return err
	}

	if len(locations.Results) == 0 {
		fmt.Println("No locations found.")
		return nil
	}

	c.Next = locations.Next
	c.Prev = locations.Previous

	for _, location := range locations.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}

func commandExplore(c *Config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: explore <location-area>")
	}

	locationArea := args[0]
	url := "https://pokeapi.co/api/v2/location-area/" + locationArea

	fmt.Printf("Exploring %s...\n", locationArea)
	locations, err := api.FetchLocationArea(url)

	if err != nil {
		return err
	}

	if len(locations.PokemonEncounters) == 0 {
		fmt.Println("No locations found.")
		return nil
	}

	fmt.Print("Found Pokemon: \n")
	for _, pokemon := range locations.PokemonEncounters {
		fmt.Println("- ", pokemon.Pokemon.Name)
	}

	return nil
}

func commandCatch(c *Config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: catch <pokemon-name>")
	}

	pokemonName := strings.ToLower(args[0])
	if _, alreadyCaught := caughtPokedex[pokemonName]; alreadyCaught {
		fmt.Printf("%s is already caught!\n", pokemonName)
		return nil
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	pokemon, err := api.FetchPokemon(url)
	if err != nil {
		return err
	}

	baseExp := pokemon.BaseExperience
	catchChance := calculateCatchChance(baseExp)

	roll := 100.0 * rand.Float64()

	if roll >= catchChance {
		caughtPokedex[pokemonName] = *pokemon
		fmt.Printf("%v was caught! \n", pokemon.Name)
	} else {
		fmt.Printf("%v was escaped! \n", pokemon.Name)
	}
	return nil
}

func calculateCatchChance(baseExp int) float64 {
	catchChance := 100.0 - (float64(baseExp) * 0.3)

	if catchChance < 5.0 {
		catchChance = 5.0
	}
	if catchChance > 95.0 {
		catchChance = 95.0
	}

	return catchChance
}

func commandInspect(c *Config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: inspect <pokemon-name>")
	}

	pokemonName := strings.ToLower(args[0])
	if _, alreadyCaught := caughtPokedex[pokemonName]; !alreadyCaught {
		fmt.Printf("%s is not caught!\n", pokemonName)
		return nil
	}

	pokemon := caughtPokedex[pokemonName]

	fmt.Println("Name: ", pokemon.Name)
	fmt.Println("Height:", pokemon.Height)
	fmt.Println("Weight:", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typ := range pokemon.Types {
		fmt.Printf("  -%s\n", typ.Type.Name)
	}
	return nil
}

func commandPokedex(c *Config, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("usage: pokedex")
	}

	fmt.Println("Caught Pokemon:")
	for _, pokemon := range caughtPokedex {
		fmt.Printf("  -%s\n", pokemon.Name)
	}
	return nil
}
