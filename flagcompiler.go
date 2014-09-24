// 24 september 2014

package main

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

	// TODO append flags from environment variables (TODO figure out how to handle quotes)
	// TODO read each file and append flags
}
