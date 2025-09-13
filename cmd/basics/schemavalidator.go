package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

func schemaValidator() {
	// Compile the schema
	compiler := jsonschema.NewCompiler()

	schema, err := compiler.Compile("cmd/basics/data/schema.json")
	if err != nil {
		fmt.Println("Schema error:", err)
		os.Exit(1)
	}

	// Open the JSON file to validate
	dataFile, err := os.Open("cmd/basics/data/payload.json")
	if err != nil {
		fmt.Println("Failed to open data:", err)
		os.Exit(1)
	}
	defer dataFile.Close()

	// Decode JSON
	var data interface{}
	if err := json.NewDecoder(dataFile).Decode(&data); err != nil {
		fmt.Println("Failed to decode JSON:", err)
		os.Exit(1)
	}

	// Validate
	if err := schema.Validate(data); err != nil {
		fmt.Println("Validation error:", err)
	} else {
		fmt.Println("Validation successful!")
	}
}
