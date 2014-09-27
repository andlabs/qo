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
	RC			string
	CFLAGS		[]string
	CXXFLAGS	[]string
	LDFLAGS		[]string
	LDCXXFLAGS	[]string	// appended to LDFLAGS if at least one C++ file is present
	CDEBUG		[]string	// appended to CFLAGS *and* CXXFLAGS
	LDDEBUG		[]string
	COUTPUT		[]string	// prepended to output filename on both CFLAGS *and CXXFLAGS
	LDOUTPUT	[]string
	LIBPREFIX		string	// for #qo LIBS: ...
	LIBSUFFIX		string
}

// toolchains[name][arch]
var toolchains = make(map[string]map[string]*Toolchain)

// values for CFLAGS/CXXFLAGS/LDFLAGS shared by all gcc and clang variants
// this is a function so the backing array of each slice is new each time (safe for append)
// specify -c here to keep things clean
// we specify -Wno-unused-parameter for the case where we are defining an interface and are not using some parameter
// I refuse to support C11.
func gcc1(exe *Toolchain, archflag string) *Toolchain {
	return &Toolchain{
		CC:			exe.CC,
		CXX:			exe.CXX,
		LD:			exe.LD,
		LDCXX:		exe.LDCXX,
		RC:			exe.RC,
		CFLAGS:		[]string{"-c", "--std=c99", "-Wall", "-Wextra", "-Wno-unused-parameter", archflag},
		CXXFLAGS:	[]string{"-c", "--std=c++11", "-Wall", "-Wextra", "-Wno-unused-parameter", archflag},
		LDFLAGS:		[]string{archflag},
		CDEBUG:		[]string{"-g"},
		LDDEBUG:		[]string{"-g"},
		COUTPUT:		[]string{"-o"},
		LDOUTPUT:	[]string{"-o"},
		LIBPREFIX:	"-l",
		LIBSUFFIX:	"",
	}
}

// TODO:
// - MinGW static libgcc/libsjlj/libwinpthread/etc.

func gcc(exe *Toolchain) map[string]*Toolchain {
	m := make(map[string]*Toolchain)
	m["386"] = gcc1(exe, "-m32")
	m["amd64"] = gcc1(exe, "-m64")
	return m
}

func init() {
	toolchains["gcc"] = gcc(&Toolchain{
		CC:			"gcc",
		CXX:			"g++",
		LD:			"gcc",
		LDCXX:		"g++",
		RC:			"windres",
	})
	toolchains["clang"] = gcc(&Toolchain{
		CC:			"clang",
		CXX:			"clang++",
		LD:			"clang",
		LDCXX:		"clang++",
		// TODO rc
	})
	toolchains["mingwcc"] = gcc(&Toolchain{
		CC:			"i686-w64-mingw32-gcc",
		CXX:			"i686-w64-mingw32-g++",
		LD:			"i686-w64-mingw32-gcc",
		LDCXX:		"i686-w64-mingw32-g++",
		RC:			"i686-w64-mingw32-windres",
	})
	// patch up the amd64 mingw cross compiler to use the correct executable names
	mingw64 := toolchains["mingwcc"]["amd64"]
	mingw64.CC = "x86_64-w64-mingw32-gcc"
	mingw64.CXX = "x86_64-w64-mingw32-g++"
	mingw64.LD = "x86_64-w64-mingw32-gcc"
	mingw64.LDCXX = "x86_64-w64-mingw32-g++"
	mingw64.RC = "x86_64-w64-mingw32-windres"

	toolchain["msvc"]["UNFINISHED"] = &Toolchain{
		CC:			"cl",
		CXX:			"cl",
		LD:			"link",
		LDCXX:		"link",
		RC:			"rc",
		// TODO /bigobj?
		CFLAGS:		[]string{"/c", "/analyze", "/nologo", "/RTC1", "/RTCc", "/RTCs", "/RTCu", "/sdl", "/TC", "/Wall", "/Wp64"},
		CXXFLAGS:	[]string{"/c", "/analyze", "/nologo", "/RTC1", "/RTCc", "/RTCs", "/RTCu", "/sdl", "/TP", "/Wall", "/Wp64"},
		// TODO keep /largeaddressaware?
		LDFLAGS:		[]string{"/largeaddressaware", "/nologo"},
		CDEBUG:		[]string{"/Z7"},		// embedded debug information
		LDDEBUG:		nil,				// TODO MSDN claims it's not possible to have embedded debug symbols (apparently COFF doesn't exist)
		COUTPUT:		[]string{"/Fo"},		// TODO is one argument
		LDOUTPUT:	[]string{"/OUT:"},	// TODO is one argument
		LIBPREFIX:	"",
		LIBSUFFIX:	".lib",
		// TODO resource compiling is a two-step process:
		// 1) rc /nologo /fo file.res file.rc
		// 2) cvtres /nologo /out:file.o file.res
	}
	// TODO: Plan 9 compilers
}

var selectedToolchain = flag.String("tc", "",  "select toolchain; list for a full list")

func listToolchains() {
	tc := make([]string, 0, len(toolchains))
	for k, _ := range toolchains {
		tc = append(tc, k)
	}
	sort.Strings(tc)
	for _, t := range tc {
		fmt.Printf("%s\n", t)
	}
}
