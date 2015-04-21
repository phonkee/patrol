package core

// Command interface defines cli commands
type Commander interface {

	// returns id of command which will be used
	// $ patrol <command_id>
	ID() string

	// Short description what command does for <list_commands> command
	Description() string

	// Runs command with application instance (config, database, queue, etc..)
	Run() error

	// parses args
	ParseArgs(args []string) error
}

type Command struct{}

func (c Command) Description() string      { return "" }
func (c Command) ParseArgs([]string) error { return nil }
