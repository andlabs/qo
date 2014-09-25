// 24 september 2014

package main

import (
	"fmt"
	"flag"
	"sort"
)

type Toolchain struct {
	CC			string
	CXX			string
	LD			string
	LDCXX		string
	CFLAGS		[]string
	CXXFLAGS	[]string
	LDFLAGS		[]string
	LDCXXFLAGS	[]string	// appended to LDFLAGS if at least one C++ file is present
	IsGCC		bool		// for the flag compiler
}

var toolchains = make(map[string]*Toolchain)

// values for CFLAGS/CXXFLAGS/LDFLAGS shared by all gcc and clang variants
// we specify -Wno-unused-parameter for the case where we are defining an interface and are not using some parameter
// I refuse to support C11.
var gccbase = &Toolchain{
	CFLAGS:		[]string{"--std=c99", "-Wall", "-Wextra", "-Wno-unused-parameter"},
	CXXFLAGS:	[]string{"--std=c++11", "-Wall", "-Wextra", "-Wno-unused-parameter"},
	LDFLAGS:		nil,
}

var gccarchflags = map[string]string{
	"386":		"-m32",
	"amd64":		"-m64",
}

// TODO:
// - MinGW static libgcc/libsjlj/libwinpthread/etc.

func init() {
	toolchains["gcc"] = &Toolchain{
		CC:			"gcc",
		CXX:			"g++",
		LD:			"gcc",
		LDCXX:		"g++",
		IsGCC:		true,
	}
	toolchains["clang"] = &Toolchain{
		CC:			"clang",
		CXX:			"clang++",
		LD:			"clang",
		LDCXX:		"clang++",
		IsGCC:		true,
	}
	// TODO: MinGW cross-compiling, MSVC, Plan 9 compilers
}

var selectedToolchain = flag.String("tc", "",  "select toolchain; list for a full list")

func listToolchains() {
	tc := make([]string, 0, len(toolchains))
	for k := range toolchains {
		tc = append(tc, k)
	}
	sort.Strings(tc)
	for _, t := range tc {
		fmt.Printf("%s\n", t)
	}
}
