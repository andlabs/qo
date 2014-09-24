// 24 september 2014

package main

import (
	"flag"
)

type Toolchain struct {
	CC			string
	CPP			string
	LD			string
	CFLAGS		[]string
	CPPFLAGS		[]string
	LDFLAGS		[]string
	LDCPPFLAGS	[]string	// appended to LDFLAGS if at least one C++ file is present
}

var toolchains = make(map[string]*Toolchain)

// values for CFLAGS/CPPFLAGS/LDFLAGS shared by all gcc and clang variants
// we specify -Wno-unused-parameter for the case where we are defining an interface and are not using some parameter
// I refuse to support C11.
var gccbase = &Toolchain{
	CFLAGS:		[]string{"--std=c99", "-Wall", "-Wextra", "-Wno-unused-parameter"},
	CPPFLAGS:	[]string{"--std=c++11", "-Wall", "-Wextra", "-Wno-unused-parameter"},
	LDFLAGS:		nil,
}

// TODO:
// - MinGW static libgcc/libsjlj/libwinpthread/etc.
// - CXX instead of CPP?

func init() {
	toolchains["gcc"] = &Toolchain{
		CC:			"gcc",
		CPP:			"g++",
		LD:			"gcc",
		CFLAGS:		gccbase.CFLAGS,
		CPPFLAGS:	gccbase.CPPFLAGS,
		LDFLAGS:		gccbase.LDFLAGS,
	}
	toolchains["clang"] = &Toolchain{
		CC:			"clang",
		CPP:			"clang++",
		LD:			"clang",
		CFLAGS:		gccbase.CFLAGS,
		CPPFLAGS:	gccbase.CPPFLAGS,
		LDFLAGS:		gccbase.LDFLAGS,
	}
	// TODO: MinGW cross-compiling, MSVC, Plan 9 compilers
}

var selectedToolchain *Toolchain
var selectToolchain = flag.String("tc", "",  "select toolchain; list for a full list")
