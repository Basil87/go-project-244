package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "difference-calculator",
		Usage: "Compare and get diff of structs",

		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("Oh, Hi Mark!")
			return nil
		},
	}

	cli.RootCommandHelpTemplate = `
	./bin/gendiff --help
	NAME:
	   gendiff - Compares two configuration files and shows a difference.
	
	USAGE:
	   gendiff [global options]
	
	GLOBAL OPTIONS:
	   --format string, -f string  output format (default: "stylish")
	   --help, -h                  show help
`

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
