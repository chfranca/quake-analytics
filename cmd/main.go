package main

import (
	"fmt"
	analytics "quake-logs-parser/analytics"
)

func main() {

	file := "qgames.log"
	report := analytics.Report{Games: []analytics.Game{}}
	analytics.ProcessLog(file, &report)

	for i, v := range report.Games {
		f := fmt.Sprintf(`{"game_%d": %s}`, i+1, v)
		fmt.Println(f)
		fmt.Printf("Ranking:\n%s", v.Ranking())
	}
}
