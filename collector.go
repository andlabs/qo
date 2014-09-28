// 24 september 2014

package main

import (
	"os"
	"strings"
	"path/filepath"
)

var cfiles []string
var cppfiles []string
var hfiles []string
var mfiles []string
var mmfiles []string
var rcfiles []string

var base string

func consider(list *[]string, path string) {
	if excludeFile(path) {
		return
	}
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
		if excludeDir(path) {
			return filepath.SkipDir
		}
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
	case ".rc":
		consider(&rcfiles, path)
	}
	return nil
}
