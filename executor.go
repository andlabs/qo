// 24 september 2014

package main

import (
	"os"
	"os/exec"
	"bytes"
)

type Executor struct {
	Name	string
	Line		[]string
	Stdout	*bytes.Buffer
	Stderr	*bytes.Buffer
	Error		error
}

func (e *Executor) Do() {
	cmd := exec.Command(e.Line[0], e.Line[1:]...)
	cmd.Env = os.Environ()
	e.Stdout = new(bytes.Buffer)
	cmd.Stdout = e.Stdout
	e.Stderr = new(bytes.Buffer)
	cmd.Stderr = e.Stderr
	e.Error = cmd.Run()
	builder <- e
}

var builder = make(chan *Executor)

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
