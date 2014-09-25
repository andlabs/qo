// 24 september 2014

package main

import (
	"os"
)

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
	}
	return true
}

func run() {
	for _, stage := range script {
		if !runStage(stage) {
			// TODO alert failure
			os.Exit(1)
		}
	}
}
