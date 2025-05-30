package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sanntintdev/pokedex/internal/models"
)

// Global variable for location cache
var pokeCache = NewCache(5 * time.Minute)

func FetchLocations(url string) (*models.Location, error) {
	var locations models.Location
	cacheData, cached := pokeCache.Get(url)

	if cached != false {
		err := json.Unmarshal(cacheData, &locations)
		if err != nil {
			return nil, err
		}

		return &locations, nil
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&locations)
	if err != nil {
		return nil, err
	}

	dataToCache, err := json.Marshal(res.Body)
	if err != nil {
		return nil, err
	}
	pokeCache.Add(url, dataToCache)

	return &locations, nil
}

func FetchLocationArea(url string) (*models.LocationArea, error) {
	var locationArea models.LocationArea

	cacheData, cached := pokeCache.Get(url)

	if cached != false {
		err := json.Unmarshal(cacheData, &locationArea)
		if err != nil {
			return nil, err
		}

		return &locationArea, nil
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&locationArea)
	if err != nil {
		return nil, err
	}

	dataToCache, err := json.Marshal(res.Body)
	if err != nil {
		return nil, err
	}
	pokeCache.Add(url, dataToCache)

	return &locationArea, nil
}

func FetchPokemon(url string) (*models.Pokemon, error) {
	var pokemon models.Pokemon

	cacheData, cached := pokeCache.Get(url)

	if cached == true {
		err := json.Unmarshal(cacheData, &pokemon)
		if err != nil {
			return nil, err
		}

		return &pokemon, nil
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&pokemon)
	if err != nil {
		return nil, err
	}

	dataToCache, err := json.Marshal(pokemon)
	if err != nil {
		return nil, err
	}

	pokeCache.Add(url, dataToCache)

	return &pokemon, nil
}
