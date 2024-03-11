package main

import (
	"fmt"
	"log"
	"os"
	analytics "quake-logs-parser/analytics"
)

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal("You must inform a filepath as argument the program")
		return
	}

	filepath := args[0]

	report := analytics.Report{Games: []analytics.Game{}}
	analytics.ProcessLog(filepath, &report)

	for i, v := range report.Games {
		f := fmt.Sprintf(`{"game_%d": %s}`, i+1, v)
		fmt.Println(f)
		fmt.Printf("Ranking:\n%s", v.Ranking())
	}
}
