package main

import (
	"github.com/kohkimakimoto/cli"
	"os"
	"fmt"
)

func main() {
	CLI := cli.NewCLI()
	CLI.Name = "example"
	CLI.Commands = []cli.Command{
		cli.Command{
			Name: "hello",
			Action: func(ctx *cli.Context) error {
				fmt.Println("Hello world")
				return nil
			},
		},
	}

	err := CLI.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Got a error: %s\n", err)
		os.Exit(1)
	}
}
