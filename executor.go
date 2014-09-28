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

type Executor struct {
	Name	string
	Line		[]string
	Output	*bytes.Buffer
	Error		error
}

func (e *Executor) Do() {
	if *showall {
		fmt.Printf("%s\n", strings.Join(e.Line, " "))
	}
	cmd := exec.Command(e.Line[0], e.Line[1:]...)
	cmd.Env = os.Environ()
	e.Output = new(bytes.Buffer)
	cmd.Stdout = e.Output
	cmd.Stderr = e.Output
	e.Error = cmd.Run()
	builder <- e
}
