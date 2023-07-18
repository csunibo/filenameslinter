package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/csunibo/synta"
	log "golang.org/x/exp/slog"

	"github.com/csunibo/filenameslinter"
)

func main() {

	recursive := flag.Bool("recursive", true, "Recursively check all files")
	ensureKebabCasing := flag.Bool("ensure-kebab-casing", true, "Check if directory names are in kebab-case")
	ignoreDotfiles := flag.Bool("ignore-dotfiles", true, "Ignore files and folders that start with a dot")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Printf("Usage: %s [file.synta] [<folder>]\n", os.Args[0])
		os.Exit(1)
	}

	data, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		log.Error("could not read synta definition file", "err", err)
		os.Exit(3)
	}

	syntaFile, err := synta.ParseSynta(string(data))
	if err != nil {
		log.Error("invalid synta definiton file", "err", err)
		os.Exit(4)
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Error("Could not get current working directory", "err", err)
		os.Exit(5)
	}

	dirPath := "."
	if len(flag.Args()) > 1 {
		dirPath = flag.Arg(1)
	}

	opts := filenameslinter.Options{
		Recursive:         *recursive,
		EnsureKebabCasing: *ensureKebabCasing,
		IgnoreDotfiles:    *ignoreDotfiles,
	}
	err = filenameslinter.CheckDir(syntaFile, os.DirFS(pwd), dirPath, &opts)
	if err != nil {
		extra := ""
		if *recursive {
			extra = "recursively "
		}
		log.Error("Error while "+extra+"checking directory", "err", err)
		os.Exit(6)
	}
	os.Exit(0)
}
