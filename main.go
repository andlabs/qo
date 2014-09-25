// 24 september 2014

package main

import (
	"os"
	"flag"
	"path/filepath"
)

func main() {
	flag.Parse()
	if *selectedToolchain == "list" {
		listToolchains()
		os.Exit(0)
	}
	err := filepath.Walk(".", walker)
	if err != nil {
		panic(err)
	}
	compileFlags()
	buildScript()
	run()
	os.Exit(0)		// success
}
