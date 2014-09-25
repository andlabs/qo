// 24 september 2014

package main

import (
	"fmt"
	"os"
)

var percentPer float64
var progress float64

func runStage(s Stage) (success bool) {
	indices := make(map[*Executor]int)

	for i, e := range s {
		indices[e] = i
		go e.Do()
	}
	for len(indices) != 0 {
		e := <-builder
		delete(indices, e)
		// TODO check error
		progress += percentPer
		fmt.Printf("[%3d%%] %s\n", int(progress), e.Name)
	}
	return true
}

func run() {
	percentPer = 100 / float64(nStages)
	progress = 0
	for _, stage := range script {
		if !runStage(stage) {
			// TODO alert failure
			os.Exit(1)
		}
	}
}
