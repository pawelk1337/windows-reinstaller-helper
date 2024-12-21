package main

import (
	"os"
	"path"
	"strings"

	cp "github.com/otiai10/copy"
)

func Resotre(basePath string) {
	filesJson := LoadJson(path.Join(basePath, "files.json"))
	os.MkdirAll(path.Join(basePath, "resotred"), os.ModePerm)

	opts := cp.Options{
		PreserveTimes: true,
		PreserveOwner: true,
		OnSymlink: func(src string) cp.SymlinkAction {
			return cp.Deep
		},
		OnDirExists: func(src, dest string) cp.DirExistsAction {
			panic("DIR EXISTS TRY AGAIN")
		},
		Skip: func(srcinfo os.FileInfo, src, dest string) (bool, error) {
			return path.Ext(src) == "skip", nil
		},
	}

	for _, file := range filesJson.Files {
		// Copy the file to the original location
		if file.Restore {
			err := cp.Copy(file.NewPath, file.OriginalPath, opts)
			if err != nil {
				panic(err)
			}
		} else {
			split := strings.Split(file.OriginalPath, "\\")
			name := split[len(split)-1]
			err := cp.Copy(file.NewPath, path.Join(basePath, "resotred", name), opts)
			if err != nil {
				panic(err)
			}
		}
	}

	var names string

	for _, name := range filesJson.Names {
		split := strings.Split(name, "\\")
		names += split[len(split)-1] + "\n"
	}

	if names != "" {
		err := os.WriteFile(path.Join(basePath, "names.txt"), []byte(names), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	os.RemoveAll(path.Join(basePath, "files"))
	os.Remove(path.Join(basePath, "files.json"))
}
