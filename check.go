package filenameslinter

import (
	"fmt"
	"path"
    "io/fs"
    "regexp"

	"github.com/csunibo/synta"
	syntaRegexp "github.com/csunibo/synta/regexp"
)

func CheckDir(synta synta.Synta, fs fs.ReadDirFS, dirPath string) (err error) {
    entries, err := fs.ReadDir(dirPath)
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
			if err := CheckName(synta, file.Name(), true); err != nil {
				dirPath = path.Join(dirPath, file.Name())
				if err := CheckDir(synta, fs, dirPath); err != nil {
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
