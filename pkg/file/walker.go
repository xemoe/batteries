package file

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type File string

type Walker interface {
	List() []File
}

type FileWalker struct {
	BaseDir  string
	FileExt  string
	MaxDepth int
	Walker
}

func (f FileWalker) List() []File {
	files := dirwalk(f.BaseDir, f.FileExt, f.MaxDepth)
	sort.Strings(files)
	return str2file(files)
}

func dirwalk(basedir string, fileext string, maxdepth int) []string {
	files := []string{}
	err := filepath.Walk(basedir, func(path string, info os.FileInfo, err error) error {

		// TODO: refactor errors handler
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != fileext {
			return nil
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		if maxdepth == 0 {
			return nil
		}

		if maxdepth > 0 && getPathDepth(path, string(os.PathSeparator)) > maxdepth {
			return nil
		}

		files = append(files, path)

		return nil
	})

	// TODO: refactor errors handler
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", basedir, err)
		return nil
	}

	return files
}

func getPathDepth(path string, separator string) int {
	return strings.Count(path, separator)
}

func str2file(s []string) []File {
	f := []File{}

	for _, v := range s {
		f = append(f, File(v))
	}

	return f
}
