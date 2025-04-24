package api

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io"
)

type NamedAPIResource struct {
    Name string 
    URL  string 
}

type PokeapiList struct {
    Count    int                
    Next     *string            
    Previous *string            
    Results  []NamedAPIResource 
}

func GetPokeapiList(url string) (PokeapiList, error) {
	res, err := http.Get(url)
	if err != nil {
		return PokeapiList{}, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return PokeapiList{}, fmt.Errorf("error reading response body from %s: %w", url, err)
	}

	var pokeList PokeapiList

	err = json.Unmarshal(data, &pokeList)
	if err != nil {
		return PokeapiList{}, fmt.Errorf("error unmarshalling location list from pokeapi %w", err)
	}

	return pokeList, nil
}
