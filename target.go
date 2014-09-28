// 24 september 2014

package main

import (
	"flag"
	"runtime"
	"os"
	"path/filepath"
	"strings"
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

var supportedOSs = strings.Fields("windows darwin linux freebsd openbsd netbsd dragonfly solaris")
var supportedArchs = strings.Fields("386 amd64")

var excludeSuffixes []string
var excludeFolders []string

func computeExcludeSuffixes() {
	for _, os := range supportedOSs {
		if os == *targetOS {
			continue
		}
		excludeSuffixes = append(excludeSuffixes, "_" + os)
		excludeFolders = append(excludeFolders, os)
		for _, arch := range supportedArchs {
			excludeSuffixes = append(excludeSuffixes, "_" + os + "_" + arch)
			excludeFolders = append(excludeFolders, os + "_" + arch)
		}
	}
	for _, arch := range supportedArchs {
		if arch == *targetArch {
			continue
		}
		excludeSuffixes = append(excludeSuffixes, "_" + arch)
		excludeFolders = append(excludeFolders, arch)
	}
}

func excludeFile(filename string) bool {
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	base = base[:len(base) - len(ext)]
	for _, suffix := range excludeSuffixes {
		if strings.HasSuffix(base, suffix) {
			return true
		}
	}
	return false
}

func excludeDir(filename string) bool {
	base := filepath.Base(filename)
	for _, exclude := range excludeFolders {
		if base == exclude {
			return true
		}
	}
	return false
}
