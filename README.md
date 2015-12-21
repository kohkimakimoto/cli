# CLI

[![Build Status](https://travis-ci.org/kohkimakimoto/cli.svg)](https://travis-ci.org/kohkimakimoto/cli)

This is a small package for CLI app in Go.
It is forked from [codegangsta/cli](https://github.com/codegangsta/cli), and modified code for my use case.

## Usage

```Go
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
		fmt.Fprintf(os.Stderr, "Got a error: %\n", err)
		os.Exit(1)
	}
}
```

## Forked from [codegangsta/cli](https://github.com/codegangsta/cli)

LICENSE

```
Copyright (C) 2013 Jeremy Saenz
All Rights Reserved.

MIT LICENSE

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```
