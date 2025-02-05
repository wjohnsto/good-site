package file

import (
	"os"
	"path"
	"path/filepath"
	"regexp"

	"good.site/pkg/utils/stack"
)

func EnsurePathExists(filePath string) {
	dir := filePath
	if filepath.Ext(dir) != "" {
		dir = path.Join(dir, "..")
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil && err != os.ErrExist {
		panic(err)
	}
}

func CreateFile(filePath string) *os.File {
	EnsurePathExists(filePath)
	f, err := os.Create(filePath)

	if err != nil {
		panic(err)
	}

	return f
}

func ReadFile(filePath string) []byte {
	EnsurePathExists(filePath)

	f, err := os.ReadFile(filePath)

	if err != nil {
		panic(err)
	}

	return f
}

func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

func WriteFile(filePath string, contents string) {
	f := CreateFile(filePath)

	_, err := f.WriteString(contents)

	if err != nil {
		defer f.Close()
		panic(err)
	}

	defer f.Close()
}

func PrunePaths(paths []string, ext string) []string {
	out := []string{}

	for _, path := range paths {
		if path[len(path)-len(ext):] == ext {
			out = append(out, path)
		}
	}

	return out
}

func FindFiles(dir string, match *regexp.Regexp) []string {
	dirs := stack.New()
	var files []string
	dirs.Push(dir)
	d := dirs.Pop()

	for d != nil {
		items, err := os.ReadDir(d.(string))

		if err != nil {
			panic(err)
		}

		for _, item := range items {
			fullPath := path.Join(d.(string), item.Name())
			if item.IsDir() {
				dirs.Push(fullPath)
			} else if match == nil || match.MatchString(fullPath) {
				files = append(files, fullPath)
			}
		}

		d = dirs.Pop()
	}

	return files
}
