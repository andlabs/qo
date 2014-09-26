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
	CDEBUG		[]string	// appended to CFLAGS *and* CXXFLAGS
	LDDEBUG		[]string
	COUTPUT		[]string	// prepended to output filename on both CFLAGS *and CXXFLAGS
	LDOUTPUT	[]string
}

var toolchains = make(map[string]*Toolchain)

// values for CFLAGS/CXXFLAGS/LDFLAGS shared by all gcc and clang variants
// specify -c here to keep things clean
// we specify -Wno-unused-parameter for the case where we are defining an interface and are not using some parameter
// I refuse to support C11.
var gccbase = &Toolchain{
	CFLAGS:		[]string{"-c", "--std=c99", "-Wall", "-Wextra", "-Wno-unused-parameter"},
	CXXFLAGS:	[]string{"-c", "--std=c++11", "-Wall", "-Wextra", "-Wno-unused-parameter"},
	LDFLAGS:		nil,
	CDEBUG:		[]string{"-g"},
	LDDEBUG:		[]string{"-g"},
	COUTPUT:		[]string{"-o"},
	LDOUTPUT:	[]string{"-o"},
}

var gccarchflags = map[string]string{
	"386":		"-m32",
	"amd64":		"-m64",
}

// TODO:
// - MinGW static libgcc/libsjlj/libwinpthread/etc.
// - simplify the below

// TODO bad for init()
func gcc(t *Toolchain) {
	t.CFLAGS = append(t.CFLAGS, gccbase.CFLAGS...)
	t.CFLAGS = append(t.CFLAGS, gccarchflags[*targetArch])
	t.CXXFLAGS = append(t.CXXFLAGS, gccbase.CXXFLAGS...)
	t.CXXFLAGS = append(t.CXXFLAGS, gccarchflags[*targetArch])
	t.LDFLAGS = append(t.LDFLAGS, gccbase.LDFLAGS...)
	t.LDFLAGS = append(t.LDFLAGS, gccarchflags[*targetArch])
	t.CDEBUG = append(t.CDEBUG, gccbase.CDEBUG...)
	t.LDDEBUG = append(t.LDDEBUG, gccbase.LDDEBUG...)
	t.COUTPUT = append(t.COUTPUT, gccbase.COUTPUT...)
	t.LDOUTPUT = append(t.LDOUTPUT, gccbase.LDOUTPUT...)
}

func init() {
	toolchains["gcc"] = &Toolchain{
		CC:			"gcc",
		CXX:			"g++",
		LD:			"gcc",
		LDCXX:		"g++",
	}
	gcc(toolchains["gcc"])
	toolchains["clang"] = &Toolchain{
		CC:			"clang",
		CXX:			"clang++",
		LD:			"clang",
		LDCXX:		"clang++",
	}
	gcc(toolchains["clang"])
	// TOOD merge this with -arch
	toolchains["mingwcc32"] = &Toolchain{
		CC:			"i686-w64-mingw32-gcc",
		CXX:			"i686-w64-mingw32-g++",
		LD:			"i686-w64-mingw32-gcc",
		LDCXX:		"i686-w64-mingw32-g++",
	}
	gcc(toolchains["mingwcc32"])
	toolchains["mingwcc64"] = &Toolchain{
		CC:			"x86_64-w64-mingw32-gcc",
		CXX:			"x86_64-w64-mingw32-g++",
		LD:			"x86_64-w64-mingw32-gcc",
		LDCXX:		"x86_64-w64-mingw32-g++",
	}
	gcc(toolchains["mingwcc64"])
	// TODO: MSVC, Plan 9 compilers
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
