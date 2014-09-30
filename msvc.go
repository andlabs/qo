// 27 september 2014

package main

import (
	// ...
)

type MSVC struct			{                                         }

// TODO /MDd, /MTd

func (m *MSVC) buildRegularFile(std string, cflags []string, filename string) (stages []Stage, object string) {
	object = objectName(filename, ".o")
	line := append([]string{
		"cl",
		filename,
		"/c",
		std,
		"/analyze",
		"/bigobj",
		"/nologo",
		"/RTC1",
		"/RTCc",
		"/RTCs",
		"/RTCu",
		"/sdl",
		"/Wall",
		"/Wp64",
	}, cflags...)
	if *debug {
		line = append(line, "/Z7")			// keep debug information in the object file
	}
	line = append(line, "/Fo" + object)		// note: one parameter
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

func (m *MSVC) BuildCFile(filename string, cflags []string) (stages []Stage, object string) {
	return m.buildRegularFile(
		"/TC",
		cflags,
		filename)
}

func (m *MSVC) BuildCXXFile(filename string, cflags []string) (stages []Stage, object string) {
	return m.buildRegularFile(
		"/TP",
		cflags,
		filename)
}

func (m *MSVC) BuildMFile(filename string, cflags []string) (stages []Stage, object string) {
	fail("Objective-C unimplemented on MSVC")
	panic("unreachable")
}

func (m *MSVC) BuildMMFile(filename string, cflags []string) (stages []Stage, object string) {
	fail("Objective-C++ unimplemented on MSVC")
	panic("unreachable")
}

func (m *MSVC) BuildRCFile(filename string, cflags []string) (stages []Stage, object string) {
	resfile := objectName(filename, ".res")
	object = objectName(filename, ".o")
	rcline := append([]string{
		"rc",
		// for rc, flags /must/ come first
		"/nologo",
		"/fo", resfile,		// note: two parameters
	}, cflags...)
	rcline = append(rcline, filename)
	e := &Executor{
		Name:	"Created RES file from " + filename,
		Line:		rcline,
	}
	cvtline := append([]string{
		"cvtres",
		// for cvtres, flags /must/ come first
		"/nologo",
		"/out:" + object,	// note: one parameter
		resfile,
	}, cflags...)
	f := &Executor{
		Name:	"Compiled object file from " + filename,
		Line:		cvtline,
	}
	stages = []Stage{
		nil,
		Stage{e},
		Stage{f},
	}
	return stages, object
}

func (m *MSVC) Link(objects []string, ldflags []string, libs []string) *Executor {
	target := targetName()
	for i := 0; i < len(libs); i++ {
		libs[i] = libs[i] + ".lib"
	}
	line := append([]string{
		"link",
		"/largeaddressaware",
		"/nologo",
	}, objects...)
	line = append(line, ldflags...)
	line = append(line, libs...)
	if *debug {
		// TODO MSDN claims it's not possible to have embedded debug symbols (apparently COFF doesn't exist)
	}
	line = append(line, "/OUT:" + target)			// note: one parameter
	return &Executor{
		Name:	"Linked " + target,
		Line:		line,
	}
}

func (m *MSVC) Prepare() {
	// nothing to do here
}

func init() {
	toolchains["msvc"] = &MSVC{}
}
