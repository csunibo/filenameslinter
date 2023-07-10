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
			return err
		}

		if file.IsDir() {
			if err := CheckName(synta, fs, file.Name()); err != nil {
				dirPath = path.Join(dirPath, file.Name())
				if err := CheckDir(synta, fs, dirPath); err != nil {
					return err
				}
			}
		} else if err = CheckName(synta, fs, file.Name()); err != nil {
			return err
		}
	}

	return
}

func CheckName(synta synta.Synta, fs fs.FS, path string) (err error) {
	file, err := fs.Open(path)

	if err != nil {
		err = fmt.Errorf("file %s does not exists: %v", path, err)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		err = fmt.Errorf("error while getting stats for file %s: %v", path, err)
		return
	}


    var reg *regexp.Regexp = nil
    if info.IsDir() {
        reg, err = syntaRegexp.ConvertWithoutExtension(synta)
        if err != nil {
            err = fmt.Errorf("can't convert synta to (dir) regexp: %v", err)
            return
        }
    } else {
        reg, err = syntaRegexp.Convert(synta)
        if err != nil {
            err = fmt.Errorf("can't convert synta to (file) regexp: %v", err)
            return
        }
    }

	if !reg.Match([]byte(info.Name())) {
		err = RegexMatchError{
			Regexp:   reg,
			Filename: info.Name(),
		}
	}

	return
}
