package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Pokemon struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Abilities []Abilities `json:"abilities"`
}

type Abilities struct {
	Ability  Ability `json:"ability"`
	IsHidden bool    `json:"is_hidden"`
	Slot     int     `json:"slot"`
}

type Ability struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func httpExample() {
	fmt.Println("// == HTTP Example =========================================")
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/ditto")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	// Option 1: Unmarshal the JSON into a map
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Accessing nested fields from an unknown JSON structure
	abilities := data["abilities"].([]interface{})
	firstAbility := abilities[0].(map[string]interface{})["ability"].(map[string]interface{})
	fmt.Printf("First ability: %v\n", firstAbility["name"])
	// Access fields like dictionary keys
	if id, exists := data["id"]; exists {
		fmt.Printf("ID: %v\n", id)
	}

	// Option 2: Unmarshal the JSON into a struct
	var pokemon Pokemon
	err = json.Unmarshal(bodyBytes, &pokemon)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	fmt.Printf("Pokemon: %v\n", pokemon.Abilities[0].Ability.Name)
}
