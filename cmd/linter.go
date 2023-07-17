package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/csunibo/filenameslinter"
	"github.com/csunibo/synta"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
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
		log.Err(err).Msg("Could not read synta definition file")
		os.Exit(3)
	}

	syntaFile, err := synta.ParseSynta(string(data))
	if err != nil {
		log.Err(err).Msg("Invalid synta definiton file")
		os.Exit(4)
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Err(err).Msg("Could not get current working directory")
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
		log.Err(err).Msg("Error while " + extra + "checking directory")
		os.Exit(6)
	}
	os.Exit(0)
}
