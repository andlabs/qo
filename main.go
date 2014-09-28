// 24 september 2014

package main

import (
	"fmt"
	"os"
	"flag"
	"path/filepath"
	"strings"
)

func main() {
	flag.Parse()
	if *selectedToolchain == "list" {
		listToolchains()
		os.Exit(0)
	}
	if *targetOS == "list" {
		fmt.Printf("%s\n", strings.Join(supportedOSs, "\n"))
		os.Exit(0)
	}
	if *targetArch == "list" {
		fmt.Printf("%s\n", strings.Join(supportedArchs, "\n"))
		os.Exit(0)
	}
	computeExcludeSuffixes()
	err := filepath.Walk(".", walker)
	if err != nil {
		panic(err)
	}
	compileFlags()
	buildScript()
	err = os.MkdirAll(".qoobj", 0755)
	if err != nil {
		// TODO
		panic(err)
	}
	run()
	os.Exit(0)		// success
}
