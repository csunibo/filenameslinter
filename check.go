package filenameslinter

import (
	"errors"
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/csunibo/synta"
	syntaRegexp "github.com/csunibo/synta/regexp"
	"github.com/rs/zerolog/log"
)

type Options struct {
	Recursive         bool
	EnsureKebabCasing bool
	IgnoreDotfiles    bool
}

var kebabRegexp = regexp.MustCompile("^[a-z0-9]+(-[a-z0-9]+)*(\\.[a-z0-9]+)?$")

// readDir uses the `readDir` method if the filesystem implements
// `fs.ReadDirFS`, otherwise opens the path and parses it using
// the `ReadDirFile` interface.
func readDir(fsys fs.FS, name string) ([]fs.DirEntry, error) {
	if fsys, ok := fsys.(fs.ReadDirFS); ok {
		return fsys.ReadDir(name)
	}

	file, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}

	dir, ok := file.(fs.ReadDirFile)
	if !ok {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: errors.New("not implemented")}
	}

	list, err := dir.ReadDir(-1)
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })

	err = file.Close()
	return list, err
}

func CheckDir(synta synta.Synta, fs fs.FS, dirPath string, opts *Options) (err error) {
	log.Info().Str("path", dirPath).Msg("Checking dir")

	entries, err := readDir(fs, dirPath)
	if err != nil {
		return
	}

	for _, entry := range entries {
		file, err := entry.Info()
		if err != nil {
			err = fmt.Errorf("Could not read directory: %v", err)
			return err
		}
		absPath := path.Join(dirPath, file.Name())
		shouldIgonre := opts.IgnoreDotfiles && strings.HasPrefix(file.Name(), ".")
		log.Info().Str("path", absPath).Bool("ignored", shouldIgonre).Msg("Checking path basename")

		// 1. (optionally) ignore files that start with a '.'
		if shouldIgonre {
			continue
		}

		// 1. (optionally) force all files to use kebab casing
		if opts.EnsureKebabCasing && !kebabRegexp.Match([]byte(file.Name())) {
			err = RegexMatchError{
				Regexp:   kebabRegexp.String(),
				Filename: file.Name(),
			}
			return err
		}

		// 3. Check the filename; if it's a directory and the name doesn't match,
		// recursively check it.
		err = CheckName(synta, file.Name(), file.IsDir())
		if err != nil && file.IsDir() && opts.Recursive {
			err = CheckDir(synta, fs, absPath, opts)
		}

		if err != nil {
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
			Regexp:   reg.String(),
			Filename: name,
		}
	}

	return
}
