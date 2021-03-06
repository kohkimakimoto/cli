package cli

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Command is a subcommand for a cli.App.
type Command struct {
	// The name of the command
	Name string
	// short name of the command. Typically one character (deprecated, use `Aliases`)
	ShortName string
	// A list of aliases for the command
	Aliases []string
	// A short description of the usage of this command
	Usage string
	// A longer explanation of how the command works
	Description string
	// An action to execute before any sub-subcommands are run, but after the context is ready
	// If a non-nil error is returned, no sub-subcommands are run
	Before func(context *Context) error
	// An action to execute after any subcommands are run, but after the subcommand has finished
	// It is run even if Action() panics
	After func(context *Context) error
	// The function to call when this command is invoked
	Action func(context *Context) error
	// List of child commands
	Subcommands []Command
	// List of flags to parse
	Flags []Flag
	// Treat all flags as normal arguments if true
	SkipFlagParsing bool
	// Boolean to hide built-in help command
	HideHelp bool
	// custom help template
	HelpTemplate string
	// Meta data to store custom values for the command.
	Meta map[string]interface{}

	commandNamePath []string
}

// Returns the full name of the command.
// For subcommands this ensures that parent commands are part of the command path
func (c Command) FullName() string {
	if c.commandNamePath == nil {
		return c.Name
	}
	return strings.Join(c.commandNamePath, " ")
}

// Invokes the command given the context, parses ctx.Args() to generate command-specific flags
func (c Command) Run(ctx *Context) error {
	if len(c.Subcommands) > 0 || c.Before != nil || c.After != nil {
		return c.startApp(ctx)
	}

	if !c.HideHelp && (HelpFlag != BoolFlag{}) {
		// append help to flags
		c.Flags = append(
			c.Flags,
			HelpFlag,
		)
	}

	set := flagSet(c.Name, c.Flags)
	set.SetOutput(ioutil.Discard)

	var err error
	if !c.SkipFlagParsing {
		firstFlagIndex := -1
		terminatorIndex := -1
		for index, arg := range ctx.Args() {
			if arg == "--" {
				terminatorIndex = index
				break
			} else if strings.HasPrefix(arg, "-") && firstFlagIndex == -1 {
				firstFlagIndex = index
			}
		}

		if firstFlagIndex > -1 {
			args := ctx.Args()
			regularArgs := make([]string, len(args[1:firstFlagIndex]))
			copy(regularArgs, args[1:firstFlagIndex])

			var flagArgs []string
			if terminatorIndex > -1 {
				flagArgs = args[firstFlagIndex:terminatorIndex]
				regularArgs = append(regularArgs, args[terminatorIndex:]...)
			} else {
				flagArgs = args[firstFlagIndex:]
			}

			err = set.Parse(append(flagArgs, regularArgs...))
		} else {
			err = set.Parse(ctx.Args().Tail())
		}
	} else {
		if c.SkipFlagParsing {
			err = set.Parse(append([]string{"--"}, ctx.Args().Tail()...))
		}
	}

	if err != nil {
		fmt.Fprintf(ctx.App.Writer, "Incorrect Usage.\n%s\n", err)
		fmt.Fprintln(ctx.App.Writer)
		ShowCommandHelp(ctx, c.Name)
		return err
	}

	nerr := normalizeFlags(c.Flags, set)
	if nerr != nil {
		fmt.Fprintln(ctx.App.Writer, nerr)
		fmt.Fprintln(ctx.App.Writer)
		ShowCommandHelp(ctx, c.Name)
		return nerr
	}
	context := NewContext(ctx.App, set, ctx)

	if checkCommandHelp(context, c.Name) {
		return nil
	}
	context.Command = c
	return c.Action(context)
}

func (c Command) Names() []string {
	names := []string{c.Name}

	if c.ShortName != "" {
		names = append(names, c.ShortName)
	}

	return append(names, c.Aliases...)
}

// Returns true if Command.Name or Command.ShortName matches given name
func (c Command) HasName(name string) bool {
	for _, n := range c.Names() {
		if n == name {
			return true
		}
	}
	return false
}

func (c Command) startApp(ctx *Context) error {
	app := NewApp()

	// set the name and usage
	app.Name = fmt.Sprintf("%s %s", ctx.App.Name, c.Name)
	if c.Description != "" {
		app.Usage = c.Description
	} else {
		app.Usage = c.Usage
	}

	// set CommandNotFound
	app.CommandNotFound = ctx.App.CommandNotFound

	// set the flags and commands
	app.Commands = c.Subcommands
	app.Flags = c.Flags
	app.HideHelp = c.HideHelp

	app.Version = ctx.App.Version
	app.HideVersion = ctx.App.HideVersion
	app.Compiled = ctx.App.Compiled
	app.Writer = ctx.App.Writer

	// set the actions
	app.Before = c.Before
	app.After = c.After
	if c.Action != nil {
		app.Action = c.Action
	} else {
		app.Action = HelpSubcommand.Action
	}

	var newCmds []Command
	for _, cc := range app.Commands {
		cc.commandNamePath = []string{c.Name, cc.Name}
		newCmds = append(newCmds, cc)
	}
	app.Commands = newCmds

	return app.RunAsSubcommand(ctx)
}
