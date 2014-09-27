// 27 september 2014

package main

import (
	// ...
)

type GCC struct {
	CC		string
	CXX		string
	LD		string
	LDCXX	string
	RC		string
	ArchFlag	string
}

func (g *GCC) buildRegularFile(cc string, std string, cflags []string, filename string) (stages Stage, object string) {
	// TODO split out
	object = strings.Replace(filename, string(filepath.Separator), "_", -1) + suffix
	object = filepath.Join(".qoobj", object)
	line := append([]string{
		cc,
		filename,
		"-c",
		std,
		"-Wall",
		"-Wextra",
		// for the case where we are implementing an interface and are not using some parameter
		"-Wno-unused-parameter",
		g.ArchFlag,
	}, cflags)
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
		nil
	}
	return stages, object
}

func (g *GCC) BuildCFile(filename string, cflags []string) (stages []Stage, object string) {
	return g.buildRegularFile(
		g.CC,
		"--std=c99",		// I refuse to support C11.
		cflags,
		filename)
}

func (g *GCC) BuildCXXFile(filename string, cflags []string) (stages []Stage, object string) {
	g.LD = g.LDCXX
	return g.buildRegularFile(
		g.CXX,
		"--std=c++11",
		cflags,
		filename)
}

// TODO .m, .mm

func (g *GCC) BuildRCFile(filename string, cflags []string) (stages Stage, object string) {
	// TODO split out
	object = strings.Replace(filename, string(filepath.Separator), "_", -1) + suffix
	object = filepath.Join(".qoobj", object)
	line := append([]string{
		g.RC,
		filename,
		object,
	}, cflags)
	e := &Executor{
		Name:	"Compiled " + filename,
		Line:		line,
	}
	stages = []Stage{
		nil,
		Stage{e},
		nil
	}
	return stages, object
}

func (g *GCC) Link(objects []string, ldflags []string, libs []string) *Executor {
	target := targetName()
	for i := 0; i < len(libs); i++ {
		libs[i] = "-l" + libs[i]
	}
	line := append([]string{
		g.LD,
		g.ArchFlag,
	}, objects...)
	line = append(line, ldflags...)
	line = append(line, libs...)
	if *debug {
		line = append(line, "-g")
	}
	line = append(line, "-o", target)
	return &Executor{
		Name:	"Linking " + target,
		Line:		line,
	}
}
