package api

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io"
	"github.com/jschasse/pokedex/internal/pokecache"
	"time"
)

type NamedAPIResource struct {
    Name string 
    URL  string 
}

type PokemonStat struct {
	Stat NamedAPIResource
	Base_Stat int
}

type PokemonTypes struct {
	Type NamedAPIResource
}

type PokemonEncounter struct {
	Pokemon NamedAPIResource
}

type PokeapiList struct {
    Count    int                
    Next     *string            
    Previous *string            
    Results  []NamedAPIResource 
}

type PokeAreaInfo struct {
	Location NamedAPIResource
	Pokemon_Encounters []PokemonEncounter
}

type PokemonInfo struct {
	Name string
	Height int
	Weight int
	Base_Experience int
	Stats []PokemonStat
	Types []PokemonTypes
}

var cache = pokecache.NewCache(5 * time.Second)

func GetPokeapiList(url string) (PokeapiList, error) {
	
	cacheData, exists := cache.Get(url)
	if exists {
		var pokeList PokeapiList
		err := json.Unmarshal(cacheData, &pokeList)
		if err == nil {
			return pokeList, nil
		} 
	}
	res, err := http.Get(url)
	if err != nil {
		return PokeapiList{}, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return PokeapiList{}, fmt.Errorf("error reading response body from %s: %w", url, err)
	}

	cache.Add(url, data)

	var pokeList PokeapiList

	err = json.Unmarshal(data, &pokeList)
	if err != nil {
		return PokeapiList{}, fmt.Errorf("error unmarshalling location list from pokeapi %w", err)
	}

	return pokeList, nil
	
}

func GetPokeAreaInfo(url string) (PokeAreaInfo, error) {
	cacheData, exists := cache.Get(url)

	if exists {
		var pokeInfo PokeAreaInfo
		err := json.Unmarshal(cacheData, &pokeInfo)
		if err == nil {
			return pokeInfo, nil
		} 
	}

	res, err := http.Get(url)
	if err != nil {
		return PokeAreaInfo{}, fmt.Errorf("error creating request: %w", err)
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return PokeAreaInfo{}, fmt.Errorf("error reading response body from %s: %w", url, err)
	}

	cache.Add(url, data)

	var pokeInfo PokeAreaInfo

	err = json.Unmarshal(data, &pokeInfo)
	if err != nil {
		return PokeAreaInfo{}, fmt.Errorf("error unmarshalling location list from pokeapi %w", err)
	}

	return pokeInfo, nil
}

func GetPokemonInfo(url string) (PokemonInfo, error) {
	res, err := http.Get(url)
	if err != nil {
		return PokemonInfo{}, fmt.Errorf("error creating request: %w", err)
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return PokemonInfo{}, fmt.Errorf("error reading response body from %s: %w", url, err)
	}

	cache.Add(url, data)

	var poke PokemonInfo

	err = json.Unmarshal(data, &poke)
	if err != nil {
		return PokemonInfo{}, fmt.Errorf("error unmarshalling location list from pokeapi %w", err)
	}

	return poke, nil
}
