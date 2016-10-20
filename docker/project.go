package docker

import (
	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/codeship/libcompose/config"
	"github.com/codeship/libcompose/docker/auth"
	"github.com/codeship/libcompose/docker/client"
	"github.com/codeship/libcompose/docker/ctx"
	"github.com/codeship/libcompose/docker/network"
	"github.com/codeship/libcompose/docker/service"
	"github.com/codeship/libcompose/docker/volume"
	"github.com/codeship/libcompose/labels"
	"github.com/codeship/libcompose/project"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// NewProject creates a Project with the specified context.
func NewProject(context *ctx.Context, parseOptions *config.ParseOptions) (project.APIProject, error) {
	if context.AuthLookup == nil {
		context.AuthLookup = auth.NewConfigLookup(context.ConfigFile)
	}

	if context.ServiceFactory == nil {
		context.ServiceFactory = service.NewFactory(context)
	}

	if context.ClientFactory == nil {
		factory, err := client.NewDefaultFactory(client.Options{})
		if err != nil {
			return nil, err
		}
		context.ClientFactory = factory
	}

	if context.NetworksFactory == nil {
		networksFactory := &network.DockerFactory{
			ClientFactory: context.ClientFactory,
		}
		context.NetworksFactory = networksFactory
	}

	if context.VolumesFactory == nil {
		volumesFactory := &volume.DockerFactory{
			ClientFactory: context.ClientFactory,
		}
		context.VolumesFactory = volumesFactory
	}

	// FIXME(vdemeester) Remove the context duplication ?
	runtime := &Project{
		clientFactory: context.ClientFactory,
	}
	p := project.NewProject(&context.Context, runtime, parseOptions)

	err := p.Parse()
	if err != nil {
		return nil, err
	}

	if err = context.LookupConfig(); err != nil {
		logrus.Errorf("Failed to open project %s: %v", p.Name, err)
		return nil, err
	}

	return p, err
}

// Project implements project.RuntimeProject and define docker runtime specific methods.
type Project struct {
	clientFactory client.Factory
}

// RemoveOrphans implements project.RuntimeProject.RemoveOrphans.
// It will remove orphan containers that are part of the project but not to any services.
func (p *Project) RemoveOrphans(ctx context.Context, projectName string, serviceConfigs *config.ServiceConfigs) error {
	client := p.clientFactory.Create(nil)
	filter := filters.NewArgs()
	filter.Add("label", labels.PROJECT.EqString(projectName))
	containers, err := client.ContainerList(ctx, types.ContainerListOptions{
		Filter: filter,
	})
	if err != nil {
		return err
	}
	currentServices := map[string]struct{}{}
	for _, serviceName := range serviceConfigs.Keys() {
		currentServices[serviceName] = struct{}{}
	}
	for _, container := range containers {
		serviceLabel := container.Labels[labels.SERVICE.Str()]
		if _, ok := currentServices[serviceLabel]; !ok {
			if err := client.ContainerKill(ctx, container.ID, "SIGKILL"); err != nil {
				return err
			}
			if err := client.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{
				Force: true,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}
