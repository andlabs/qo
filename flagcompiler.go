// 24 september 2014

package main

import (
	"os"
	"strings"
)

var toolchain Toolchain

func compileFlags() {
	if *selectedToolchain == "" {
		*selectedToolchain = "gcc"
		if *targetOS == "darwin" {
			*selectedToolchain = "clang"
		}
	}

	// copy the initial values
	toolchain = *(toolchains[*selectedToolchain])

	if toolchain.IsGCC {
		toolchain.CFLAGS = append(toolchain.CFLAGS, gccbase.CFLAGS...)
		toolchain.CFLAGS = append(toolchain.CFLAGS, gccarchflags[*targetArch])
		toolchain.CXXFLAGS = append(toolchain.CXXFLAGS, gccbase.CXXFLAGS...)
		toolchain.CXXFLAGS = append(toolchain.CXXFLAGS, gccarchflags[*targetArch])
		toolchain.LDFLAGS = append(toolchain.LDFLAGS, gccbase.LDFLAGS...)
		toolchain.LDFLAGS = append(toolchain.LDFLAGS, gccarchflags[*targetArch])
	}

	toolchain.CFLAGS = append(toolchain.CFLAGS, strings.Fields(os.Getenv("CFLAGS"))...)
	toolchain.CXXFLAGS = append(toolchain.CXXFLAGS, strings.Fields(os.Getenv("CXXFLAGS"))...)
	toolchain.LDFLAGS = append(toolchain.LDFLAGS, strings.Fields(os.Getenv("LDFLAGS"))...)

	// TODO read each file and append flags
}
