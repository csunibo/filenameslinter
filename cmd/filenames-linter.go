package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
)

func main() {
	fmt.Println("Hello World")
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
