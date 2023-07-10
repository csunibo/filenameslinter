package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/csunibo/synta"
	"github.com/csunibo/filenameslinter"
)

func main() {
	recursive := flag.Bool("recurive", false, "Recursively check all files")

	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Printf("Usage: %s [file.synta] [folder]\n", os.Args[0])
		os.Exit(1)
	}

	data, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		fmt.Printf("Could not read synta definition file: %v", err)
		os.Exit(3)
	}

	syntaFile, err := synta.ParseSynta(string(data))
	if err != nil {
		fmt.Printf("Invalid synta definiton file: %v", err)
		os.Exit(4)
	}

    pwd, err := os.Getwd()
	if err != nil {
        fmt.Printf("Could not get current working directory: %v", err)
		os.Exit(5)
	}

	// TODO: use the filenamelinter library to actually check the folder
   filenameslinter.CheckDir(syntaFile, os.DirFS(pwd), flag.Arg(1))
	os.Exit(0)
}
