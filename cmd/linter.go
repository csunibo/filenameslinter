package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/csunibo/filenameslinter"
	"github.com/csunibo/synta"
)

func main() {
	recursive := flag.Bool("recursive", true, "Recursively check all files")
	ensureKebabCasing := flag.Bool("ensure-kebab-casing", true, "Check if directory names are in kebab-case")
	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Printf("Usage: %s [file.synta] [folder]\n", os.Args[0])
		os.Exit(1)
	}

	data, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		fmt.Printf("Could not read synta definition file: %v\n", err)
		os.Exit(3)
	}

	syntaFile, err := synta.ParseSynta(string(data))
	if err != nil {
		fmt.Printf("Invalid synta definiton file: %v\n", err)
		os.Exit(4)
	}

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not get current working directory: %v\n", err)
		os.Exit(5)
	}

	err = filenameslinter.CheckDir(syntaFile, os.DirFS(pwd), flag.Arg(1), *recursive, *ensureKebabCasing)
	if err != nil {
		extra := ""
		if *recursive {
			extra = "recursively "
		}
		fmt.Printf("Error while %schecking directory: %v\n", extra, err)
		os.Exit(6)
	}
	os.Exit(0)
}
