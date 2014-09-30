// 27 september 2014

package main

import (
	"flag"
	"runtime"
)

type GCCBase struct {
	CC		string
	CXX		string
	LD		string
	RC		string
	ArchFlag	[]string
}

func (g *GCCBase) buildRegularFile(cc string, std string, cflags []string, filename string) (stages []Stage, object string) {
	object = objectName(filename, ".o")
	line := []string{
		cc,
		filename,
		"-c",
		std,
		"-Wall",
		"-Wextra",
		// for the case where we are implementing an interface and are not using some parameter
		"-Wno-unused-parameter",
	}
	if g.ArchFlag != nil {
		line = append(line, g.ArchFlag...)
	}
	line = append(line, cflags...)
	if *debug {
		line = append(line, "-g")
	}
	line = append(line, "-o", object)
	e := &Executor{
		Name:	"Compiled " + filename,
		Line:		line,
	}
	stages = []Stage{
		nil,
		Stage{e},
		nil,
	}
	return stages, object
}

func (g *GCCBase) BuildCFile(filename string, cflags []string) (stages []Stage, object string) {
	return g.buildRegularFile(
		g.CC,
		"--std=c99",		// I refuse to support C11.
		cflags,
		filename)
}

func (g *GCCBase) BuildCXXFile(filename string, cflags []string) (stages []Stage, object string) {
	g.LD = g.CXX
	return g.buildRegularFile(
		g.CXX,
		"--std=c++11",
		cflags,
		filename)
}

// apart from needing -lobjc at link time, Objective-C/C++ are identical to C/C++; the --std flags are the same (thanks Beelsebob in irc.freenode.net/#macdev)
// TODO provide -lobjc?

func (g *GCCBase) BuildMFile(filename string, cflags []string) (stages []Stage, object string) {
	return g.BuildCFile(filename, cflags)
}

func (g *GCCBase) BuildMMFile(filename string, cflags []string) (stages []Stage, object string) {
	return g.BuildCXXFile(filename, cflags)
}

func (g *GCCBase) BuildRCFile(filename string, cflags []string) (stages []Stage, object string) {
	if g.RC == "" {
		fail("LLVM/clang does not come with a Windows resource compiler (if this message appears in other situations in error, contact andlabs)")
	}
	object = objectName(filename, ".o")
	line := append([]string{
		g.RC,
		filename,
		object,
	}, cflags...)
	e := &Executor{
		Name:	"Compiled " + filename,
		Line:		line,
	}
	stages = []Stage{
		nil,
		Stage{e},
		nil,
	}
	return stages, object
}

func (g *GCCBase) Link(objects []string, ldflags []string, libs []string) *Executor {
	if g.LD == "" {
		g.LD = g.CC
	}
	target := targetName()
	for i := 0; i < len(libs); i++ {
		libs[i] = "-l" + libs[i]
	}
	line := []string{
		g.LD,
	}
	if g.ArchFlag != nil {
		line = append(line, g.ArchFlag...)
	}
	line = append(line, objects...)
	line = append(line, ldflags...)
	line = append(line, libs...)
	if *debug {
		line = append(line, "-g")
	}
	line = append(line, "-o", target)
	return &Executor{
		Name:	"Linked " + target,
		Line:		line,
	}
}

// TODO:
// - MinGW static libgcc/libsjlj/libwinpthread/etc.

var triplet = flag.String("triplet", "", "gcc/clang target triplet to use; see README")

var garchs = map[string]string{
	"386":		"i686",
	"amd64":		"x86_64",
}

// 386 and amd64 are commonly configured using multilib rather than targets
// everything else will be a nil slice
var garchflags = map[string][]string{
	"386":		[]string{"-m32"},
	"amd64":		[]string{"-m64"},
}

type GCC struct {
	*GCCBase
}

func isMultilib() bool {
	if *targetOS == runtime.GOOS && runtime.GOOS != "windows" {
		// MinGW for Windows is not multilib
		// assume everything else is
		if runtime.GOARCH == "386" || runtime.GOARCH == "amd64" {
			if *targetArch == "386" || *targetArch == "amd64" {
				// multilib only
				return true
			}
		}
	}
	return false
}

func (g *GCC) Prepare() {
	garch := garchs[*targetArch]
	prefix := ""
	if *triplet != "" {
		prefix = *triplet + "-"
		goto out
	}
	// set this before any of the following in case target == host
	g.ArchFlag = garchflags[*targetArch]
	if *targetOS == runtime.GOOS && *targetArch == runtime.GOARCH {
		return
	}
	if isMultilib() {
		return
	}
	switch *targetOS {
	case "windows":
		prefix = garch + "-w64-mingw32-"
	case "linux":
		// TODO abi override
		prefix = garch + "-linux-gnu-"
	default:
		fail("Sorry, cross-compiling for gcc on %s requires specifying an explicit target triple with -target", *targetOS)
	}
out:
	g.CC = prefix + g.CC
	g.CXX = prefix + g.CXX
	g.RC = prefix + g.RC
}

type Clang struct {
	*GCCBase
}

// see the LLVM source; file include/llvm/ADT/Triple.h (thanks jroelofs in irc.oftc.net/#llvm)
var clangOS = map[string]string{
	"windows":	"win32",
	"linux":		"linux",
	"darwin":		"darwin",
	"freebsd":		"freebsd",
	"openbsd":	"openbsd",
	"netbsd":		"netbsd",
	"dragonfly":	"dragonfly",
	"solaris":		"solaris",
}

func (g *Clang) Prepare() {
	if *triplet != "" {
		g.ArchFlag = []string{"-target", *triplet}
		return
	}
	// set this before any of the following in case target == host
	g.ArchFlag = garchflags[*targetArch]
	if *targetOS == runtime.GOOS && *targetArch == runtime.GOARCH {
		return
	}
	if isMultilib() {
		return
	}
	// clang makes the job easier
	// TODO abi override
	g.ArchFlag = append(g.ArchFlag, "-target", garchs[*targetArch] + "-" + clangOS[*targetOS])
}

func init() {
	toolchains["gcc"] = &GCC{
		GCCBase:		&GCCBase{
			CC:		"gcc",
			CXX:		"g++",
			RC:		"windres",
		},
	}
	toolchains["clang"] = &Clang{
		GCCBase:		&GCCBase{
			CC:		"clang",
			CXX:		"clang++",
			// no RC in clang
		},
	}
}
