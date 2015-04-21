package commands

import (
	"errors"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/static"
)

func NewStaticDumpCommand(context *context.Context) *StaticDumpCommand {
	return &StaticDumpCommand{
		context: context,
	}
}

type StaticDumpCommand struct {
	core.Command
	context   *context.Context
	targetDir string
}

func (c *StaticDumpCommand) ID() string { return "dump" }
func (c *StaticDumpCommand) Description() string {
	return `Dumps static files to given directory
e.g.: $ ./patrol static:dump .`
}

func (c *StaticDumpCommand) Run() error {
	if static.IsEmbedded == false {
		return errors.New("patrol was not compiled with embedded static (build flag EMBED_STATIC).")
	}

	return nil
}

func (c *StaticDumpCommand) ParseArgs(args []string) (err error) {
	if len(args) != 1 {
		return errors.New("Please provide just one argument - target directory")
	}

	c.targetDir = args[0]

	return
}
