// 24 september 2014

package main

import (
	"fmt"
	"os"
)

var builder = make(chan *Executor)

var percentPer float64
var progress float64

func printProgress(step string) {
	fmt.Printf("[%3d%%] %s\n", int(progress), step)
}

func runStage(s Stage) (success bool) {
	indices := make(map[*Executor]int)
	for i, e := range s {
		indices[e] = i
		go e.Do()
	}
	for len(indices) != 0 {
		e := <-builder
		delete(indices, e)
		if e.Error != nil {
			fmt.Fprintf(os.Stderr, "%s", e.Output.Bytes())
			// ensure only one newline
			if e.Output.Len() == 0 || e.Output.Bytes()[e.Output.Len() - 1] != '\n' {
				fmt.Fprintf(os.Stderr, "\n")
			}
			fmt.Fprintf(os.Stderr, "[FAIL] Step %q failed with error: %v\n", e.Name, e.Error)
			return false
		}
		progress += percentPer
		printProgress(e.Name)
	}
	return true
}

func run() {
	percentPer = 100 / float64(nStages)
	progress = 0
	printProgress("Beginning build")
	for _, stage := range script {
		if !runStage(stage) {
			// TODO alert failure
			os.Exit(1)
		}
	}
}
