package commands

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/phonkee/patrol/context"

	"github.com/golang/glog"
	"github.com/mgutz/ansi"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/signals"
	"github.com/phonkee/patrol/utils"
)

/*
Commands
*/

const (
	HTTP_SERVER_COMMAND_HELP = `
Runs event process worker.
event:worker has following positional arguments
	./patrol event:http [listen=` + settings.HTTP_SERVER_DEFAULT_HOST + `]`
)

var (
	ErrPendingMigrations = errors.New("you have unapplied migrations, please run migrate.")

	// colors
	sectioncolor = ansi.ColorFunc("yellow+h:black")
	itemcolor    = ansi.ColorFunc("green+h:black")
)

/*
	List Commands command
		prints all available commands in patrol to stdout
*/
func NewCommonListCommandsCommand(c *context.Context, pr *core.PluginRegistry) *CommonListCommandsCommand {
	return &CommonListCommandsCommand{
		context: c,
		pr:      pr,
	}
}

type CommonListCommandsCommand struct {
	core.Command
	context *context.Context
	pr      *core.PluginRegistry
}

func (lcc *CommonListCommandsCommand) ID() string { return "list_commands" }
func (lcc *CommonListCommandsCommand) Description() string {
	return "Lists all available commands"
}
func (lcc *CommonListCommandsCommand) Run() error {
	fmt.Println("Available patrol commands:")
	fmt.Println()
	lcc.pr.Do(func(plugin core.Pluginer) error {
		commands := plugin.Commands()
		if len(commands) == 0 {
			return nil
		}
		fmt.Println(sectioncolor("[" + plugin.ID() + "]"))
		for _, command := range commands {
			line := "    " + utils.StringPadRight(command.ID(), " ", settings.LIST_COMMANDS_COMMAND_PADDING)
			lineLen := len(line)
			line = itemcolor(line)
			description := strings.Trim(strings.TrimSpace(command.Description()), "\n")

			if description != "" {
				description = strings.Replace(description, "\n", "\n"+strings.Repeat(" ", lineLen), -1)
				line = fmt.Sprintf("%s%s", line, description)
			}
			fmt.Println(line)
		}
		fmt.Println()
		return nil
	})
	return nil
}

/*
	Migrate command
		applies database migrations
*/
func NewCommonMigrateCommand(context *context.Context, pr *core.PluginRegistry) *CommonMigrateCommand {
	return &CommonMigrateCommand{
		context: context,
		pr:      pr,
	}
}

type CommonMigrateCommand struct {
	core.Command
	context *context.Context
	pr      *core.PluginRegistry
}

func (m *CommonMigrateCommand) ID() string          { return "migrate" }
func (m *CommonMigrateCommand) Description() string { return "Applies database migrations" }
func (m *CommonMigrateCommand) Run() error {
	se := core.NewSchemaEditor(m.context, m.pr)
	return se.Migrate()
}

/*HttpServer command
Runs http server
*/

func NewCommonHttpServerCommand(context *context.Context, pr *core.PluginRegistry) *CommonHttpServerCommand {
	return &CommonHttpServerCommand{
		context: context,
		pr:      pr,
	}
}

type CommonHttpServerCommand struct {
	core.Command
	context *context.Context
	pr      *core.PluginRegistry

	// cli args
	host   string
	static string
}

func (hsc *CommonHttpServerCommand) ID() string          { return "http" }
func (hsc *CommonHttpServerCommand) Description() string { return HTTP_SERVER_COMMAND_HELP }
func (hsc *CommonHttpServerCommand) Run() error {
	se := core.NewSchemaEditor(hsc.context, hsc.pr)
	if count, err := se.PendingMigrations(); err != nil {
		return err
	} else {
		if count > 0 {
			return ErrPendingMigrations
		}
	}

	glog.Infof("patrol: start listening on http://%s", hsc.host)

	/* send server start signal
	@TODO: run server in separate goroutine and fire signal after server has been started
	*/
	hsc.pr.Do(func(plugin core.Pluginer) error {
		if t, ok := plugin.(signals.OnHttpServerStartSignalHandler); ok {
			t.OnHttpServerStart()
		}
		return nil
	})

	// run http server
	if err := http.ListenAndServe(hsc.host, hsc.context.Router); err != nil {
		return err
	}

	return nil
}
func (hsc *CommonHttpServerCommand) ParseArgs(args []string) error {
	hsc.host = hsc.getHost(args)
	return nil
}

// Returns host from positional arguments
func (hsc *CommonHttpServerCommand) getHost(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return settings.HTTP_SERVER_DEFAULT_HOST
}

// Returns host from positional arguments
func (hsc *CommonHttpServerCommand) getStaticDir(args []string) string {
	if len(args) > 1 {
		return args[1]
	}
	return "./static/"
}

/*HttpServer command
Runs http server
*/

func NewCommonListRoutesCommand(context *context.Context, pr *core.PluginRegistry) *CommonListRoutesCommand {
	return &CommonListRoutesCommand{
		context: context,
		pr:      pr,
	}
}

type CommonListRoutesCommand struct {
	core.Command
	context *context.Context
	pr      *core.PluginRegistry
}

func (c *CommonListRoutesCommand) ID() string          { return "list_routes" }
func (c *CommonListRoutesCommand) Description() string { return "Lists all router routes" }
func (c *CommonListRoutesCommand) Run() error {

	fmt.Printf("List of registered routes:\n\n")
	c.pr.Do(func(plugin core.Pluginer) error {

		// no urls, no entry
		if len(plugin.URLs()) == 0 {
			return nil
		}

		fmt.Println(sectioncolor("[" + plugin.ID() + "]"))

		for _, uv := range plugin.URLs() {
			url := utils.StringPadRight(fmt.Sprintf("    %s", uv.URL()), " ", settings.LIST_ROUTES_COMMAND_PADDING)
			colorized := itemcolor(url)
			if len(url) > settings.LIST_ROUTES_COMMAND_PADDING {
				colorized = colorized + "\n" + strings.Repeat(" ", settings.LIST_ROUTES_COMMAND_PADDING)
			} else {
				colorized = utils.StringPadRight(colorized, " ", settings.LIST_ROUTES_COMMAND_PADDING)
			}
			fmt.Printf("%s %s:%s\n", colorized, plugin.ID(), uv.GetName())
		}
		fmt.Println()
		return nil
	})
	fmt.Println()
	return nil
}
