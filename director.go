// 24 september 2014

package main

import (
	"path/filepath"
	"strings"
)

type Stage []*Executor

var script []Stage
var nStages int

func objectName(filename string, suffix string) string {
	object := strings.Replace(filename, string(filepath.Separator), "_", -1) + suffix
	object = filepath.Join(".qoobj", object)
	return object
}

func buildScript() {
	script = nil
	nStages = 0

	stage1 := Stage(nil)
	stage2 := Stage(nil)
	stage3 := Stage(nil)
	objects := []string(nil)

	for _, f := range cfiles {
		s, obj := toolchain.BuildCFile(f, cflags)
		stage1 = append(stage1, s[0]...)
		stage2 = append(stage2, s[1]...)
		stage3 = append(stage3, s[2]...)
		objects = append(objects, obj)
	}
	for _, f := range cppfiles {
		s, obj := toolchain.BuildCXXFile(f, cxxflags)
		stage1 = append(stage1, s[0]...)
		stage2 = append(stage2, s[1]...)
		stage3 = append(stage3, s[2]...)
		objects = append(objects, obj)
	}
	for _, f := range mfiles {
		s, obj := toolchain.BuildMFile(f, cflags)
		stage1 = append(stage1, s[0]...)
		stage2 = append(stage2, s[1]...)
		stage3 = append(stage3, s[2]...)
		objects = append(objects, obj)
	}
	for _, f := range mmfiles {
		s, obj := toolchain.BuildMMFile(f, cxxflags)
		stage1 = append(stage1, s[0]...)
		stage2 = append(stage2, s[1]...)
		stage3 = append(stage3, s[2]...)
		objects = append(objects, obj)
	}
	for _, f := range rcfiles {
		s, obj := toolchain.BuildRCFile(f, nil)
		stage1 = append(stage1, s[0]...)
		stage2 = append(stage2, s[1]...)
		stage3 = append(stage3, s[2]...)
		objects = append(objects, obj)
	}

	if len(stage1) > 0 {
		script = append(script, stage1)
		nStages += len(stage1)
	}
	if len(stage2) > 0 {
		script = append(script, stage2)
		nStages += len(stage2)
	}
	if len(stage3) > 0 {
		script = append(script, stage3)
		nStages += len(stage3)
	}

	e := toolchain.Link(objects, ldflags, libs)
	script = append(script, Stage{e})
	nStages++
}
