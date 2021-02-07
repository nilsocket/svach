package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	color "github.com/logrusorgru/aurora/v3"
	svach "github.com/nilsocket/svach"
	"github.com/xlab/treeprint"
)

var (
	recursive bool
	change    bool
	name      bool
	rootDir   string
	treeMap   map[string]treeprint.Tree
)

var separator = string(filepath.Separator)
var arrow = "━━▶ "

func init() {
	flag.BoolVar(&name, "n", false, "name")
	flag.BoolVar(&recursive, "r", false, "recursive")
	flag.BoolVar(&change, "c", false, "change fileNames, used in conjunction")
}

var s *svach.Svach

func main() {

	flag.Parse()

	s = svach.New()

	args := flag.Args()

	if len(args) > 0 {
		for _, n := range args {

			var nn string

			if !name {
				nn = s.Clean(n)
			} else {
				nn = s.Name(n)
			}

			fmt.Println(nn)
		}
	} else {

		rootDir, _ = os.Getwd()
		if rootDir != "/" {
			rootDir += string(filepath.Separator)
		}

		treeMap = make(map[string]treeprint.Tree)
		treeMap[rootDir] = treeprint.New()

		filepath.Walk(rootDir, rename)

		fmt.Print(treeMap[rootDir].String())
	}

}

func rename(old string, info os.FileInfo, err error) error {

	if err != nil {
		return nil
	}

	var new string

	dir, file := filepath.Split(old)

	if !recursive {
		// End, when trying to walk another directory
		if dir != rootDir {
			return filepath.SkipDir
		}
	}

	if file != "" {
		if !name {
			new = s.Clean(file)
		} else {
			new = s.Name(file)
		}

		if file != new {
			if change {

				filePath := filepath.Join(dir, file)
				newPath := filepath.Join(dir, new)

				if err = os.Rename(filePath, newPath); err != nil {
					log.Println(file, err)
				}
			}

			if treeMap[dir] == nil {
				createLink(dir)
			}

			if info.IsDir() {
				treeMap[old+separator] = treeMap[dir].AddBranch(
					fileChange(file, new),
				)
			} else {
				treeMap[dir].AddNode(
					fileChange(file, new),
				)
			}
		}
	}

	return nil
}

func createLink(dir string) {
	roots := strings.SplitAfter(rootDir, separator)
	dirs := strings.SplitAfter(dir, separator)

	rootLen := len(roots) - 1

	roots = roots[:rootLen]
	dirs = dirs[:len(dirs)-1]

	prevDir := rootDir
	for i := range dirs[rootLen:] {
		curDir := prevDir + dirs[rootLen+i]
		if treeMap[curDir] == nil {
			_, file := filepath.Split(curDir[:len(curDir)-1])
			treeMap[curDir] = treeMap[prevDir].AddBranch(color.Cyan(file))
		}
		prevDir = curDir
	}
}

func fileChange(file, new string) string {
	return fmt.Sprint(color.Red(file), arrow, color.Green(new))
}
