package main

import (
	"fmt"
	"log"
	"os"
	"quake-logs-parser/analytics"
)

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal("You must inform a filepath as argument the program")
		return
	}

	path := args[0]
	report := analytics.ProcessLog(path)

	for i, v := range report.Games {
		fmt.Printf("{\"game_%d\": %s}\n", i+1, v)
	}
}
