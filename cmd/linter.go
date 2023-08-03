package main

import (
	"flag"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/csunibo/synta"
	log "golang.org/x/exp/slog"

	"github.com/csunibo/filenameslinter"
)

func main() {
	recursive := flag.Bool("recursive", true, "Recursively check all files")
	ensureKebabCasing := flag.Bool("ensure-kebab-casing", true, "Check if directory names are in kebab-case")
	ignoreDotfiles := flag.Bool("ignore-dotfiles", true, "Ignore files and folders that start with a dot")
	syntaDefinition := flag.String("definition", "", "Synta definition file to check filenames against")
	flag.Parse()

	pwd, err := os.Getwd()
	if err != nil {
		log.Error("could not get current working directory", "err", err)
		os.Exit(1)
	}

	dirPath := "."
	parent := pwd
	if len(flag.Args()) > 0 {
		absDir := path.Join(pwd, flag.Arg(0))
		parent = filepath.Dir(strings.TrimSuffix(absDir, string(os.PathSeparator)))
		dirPath, err = filepath.Rel(parent, absDir)

		if err != nil {
			log.Error("could not make the path relative", "err", err)
			os.Exit(2)
		}
	}

	var syntaFile *synta.Synta = nil
	if *syntaDefinition != "" {
		data, err := os.ReadFile(*syntaDefinition)
		if err != nil {
			log.Error("could not read synta definition file", "err", err)
			os.Exit(3)
		}
		s, err := synta.ParseSynta(string(data))
		if err != nil {
			log.Error("invalid synta definiton file", "err", err)
			os.Exit(4)
		}
		syntaFile = &s
	}

	opts := filenameslinter.Options{
		Recursive:         *recursive,
		EnsureKebabCasing: *ensureKebabCasing,
		IgnoreDotfiles:    *ignoreDotfiles,
	}
	err = filenameslinter.CheckDir(syntaFile, os.DirFS(parent), dirPath, &opts)
	if err != nil {
		log.Error("error while checking directory", "recursive", *recursive, "err", err)
		os.Exit(5)
	}
	os.Exit(0)
}
