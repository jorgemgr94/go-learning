package main

import "sync"

func main() {
	httpExample()
	variables()
	maps()
	slices()
	interfaces()

	// older way to use wait groups
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		channels()
	}()

	structs()
	schemaValidator()
	closures()
	mutex()

	// wait for the channels goroutine to complete
	wg.Wait()
}
