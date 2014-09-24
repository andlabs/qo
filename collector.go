// 24 september 2014

package main

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"
)

var cfiles []string
var cppfiles []string
var hfiles []string
var mfiles []string
var mmfiles []string

var base string

func consider(list *[]string, path string) {
	// TODO operating system filters
	*list = append(*list, path)
}

func walker(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if base := filepath.Base(path); base != "." && base != ".." && base[0] == '.'  {
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}
	if info.IsDir() {
		return nil
	}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".c":
		consider(&cfiles, path)
	case ".cpp", ".cxx", ".c++", ".cc":
		consider(&cppfiles, path)
	case ".h", ".hpp", ".hxx", ".h++", ".hh":
		consider(&hfiles, path)
	case ".m":
		consider(&mfiles, path)
	case ".mm":
		consider(&mmfiles, path)
	}
	return nil
}

func main() {
//	flags.Parse()
	err := filepath.Walk(".", walker)
	if err != nil {
		panic(err)
	}
	compileFlags()
	buildScript()
	for _, s := range script {
		for _, e := range *s {
			fmt.Printf("%#v\n", e)
		}
	}
}
