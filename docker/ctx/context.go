package ctx

import (
	"github.com/codeship/libcompose/docker/auth"
	"github.com/codeship/libcompose/docker/client"
	"github.com/codeship/libcompose/project"
	"github.com/docker/docker/cliconfig"
	"github.com/docker/docker/cliconfig/configfile"
)

// Context holds context meta information about a libcompose project and docker
// client information (like configuration file, builder to use, …)
type Context struct {
	project.Context
	ClientFactory client.Factory
	ConfigDir     string
	ConfigFile    *configfile.ConfigFile
	AuthLookup    auth.Lookup
}

// LookupConfig tries to load the docker configuration files, if any.
func (c *Context) LookupConfig() error {
	if c.ConfigFile != nil {
		return nil
	}

	config, err := cliconfig.Load(c.ConfigDir)
	if err != nil {
		return err
	}

	c.ConfigFile = config

	return nil
}
