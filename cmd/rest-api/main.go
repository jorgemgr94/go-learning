package main

import (
	"fmt"
	"go-learning/internal/config"
)

func main() {
	config := config.LoadConfig()

	fmt.Println("Config:", config)

}
