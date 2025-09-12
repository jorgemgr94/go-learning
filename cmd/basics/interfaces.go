package main

import "fmt"

// This interface says that any other type that has a function getGreeting
// with no arguments and returns a string is also of type bot
type bot interface {
	getGreeting() string
}

type englishBot struct{}
type spanishBot struct{}

func interfaces() {
	fmt.Println("// == Interfaces ===========================================")
	eb := englishBot{}
	sb := spanishBot{}

	printGreeting(eb)
	printGreeting(sb)
}

func (englishBot) getGreeting() string {
	return "Hi there!"
}

func (spanishBot) getGreeting() string {
	return "Hola!"
}

func printGreeting(b bot) {
	println(b.getGreeting())
}
