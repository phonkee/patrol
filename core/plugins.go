/* Plugin system

This is implementation of simple plugin system for patrol. We currently support
only processing of incoming messages/requests and custom http handlers.
There is no support for frontend but it's planned in near future.
*/

package core

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"github.com/phonkee/patrol/utils"
)

var (
	ErrPluginNotFound            = errors.New("plugin not found")
	ErrPluginAlreadyRegistered   = errors.New("plugin already registered")
	ErrCommandNotFound           = errors.New("command not found.")
	ErrMultipleCommandsFound     = errors.New("multiple commands found.")
	ErrPluginsAlreadyInitialized = errors.New("plugins already initialized")
)

// Pluginer interface defines plugin functionality
type Pluginer interface {

	// returns id of plugin
	// must be unique across all plugins
	// so maybe <author>-<plugin> should be used
	// e.g. phonkee-analytics
	ID() string

	// initializes plugin with application
	Init() error

	// returns list of available commands
	Commands() []Commander

	// returns registered views
	URLViews() []*URLView

	// returns migrations
	Migrations() []Migrationer
}

/* Base plugin implementation
it's recommended to use it in all plugins, so in case of new functionality
blank one will be already implemented in base plugin and there will be no need
to write methods that will satisfy Pluginer interface.

``go
	type AwesomePlugin struct {
		patrol.Plugin
	}
``
*/
type Plugin struct{}

// Returns list of commands that plugin exposes to patrol cli
func (p *Plugin) Commands() []Commander     { return []Commander{} }
func (p *Plugin) Init() error               { return nil }
func (p *Plugin) URLViews() []*URLView      { return []*URLView{} }
func (p *Plugin) Migrations() []Migrationer { return []Migrationer{} }

/*
	Plugin registry
*/
type PluginRegistry struct {
	plugins     []Pluginer
	initialized bool
}

// helper function to return plugin registry
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{}
}

// Registers plugin to registry
// 	basic checks are performed to be sure that no plugin is registered multiple
// 	times
func (pr *PluginRegistry) RegisterPlugin(plugin Pluginer) (err error) {
	// no possibility to register plugin after plugins have been initialized
	if pr.initialized {
		return ErrPluginsAlreadyInitialized
	}
	glog.V(2).Infof("patrol: register plugin %T as %s.", plugin, plugin.ID())

	// validate plugin
	if err = pr.validatePlugin(plugin); err != nil {
		return err
	}

	pr.plugins = append(pr.plugins, plugin)
	return nil
}

// performs function on all plugins
func (pr *PluginRegistry) Do(f func(Pluginer) error) error {
	for _, p := range pr.plugins {
		e := f(p)
		if e != nil {
			return e
		}
	}
	return nil
}

// performs function on all plugins in reverse order
func (pr *PluginRegistry) DoReverse(f func(Pluginer) error) error {
	for i := len(pr.plugins) - 1; i >= 0; i-- {
		e := f(pr.plugins[i])
		if e != nil {
			return e
		}
	}
	return nil
}

/*
	Returns plugin by plugin id.
*/
func (pr *PluginRegistry) Plugin(pluginid string) (Pluginer, error) {
	for _, plugin := range pr.plugins {
		if plugin.ID() == pluginid {
			return plugin, nil
		}
	}
	return nil, ErrPluginNotFound
}

/*
	returns command from identifier
	identifier is ":" delimited string in form <pluginid>:<command>
	if no pluginid is given all plugins are searched. If multiple commands with
	name are found, error is returned
*/
func (pr *PluginRegistry) Command(identifier string) (command Commander, err error) {
	var pid, cid string
	if pid, cid, err = utils.SplitIdentifier(identifier, ""); err != nil {
		return nil, err
	}

	foundCommands := make([]Commander, 0)

	pr.Do(func(plugin Pluginer) error {
		if pid != "" {
			if plugin.ID() != pid {
				return nil
			}
		}

		for _, c := range plugin.Commands() {
			if c.ID() == cid {
				foundCommands = append(foundCommands, c)
				return nil
			}
		}
		return nil
	})

	// command not found
	if len(foundCommands) == 0 {
		err = ErrCommandNotFound
	} else if len(foundCommands) > 1 {
		err = ErrMultipleCommandsFound
	} else if len(foundCommands) == 1 {
		command = foundCommands[0]
	}

	return
}

// Initializes plugins
func (p *PluginRegistry) InitPlugins() (err error) {
	// can initialize plugins only once
	if p.initialized {
		return ErrPluginsAlreadyInitialized
	}
	err = p.Do(func(plugin Pluginer) error {
		return plugin.Init()
	})
	p.initialized = true
	return
}

/* Validates plugin. following steps
1./ Check whether has plugin already been registered
2./ Check duplicate migrations
3./ Check duplicate commands
*/
func (p *PluginRegistry) validatePlugin(plugin Pluginer) (err error) {
	for _, p := range p.plugins {
		if p.ID() == plugin.ID() {
			return ErrPluginAlreadyRegistered
		}
	}

	mmap := map[string]bool{}
	for _, migration := range plugin.Migrations() {
		if _, ok := mmap[migration.ID()]; ok {
			return fmt.Errorf("duplicate migration %s:%s found.", plugin.ID(), migration.ID())
		}
		mmap[migration.ID()] = true
	}

	cmap := map[string]bool{}
	for _, command := range plugin.Commands() {
		if _, ok := cmap[command.ID()]; ok {
			return fmt.Errorf("duplicate command %s:%s found.", plugin.ID(), command.ID())
		}
		cmap[command.ID()] = true
	}

	return nil
}
