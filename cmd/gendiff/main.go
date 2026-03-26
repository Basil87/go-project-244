package main

import (
	"code"
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

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "stylish",
				Usage:   "use human-readable format",
			},
		},

		Action: func(ctx context.Context, cmd *cli.Command) error {

			if cmd.Args().Len() == 0 {
				fmt.Println("path is required")
				return nil
			}

			file1 := cmd.Args().Get(0)
			file2 := cmd.Args().Get(1)
			result, err := code.GetDiff(file1, file2)
			if err != nil {
				return err
			}

			fmt.Println(result)
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
