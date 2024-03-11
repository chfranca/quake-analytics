package main

import (
	"fmt"
	quakeanalytics "quake-logs-parser/quake-analytics"
)

func main() {

	report := quakeanalytics.Report{Games: []quakeanalytics.Game{}}
	quakeanalytics.ProcessLog("qgames.log", &report)

	for i, v := range report.Games {
		f := fmt.Sprintf(`{"game_%d": %s}`, i+1, v)
		fmt.Println(f)
		fmt.Printf("Ranking:\n%s", v.Ranking())
	}
}
