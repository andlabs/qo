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

func pkgconfig(which string, pkgs []string) []string {
	cmd := exec.Command("pkg-config", append([]string{which}, pkgs...)...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		// TODO
		panic(err)
	}
	return strings.Fields(string(output))
}

func parseFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		// TODO
		panic(err)
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
			toolchain.CFLAGS = append(toolchain.CFLAGS, parts[1:]...)
		case "CXXFLAGS:":
			toolchain.CXXFLAGS = append(toolchain.CXXFLAGS, parts[1:]...)
		case "LDFLAGS:":
			toolchain.LDFLAGS = append(toolchain.LDFLAGS, parts[1:]...)
		case "pkg-config:":
			cflags := pkgconfig("--cflags", parts[1:])
			libs := pkgconfig("--libs", parts[1:])
			toolchain.CFLAGS = append(toolchain.CFLAGS, cflags...)
			toolchain.CXXFLAGS = append(toolchain.CXXFLAGS, cflags...)
			toolchain.LDFLAGS = append(toolchain.LDFLAGS, libs...)
		case "LIBS:":
			for i := 1; i < len(parts); i++ {
				toolchain.LDFLAGS = append(toolchain.LDFLAGS, toolchain.LIBPREFIX + parts[i] + toolchain.LIBSUFFIX)
			}
		default:
			// TODO
			panic("invalid line")
		}
	}
	if err := r.Err(); err != nil {
		// TODO
		panic(err)
	}
}

func compileFlags() {
	if *selectedToolchain == "" {
		*selectedToolchain = "gcc"
		if *targetOS == "darwin" {
			*selectedToolchain = "clang"
		}
	}

	// copy the initial values
	toolchain = *(toolchains[*selectedToolchain][*targetArch])

	toolchain.CFLAGS = append(toolchain.CFLAGS, strings.Fields(os.Getenv("CFLAGS"))...)
	toolchain.CXXFLAGS = append(toolchain.CXXFLAGS, strings.Fields(os.Getenv("CXXFLAGS"))...)
	toolchain.LDFLAGS = append(toolchain.LDFLAGS, strings.Fields(os.Getenv("LDFLAGS"))...)

	// TODO read each file and append flags

	if *debug {
		toolchain.CFLAGS = append(toolchain.CFLAGS, toolchain.CDEBUG...)
		toolchain.CXXFLAGS = append(toolchain.CXXFLAGS, toolchain.CDEBUG...)
		toolchain.LDFLAGS = append(toolchain.LDFLAGS, toolchain.LDDEBUG...)
	}

	for _, f := range cfiles {
		parseFile(f)
	}
	for _, f := range cppfiles {
		parseFile(f)
	}
	for _, f := range hfiles {
		parseFile(f)
	}
	// TODO .m, .mm
	for _, f := range rcfiles {
		parseFile(f)
	}
}
