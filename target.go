// 24 september 2014

package main

import (
	"flag"
	"runtime"
	"os"
	"path/filepath"
)

// TODO list
var targetOS = flag.String("os", runtime.GOOS, "select target OS")
var targetArch = flag.String("arch", runtime.GOARCH, "select target architecture")

func targetName() string {
	pwd, err := os.Getwd()
	if err != nil {
		// TODO
		panic(err)
	}
	return filepath.Base(pwd)
}
