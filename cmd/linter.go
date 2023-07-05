package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/csunibo/synta"
)

func main() {
	recursive := flag.Bool("recurive", false, "Recursively check all files")

	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Printf("Usage: %s [file.synta] [folder]\n", os.Args[0])
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		fmt.Printf("Could not read synta definition file: %v", err)
		os.Exit(3)
	}

	syntaFile, err := synta.ParseSynta(string(data))
	if err != nil {
		fmt.Printf("Invalid synta definiton file: %v", err)
		os.Exit(4)
	}

	// TODO: use the filenamelinter library to actually check the folder
	fmt.Printf("%v %f\n", syntaFile, recursive)

	os.Exit(0)
}
