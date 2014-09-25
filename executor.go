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

/*
func main() {
	go (&Executor{
		Name:	"echo",
		Line:		[]string{"echo", "hello,", "world"},
	}).Do()
	go (&Executor{
		Name:	"sleep",
		Line:		[]string{"sleep", "5"},
	}).Do()
	go (&Executor{
		Name:	"badcommand",
		Line:		[]string{"badcommand"},
	}).Do()
	go (&Executor{
		Name:	"stderr",
		Line:		[]string{"gcc", "--qwertyuiop"},
	}).Do()
	for i := 0; i < 4; i++ {
		e := <-builder
		fmt.Printf("done %q %v\n", e.Name, e.Error)
		fmt.Printf("%q %q\n", e.Stdout.String(), e.Stderr.String())
	}
}
*/
