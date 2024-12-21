package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"

	cp "github.com/otiai10/copy"
)

type File struct {
	OriginalPath string `json:"originalpath"`
	NewPath      string `json:"newpath"`
	Restore      bool   `json:"restore"`
	IsDir        bool   `json:"is_dir"`
}

type FilesJson struct {
	Files      []File   `json:"files"`
	Names      []string `json:"names"`
	configPath string   `json:"-"`
}

func (fj *FilesJson) Save() {
	data, err := json.MarshalIndent(*fj, "", "\t")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(fj.configPath, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func LoadJson(filePath string) FilesJson {
	// Check if exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return FilesJson{configPath: filePath}
	}

	// Read the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var jsonFiles FilesJson
	err = json.Unmarshal([]byte(data), &jsonFiles)
	if err != nil {
		panic(err)
	}
	jsonFiles.configPath = filePath

	return jsonFiles
}

func SaveName(basePath, name string) {
	filesJson := LoadJson(path.Join(basePath, "files.json"))
	filesJson.Names = append(filesJson.Names, name)
	filesJson.Save()
}

func SaveFile(basePath, filepath string, noSaved, restore bool) {
	// Generate a random 8 character string
	str := RandStringRunes(8)

	newFilePath := path.Join(basePath, "files", str+".wrh")

	opts := cp.Options{
		PreserveTimes: true,
		PreserveOwner: true,
		OnSymlink: func(src string) cp.SymlinkAction {
			return cp.Deep
		},
		OnDirExists: func(src, dest string) cp.DirExistsAction {
			panic("FILE EXISTS TRY AGAIN")
		},
		Skip: func(srcinfo os.FileInfo, src, dest string) (bool, error) {
			return path.Ext(src) == "skip", nil
		},
	}

	// Copy the file to the new location
	err := cp.Copy(filepath, newFilePath, opts)
	if err != nil {
		panic(err)
	}

	// Mark the file as saved
	// by creating a .save file in the same directory
	// as the original file
	if !noSaved {
		_ = os.WriteFile(filepath+".saved", []byte("win-reinstaller-helper: saved\npath: "+newFilePath), os.ModePerm)
	}

	// Save the new file path to the files.json file
	filesJson := LoadJson(path.Join(basePath, "files.json"))
	filesJson.Files = append(filesJson.Files, File{
		OriginalPath: filepath,
		NewPath:      newFilePath,
		Restore:      restore,
		IsDir:        false,
	})
	filesJson.Save()
}

func SaveDir(basePath, filepath string, restore, noSaved bool) {
	// Generate a random 8 character string
	str := RandStringRunes(8)

	newFilePath := path.Join(basePath, "files", str+".wrh")

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

	// Copy the file to the new location
	err := cp.Copy(filepath, newFilePath, opts)
	if err != nil {
		panic(err)
	}

	// Mark the file as saved
	// by creating a .save file in the same directory
	// as the original file
	if !noSaved {
		_ = os.WriteFile(strings.TrimSuffix(filepath, "/")+".saved", []byte("win-reinstaller-helper: saved\npath: "+newFilePath), os.ModePerm)
	}

	// Save the new file path to the files.json file
	filesJson := LoadJson(path.Join(basePath, "files.json"))
	filesJson.Files = append(filesJson.Files, File{
		OriginalPath: filepath,
		NewPath:      newFilePath,
		Restore:      restore,
		IsDir:        false,
	})
	filesJson.Save()
}

func Delete(basePath, filePath string) {
	filesJson := LoadJson(path.Join(basePath, "files.json"))
	var file File

	for i, ffile := range filesJson.Files {
		if ffile.OriginalPath == filePath {
			file = ffile
			filesJson.Names = append(filesJson.Names[:i], filesJson.Names[i+1:]...)
			break
		}
	}
	if (file == File{}) {
		for i, name := range filesJson.Names {
			if name == filePath {
				filesJson.Names = append(filesJson.Names[:i], filesJson.Names[i+1:]...)
				return
			}
		}
	}

	err := os.RemoveAll(file.NewPath)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(filePath + ".saved"); os.IsExist(err) {
		err := os.RemoveAll(strings.TrimSuffix(filePath, "/") + ".saved")
		if err != nil {
			panic(err)
		}
	}
}
