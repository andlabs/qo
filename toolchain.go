// 24 september 2014

package main

import (
	"fmt"
	"flag"
	"sort"
)

type Toolchain interface {
	Prepare()
	BuildCFile(filename string, cflags []string) (stages []Stage, object string)
	BuildCXXFile(filename string, cflags []string) (stages []Stage, object string)
	BuildMFile(filename string, cflags []string) (stages []Stage, object string)
	BuildMMFile(filename string, cflags []string) (stages []Stage, object string)
	BuildRCFile(filename string, cflags []string) (stages []Stage, object string)
	Link(objects []string, ldflags []string, libs []string) *Step
}

var toolchains = make(map[string]Toolchain)

// TODO: Plan 9 compilers (plan9)

var selectedToolchain = flag.String("tc", "",  "select toolchain; list for a full list of supported toolchains")

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
