package filenameslinter

import (
	"errors"
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"sort"

	"github.com/csunibo/synta"
	syntaRegexp "github.com/csunibo/synta/regexp"
)

func CustomReadDir(fsys fs.FS, name string) ([]fs.DirEntry, error) {
	if fsys, ok := fsys.(fs.ReadDirFS); ok {
		return fsys.ReadDir(name)
	}

	file, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dir, ok := file.(fs.ReadDirFile)
	if !ok {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: errors.New("not implemented")}
	}

	list, err := dir.ReadDir(-1)
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, err
}

func CheckDir(synta synta.Synta, fs fs.FS, dirPath string, recursive bool, kebab_dir bool) (err error) {
	entries, err := CustomReadDir(fs, dirPath)
	if err != nil {
		return
	}

	for _, entry := range entries {
		file, err := entry.Info()
		if err != nil {
			err = fmt.Errorf("Could not read directory: %v", err)
			return err
		}

		if file.IsDir() {
			kebab_regexp := regexp.MustCompile("^[a-z0-9]+(-[a-z0-9]+)*$")
			if kebab_dir && !kebab_regexp.Match([]byte(file.Name())) {
				err = fmt.Errorf("Directories need to be in kebab-case")
				return err
			}

			if recursive {
				if err := CheckName(synta, file.Name(), true); err != nil {
					dirPath = path.Join(dirPath, file.Name())
					if err := CheckDir(synta, fs, dirPath, recursive, kebab_dir); err != nil {
						return err
					}
				}
			} else {
				if err := CheckName(synta, file.Name(), true); err != nil {
					return err
				}
			}
		} else if err = CheckName(synta, file.Name(), false); err != nil {
			return err
		}
	}

	return
}

func CheckName(synta synta.Synta, name string, isDir bool) (err error) {
	var reg *regexp.Regexp = nil
	if isDir {
		reg, err = syntaRegexp.ConvertWithoutExtension(synta)
		if err != nil {
			err = fmt.Errorf("Could not convert synta to (dir) regexp: %v", err)
			return
		}
	} else {
		reg, err = syntaRegexp.Convert(synta)
		if err != nil {
			err = fmt.Errorf("Could not convert synta to (file) regexp: %v", err)
			return
		}
	}

	if !reg.Match([]byte(name)) {
		err = RegexMatchError{
			Regexp:   reg,
			Filename: name,
		}
	}

	return
}
