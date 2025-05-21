// Package command provides a modular command registry and execution system.
package command

// Command defines the interface for executable commands.
type Command interface {
	Execute(args []string) error
	Metadata() MetaData
}

// Metadata holds descriptive and configuration data for a command.
type MetaData struct {
	Name             string
	ShortDescription string
	LongDescription  string
	HasNodes         bool
	Flags            []Flag
}

// Flag represents a command-line flag with its properties.
type Flag struct {
	Name         string
	Shorthand    string
	Description  string
	DefaultValue string
	Required     bool
}

// Registry manages a collection of commands and their hierarchy.
type Registry struct {
	commands map[string]*CommandNode
}

// CommandNode represents a node in the command hierarchy.
type CommandNode struct {
	Command     Command
	Subcommands map[string]*CommandNode
}
