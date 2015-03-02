// 24 september 2014

package main

import (
	"fmt"
	"os"
	"os/exec"
	"flag"
	"bytes"
	"strings"
)

var showall = flag.Bool("x", false, "show all commands as they run")

// Step represents a single step in the build process.
type Step struct {
	Name	string
	Line		[]string
	Output	*bytes.Buffer
	Error		error
}

func (s *Step) Do() {
	if *showall {
		fmt.Printf("%s\n", strings.Join(s.Line, " "))
	}
	cmd := exec.Command(s.Line[0], s.Line[1:]...)
	cmd.Env = os.Environ()
	s.Output = new(bytes.Buffer)
	cmd.Stdout = s.Output
	cmd.Stderr = s.Output
	s.Error = cmd.Run()
}

// Stage is a list of Steps.
// Each Step in the Stage is run concurrently.
// A build process consists of multiple Stages.
// All Steps in the current Stage must be run to completion before other Steps will run.
type Stage []*Step
