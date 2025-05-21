package command

import (
	"errors"
	"fmt"
)

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]*CommandNode),
	}
}

// Register adds a command to the registry with optional subcommands.
func (r *Registry) New(path []string, cmd Command, setup func()) error {
	if cmd.Metadata().Name == "" {
		return errors.New("command name cannot be empty")
	}
	if setup != nil {
		setup()
	}

	node := &CommandNode{
		Command:     cmd,
		Subcommands: make(map[string]*CommandNode),
	}

	// If the path slice has no values (or the "HasNodes" field is set to false)
	// we will assume this is the parent command and set it as so
	if len(path) == 0 || !cmd.Metadata().HasNodes {
		if _, exists := r.findCommand([]string{cmd.Metadata().Name}); exists {
			return fmt.Errorf("command %q already exists", cmd.Metadata().Name)
		}
		r.commands[cmd.Metadata().Name] = node
		return nil
	}

	// Register as a sub-command
	current := r.commands

	// If the path has only one value, only the parent command came along.
	// They either forgot their children (subcommands), or its not their weekend yet.
	//
	// ["parent"]
	// Instead of
	// ["parent", "child", ...n]
	if len(path) == 1 {
		return errors.New("no subcommands provided")
	}
	for i, part := range path {
		if i == len(path)-1 {
			if _, exists := r.findCommand([]string{part}); !exists {
				return fmt.Errorf("parent command %q not found", part)
			}
			if _, exists := current[part].Subcommands[cmd.Metadata().Name]; exists {
				return fmt.Errorf("subcommand %q already exists under %q", cmd.Metadata().Name, part)
			}
			current[part].Subcommands[cmd.Metadata().Name] = node
			return nil
		}
		if _, exists := current[part]; !exists {
			return fmt.Errorf("parent command %q not found", part)
		}
		current = current[part].Subcommands
	}
	return nil
}

// findCommand looks for a command (or subcommand) and returns it, along with a confirmation of existance.
func (r *Registry) findCommand(path []string) (*CommandNode, bool) {
	current := r.commands
	for i, part := range path {
		node, exists := current[part]
		if !exists {
			return nil, false
		}
		if i == len(path)-1 {
			return node, true
		}
		current = node.Subcommands
	}
	return nil, false
}

// Execute triggers the command's handler.
func (r *Registry) Execute(path, args []string) error {
	if len(path) == 0 {
		return errors.New("command path cannot not be empty")
	}

	current := r.commands
	for i, part := range path {
		node, exists := current[part]
		if !exists {
			return fmt.Errorf("command %q not found", part)
		}
		if i == len(path)-1 {
			return node.Command.Execute(args)
		}
		current = node.Subcommands
	}

	return errors.New("command not found")
}
