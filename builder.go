// 24 september 2014

package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
)

var percentPer float64
var progress float64

func printProgress(step string) {
	fmt.Printf("[%3d%%] %s\n", int(progress), step)
}

func runner(in <-chan *Step, out chan<- *Step, wg *sync.WaitGroup) {
	for s := range in {
		s.Do()
		out <- s
	}
	wg.Done()
}

func runStage(s Stage) {
	wg := new(sync.WaitGroup)
	in := make(chan *Step)
	out := make(chan *Step)
	// TODO make an option
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker(in, out, wg)
	}
	wg.Add(1)
	go func() {
		for i := 0; i < len(s); i++ {
			in <- s[i]
		}
		close(in)
		wg.Done()
	}()
	for i := 0; i < len(s); i++ {
		s := <-out
		fmt.Fprintf(os.Stderr, "%s", s.Output.Bytes())
		// ensure only one newline
		if s.Output.Len() != 0 && s.Output.Bytes()[s.Output.Len() - 1] != '\n' {
			fmt.Fprintf(os.Stderr, "\n")
		}
		if s.Error != nil {
			fail("Step %q failed with error: %v", s.Name, s.Error)
		}
		progress += percentPer
		printProgress(s.Name)
	}
	wg.Wait()
	close(out)
}

func run() {
	percentPer = 100 / float64(nSteps)
	progress = 0
	printProgress("Beginning build")
	for _, stage := range script {
		runStage(stage)
	}
}
