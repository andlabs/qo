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
		toolchain.CPPFLAGS = append(toolchain.CPPFLAGS, gccbase.CPPFLAGS...)
		toolchain.CPPFLAGS = append(toolchain.CPPFLAGS, gccarchflags[*targetArch])
		toolchain.LDFLAGS = append(toolchain.LDFLAGS, gccbase.LDFLAGS...)
		toolchain.LDFLAGS = append(toolchain.LDFLAGS, gccarchflags[*targetArch])
	}

	toolchain.CFLAGS = append(toolchain.CFLAGS, strings.Fields(os.Getenv("CFLAGS"))...)
	toolchain.CPPFLAGS = append(toolchain.CPPFLAGS, strings.Fields(os.Getenv("CPPFLAGS"))...)
	toolchain.LDFLAGS = append(toolchain.LDFLAGS, strings.Fields(os.Getenv("LDFLAGS"))...)

	// TODO read each file and append flags
}
