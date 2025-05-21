package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/wriverarincon/command"
)

type simpleCommand struct{}

var ran bool
var registry *command.Registry
var meta command.MetaData

func new() *simpleCommand {
	ran = false
	return &simpleCommand{}
}

func (c *simpleCommand) Execute(args []string) error {
	ran = true
	return nil
}

func (c *simpleCommand) Metadata() command.MetaData {
	return meta
}

func TestMain(m *testing.M) {
	fmt.Println("Setting up command tests...")

	registry = command.NewRegistry()

	exitCode := m.Run()

	fmt.Println("Tearing down command tests...")

	os.Exit(exitCode)
}

func TestCommand(t *testing.T) {
	m := command.MetaData{
		Name:             "test",
		ShortDescription: "this is a test command",
		LongDescription:  "this is a test command for the command library",
		Flags:            nil,
	}
	meta = m
	err := registry.New(nil, new(), nil)
	registry.Execute([]string{"test"})

	if !ran {
		t.Fatalf("command didn't run\ngot error: %v\nregistry: %v", err, registry)
	}
}

func TestSubCommand(t *testing.T) {
	m := command.MetaData{
		Name:             "testing",
		ShortDescription: "this is a test subcommand",
		LongDescription:  "this is a test subcommand for the command library",
		Flags:            nil,
	}
	meta = m
	err := registry.New([]string{"test", "testing"}, new(), nil)
	registry.Execute([]string{"test", "testing"})

	if !ran {
		t.Fatalf("subcommand didn't run\ngot: %v\nregistry: %v", err, registry)
	}
}
