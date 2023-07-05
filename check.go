package filenameslinter

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/csunibo/synta"
	syntaRegexp "github.com/csunibo/synta/regexp"
)

func CheckDir(synta synta.Synta, dirPath string) (err error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return
	}

	for _, entry := range entries {
		file, err := entry.Info()
		if err != nil {
			return err
		}

		if file.IsDir() {
			if err := CheckFilePath(synta, file.Name()); err != nil {
				dirPath = path.Join(dirPath, file.Name())
				if err := CheckDir(synta, dirPath); err != nil {
					return err
				}
			}
		} else if err = CheckFilePath(synta, file.Name()); err != nil {
			return err
		}
	}

	return
}

func CheckFilePath(synta synta.Synta, path string) (err error) {
	file, err := os.Open(path)

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

	reg, err := syntaRegexp.Convert(synta)
	if err != nil {
		err = fmt.Errorf("can't convert synta to regexp: %v", err)
		return
	}

	if !reg.Match([]byte(info.Name())) {
		err = RegexMatchError{
			Regexp:   reg,
			Filename: info.Name(),
		}
	}

	return
}
