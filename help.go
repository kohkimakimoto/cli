package cli

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"text/template"
)

// The text template for the Default help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var AppHelpTemplate = `{{.Name}} {{.Version}}{{if .Usage}}

{{.Usage}}{{end}}

Usage:
  {{.Name}}{{if .Flags}} [<options>]{{end}} <command> [<arguments>]
{{if .Flags}}
Options:
  {{range .Flags}}{{.}}
  {{end}}{{end}}
Commands:
  {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
  {{end}}
`

// The text template for the command help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var CommandHelpTemplate = `{{.Name}}{{if .Usage}}

{{.Usage}}{{end}}

Usage:
  {{.Name}}{{if .Flags}} [<options>]{{end}} [<arguments>]

{{if .Flags}}Options:
  {{range .Flags}}{{.}}
  {{end}}{{end}}{{if .Description}}

Description:
  {{.Description}}
{{end}}
`

// The text template for the subcommand help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var SubcommandHelpTemplate = `{{.Name}}{{if .Usage}}

{{.Usage}}{{end}}

Usage:
  {{.Name}} <command>{{if .Flags}} [<options>]{{end}} [<arguments>]
{{if .Flags}}
Options:
  {{range .Flags}}{{.}}
  {{end}}{{end}}
Commands:
  {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
  {{end}}
`

var HelpCommand = Command{
	Name:    "help",
	Aliases: []string{"h"},
	Usage:   "Shows a list of commands or help for one command",
	Action: func(c *Context) error {
		args := c.Args()
		if args.Present() {
			ShowCommandHelp(c, args.First())
		} else {
			ShowAppHelp(c)
		}
		return nil
	},
}

var HelpSubcommand = Command{
	Name:    "help",
	Aliases: []string{"h"},
	Usage:   "Shows a list of commands or help for one command",
	Action: func(c *Context) error {
		args := c.Args()
		if args.Present() {
			ShowCommandHelp(c, args.First())
		} else {
			ShowSubcommandHelp(c)
		}
		return nil
	},
}

// Prints help for the App or Command
type helpPrinter func(w io.Writer, templ string, data interface{})

var HelpPrinter helpPrinter = printHelp

// Prints version for the App
var VersionPrinter = printVersion

func ShowAppHelp(c *Context) {
	DefaultAppHelp(c)
}

var DefaultAppHelp = func(c *Context) {
	HelpPrinter(c.App.Writer, AppHelpTemplate, c.App)
}

// Prints the list of subcommands as the default app completion method
func DefaultAppComplete(c *Context) {
	for _, command := range c.App.Commands {
		for _, name := range command.Names() {
			fmt.Fprintln(c.App.Writer, name)
		}
	}
}

// Prints help for the given command
func ShowCommandHelp(ctx *Context, command string) {
	// show the subcommand help for a command with subcommands
	if command == "" {
		HelpPrinter(ctx.App.Writer, SubcommandHelpTemplate, ctx.App)
		return
	}

	for _, c := range ctx.App.Commands {
		if c.HasName(command) {
			if c.HelpTemplate != "" {
				HelpPrinter(ctx.App.Writer, c.HelpTemplate, c)
			} else {
				HelpPrinter(ctx.App.Writer, CommandHelpTemplate, c)
			}
			return
		}
	}

	if ctx.App.CommandNotFound != nil {
		ctx.App.CommandNotFound(ctx, command)
	} else {
		fmt.Fprintf(ctx.App.Writer, "Command '%v' is not defined.\n", command)
	}
}

// Prints help for the given subcommand
func ShowSubcommandHelp(c *Context) {
	ShowCommandHelp(c, c.Command.Name)
}

// Prints the version number of the App
func ShowVersion(c *Context) {
	VersionPrinter(c)
}

func printVersion(c *Context) {
	fmt.Fprintf(c.App.Writer, "%v version %v\n", c.App.Name, c.App.Version)
}

func printHelp(out io.Writer, templ string, data interface{}) {
	funcMap := template.FuncMap{
		"join": strings.Join,
	}

	w := tabwriter.NewWriter(out, 0, 8, 1, '\t', 0)
	t := template.Must(template.New("help").Funcs(funcMap).Parse(templ))
	err := t.Execute(w, data)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func checkVersion(c *Context) bool {
	found := false
	if VersionFlag.Name != "" {
		eachName(VersionFlag.Name, func(name string) {
			if c.GlobalBool(name) || c.Bool(name) {
				found = true
			}
		})
	}
	return found
}

func checkHelp(c *Context) bool {
	found := false
	if HelpFlag.Name != "" {
		eachName(HelpFlag.Name, func(name string) {
			if c.GlobalBool(name) || c.Bool(name) {
				found = true
			}
		})
	}
	return found
}

func checkCommandHelp(c *Context, name string) bool {
	if c.Bool("h") || c.Bool("help") {
		ShowCommandHelp(c, name)
		return true
	}

	return false
}

func checkSubcommandHelp(c *Context) bool {
	if c.GlobalBool("h") || c.GlobalBool("help") {
		ShowSubcommandHelp(c)
		return true
	}

	return false
}
