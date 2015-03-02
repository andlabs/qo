// 24 september 2014

package main

import (
	"path/filepath"
	"strings"
)

var script []Stage
var nSteps int

func objectName(filename string, suffix string) string {
	object := strings.Replace(filename, string(filepath.Separator), "_", -1) + suffix
	object = filepath.Join(".qoobj", object)
	return object
}

func mergeScript(s []Stage) {
	var i int

	// first append existing steps
	for i = 0; i < len(s); i++ {
		if i > len(script) {
			break
		}
		script[i] = append(script[i], s[i]...)
		nSteps += len(s[i])
	}
	// now add new stages
	for ; i < len(s); i++ {
		script = append(script, s[i])
		nSteps += len(s[i])
	}
}

func buildScript() {
	script = nil
	nSteps = 0

	objects := []string(nil)

	for _, f := range cfiles {
		s, obj := toolchain.BuildCFile(f, cflags)
		mergeScript(s)
		objects = append(objects, obj)
	}
	for _, f := range cppfiles {
		s, obj := toolchain.BuildCXXFile(f, cxxflags)
		mergeScript(s)
		objects = append(objects, obj)
	}
	for _, f := range mfiles {
		s, obj := toolchain.BuildMFile(f, cflags)
		mergeScript(s)
		objects = append(objects, obj)
	}
	for _, f := range mmfiles {
		s, obj := toolchain.BuildMMFile(f, cxxflags)
		mergeScript(s)
		objects = append(objects, obj)
	}
	for _, f := range rcfiles {
		s, obj := toolchain.BuildRCFile(f, nil)
		mergeScript(s)
		objects = append(objects, obj)
	}

	s := toolchain.Link(objects, ldflags, libs)
	script = append(script, Stage{s})
}
