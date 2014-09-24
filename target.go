// 24 september 2014

package main

import (
	"flag"
	"runtime"
)

// TODO list
var targetOS = flag.String("os", runtime.GOOS, "target OS")
var targetArch = flag.String("arch", runtime.GOARCH, "target architecture")
