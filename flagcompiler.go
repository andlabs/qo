// 24 september 2014

package main

import (
	"flag"
	"os"
	"strings"
	"bufio"
	"os/exec"
)

var debug = flag.Bool("g", false, "build with debug symbols")

var toolchain Toolchain

var cflags []string
var cxxflags []string
var ldflags []string
var libs []string

func pkgconfig(which string, pkgs []string) []string {
	cmd := exec.Command("pkg-config", append([]string{which}, pkgs...)...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		fail("Error runing pkg-config: %v", err)
	}
	return strings.Fields(string(output))
}

func parseFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		fail("Error opening %s to scan for #qo directives: %v", filename, err)
	}
	defer f.Close()
	r := bufio.NewScanner(f)

	for r.Scan() {
		line := r.Text()
		if !strings.HasPrefix(line, "// #qo ") {
			continue
		}
		line = line[len("// #qo "):]
		parts := strings.Fields(line)
		switch parts[0] {
		case "CFLAGS:":
			cflags = append(cflags, parts[1:]...)
		case "CXXFLAGS:":
			cxxflags = append(cxxflags, parts[1:]...)
		case "LDFLAGS:":
			ldflags = append(ldflags, parts[1:]...)
		case "pkg-config:":
			xcflags := pkgconfig("--cflags", parts[1:])
			xlibs := pkgconfig("--libs", parts[1:])
			cflags = append(cflags, xcflags...)
			cxxflags = append(cxxflags, xcflags...)
			ldflags = append(ldflags, xlibs...)
		case "LIBS:":
			for i := 1; i < len(parts); i++ {
				libs = append(libs, parts[i])
			}
		default:
			fail("Invalid #qo directive %q in %s", parts[0], filename)
		}
	}
	if err := r.Err(); err != nil {
		fail("Error reading %s to scan for #qo directives: %v", filename, err)
	}
}

func compileFlags() {
	if *selectedToolchain == "" {
		*selectedToolchain = "gcc"
		if *targetOS == "darwin" {
			*selectedToolchain = "clang"
		}
	}
	toolchain = toolchains[*selectedToolchain]
	toolchain.Prepare()

	cflags = append(cflags, strings.Fields(os.Getenv("CFLAGS"))...)
	cxxflags = append(cxxflags, strings.Fields(os.Getenv("CXXFLAGS"))...)
	ldflags = append(ldflags, strings.Fields(os.Getenv("LDFLAGS"))...)

	for _, f := range cfiles {
		parseFile(f)
	}
	for _, f := range cppfiles {
		parseFile(f)
	}
	for _, f := range hfiles {
		parseFile(f)
	}
	for _, f := range mfiles {
		parseFile(f)
	}
	for _, f := range mmfiles {
		parseFile(f)
	}
	for _, f := range rcfiles {
		parseFile(f)
	}
}
