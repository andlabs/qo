// 24 september 2014

package main

import (
	"fmt"
	"flag"
	"sort"
)

type Toolchain interface {
	BuildCFile(filename string, cflags []string) (stages []Stage, object string)
	BuildCXXFile(filename string, cflags []string) (stages []Stage, object string)
//	BuildMFile(filename string, cflags []string) (stages []Stage, object string)
//	BuildMMFile(filename string, cflags []string) (stages []Stage, object string)
	BuildRCFile(filename string, cflags []string) (stages []Stage, object string)
	Link(objects []string, ldflags []string, libs []string) *Executor
}

// toolchains[name][arch]
var toolchains = make(map[string]map[string]Toolchain)

/*
TODO
func init() {
	toolchains["msvc"] = make(map[string]*Toolchain)
	toolchains["msvc"]["386"] = &Toolchain{
		CC:				"cl",
		CXX:				"cl",
		LD:				"link",
		LDCXX:			"link",
		RC:				"rc",
		CVTRES:			"cvtres",
		// TODO /bigobj?
		CFLAGS:			[]string{"/c", "/analyze", "/nologo", "/RTC1", "/RTCc", "/RTCs", "/RTCu", "/sdl", "/TC", "/Wall", "/Wp64"},
		CXXFLAGS:		[]string{"/c", "/analyze", "/nologo", "/RTC1", "/RTCc", "/RTCs", "/RTCu", "/sdl", "/TP", "/Wall", "/Wp64"},
		// TODO keep /largeaddressaware?
		LDFLAGS:			[]string{"/largeaddressaware", "/nologo"},
		RCFLAGS:			[]string{"/nologo"},
		CVTRESFLAGS:		[]string{"/nologo"},
		CDEBUG:			[]string{"/Z7"},		// embedded debug information
		LDDEBUG:			nil,				// TODO MSDN claims it's not possible to have embedded debug symbols (apparently COFF doesn't exist)
		COUTPUT:			[]string{"/Fo"},
		LDOUTPUT:		[]string{"/OUT:"},
		RCOUTPUT:		[]string{"/fo", ""},
		CVTRESOUTPUT:	[]string{"/out:"},
		LIBPREFIX:		"",
		LIBSUFFIX:		".lib",
		// TODO resource compiling is a two-step process:
		// 1) rc /nologo /fo file.res file.rc
		// 2) cvtres /nologo /out:file.o file.res
	}
	// TODO: Plan 9 compilers
}
*/

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
