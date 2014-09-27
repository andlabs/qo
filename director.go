// 24 september 2014

package main

import (
	"path/filepath"
	"strings"
	"sort"
)

type Stage []*Executor

func (s Stage) Len() int {
	return len(s)
}

func (s Stage) Less(i, j int) bool {
	return s[i].Line[1] < s[j].Line[1]
}

func (s Stage) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

var script []Stage
var nStages int

func makeCompileStep(f string, cc string, flags []string, outflag []string, suffix string) (*Executor, string) {
	object := strings.Replace(f, string(filepath.Separator), "_", -1) + suffix
	object = filepath.Join(".qoobj", object)
	e := &Executor{
		Name:	"Compiled " + f,
		Line:		[]string{cc, f},
	}
	e.Line = append(e.Line, flags...)
	e.Line = append(e.Line, outflag...)
	e.Line[len(e.Line) - 1] += object	// append to last flag for msvc
	return e, object
}

func buildScript() {
	script = nil
	nStages = 0

	// stage 1: make the object file directory
	e := &Executor{
		Name:	"Made working directory",
		Line:		[]string{"mkdir", "-p", ".qoobj"},
	}
	script = append(script, Stage{e})
	nStages++

	// stage 2: compile everything
	stage2 := Stage(nil)
	objects := []string(nil)
	linker := toolchain.LD
	for _, f := range cfiles {
		e, object := makeCompileStep(f, toolchain.CC, toolchain.CFLAGS, toolchain.COUTPUT, ".o")
		stage2 = append(stage2, e)
		objects = append(objects, object)
	}
	for _, f := range cppfiles {
		linker = toolchain.LDCXX		// run only if cppfiles isn't empty
		e, object := makeCompileStep(f, toolchain.CXX, toolchain.CXXFLAGS, toolchain.COUTPUT, ".o")
		stage2 = append(stage2, e)
		objects = append(objects, object)
	}
	// TODO .m, .mm files
	resfiles := make([]string, 0, len(rcfiles))
	for _, f := range rcfiles {
		suffix := ".o"
		if toolchain.CVTRES != "" {
			suffix = ".res"
		}
		e, object := makeCompileStep(f, toolchain.RC, toolchain.RCFLAGS, toolchain.RCOUTPUT, suffix)
		stage2 = append(stage2, e)
		if toolchain.CVTRES != "" {
			resfiles = append(resfiles, object)
		} else {
			objects = append(objects, object)
		}
	}
	sort.Sort(stage2)
	sort.Strings(objects)
	script = append(script, stage2)
	nStages += len(stage2)

	// 3) cvtres (msvc)
	stage3 := Stage(nil)
	if len(resfiles) != 0 {
		for _, f := range rcfiles {
			e, object := makeCompileStep(f, toolchain.CVTRES, toolchain.CVTRESFLAGS, toolchain.CVTRESOUTPUT, ".o")
			stage3 = append(stage3, e)
			objects = append(objects, object)
		}
		sort.Sort(stage3)
		sort.Strings(objects)
		script = append(script, stage3)
		nStages += len(stage3)
	}

	// 4) link
	target := targetName()
	e = &Executor{
		Name:	"Linked " + target,
		Line:		make([]string, 0, len(objects) + len(toolchain.LDFLAGS) + 10),
	}
	e.Line = append(e.Line, linker)
	e.Line = append(e.Line, objects...)
	e.Line = append(e.Line, toolchain.LDFLAGS...)
	e.Line = append(e.Line, toolchain.LDOUTPUT...)
	e.Line = append(e.Line, target)
	script = append(script, Stage{e})
	nStages++
}
