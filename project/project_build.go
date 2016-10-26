package project

import (
	"golang.org/x/net/context"

	"github.com/codeship/libcompose/project/events"
	"github.com/codeship/libcompose/project/options"
)

// Build builds the specified services (like docker build).
// Build builds the specified services (like docker build).
func (p *Project) Build(ctx context.Context, buildOptions options.Build, services ...string) error {
	eventWrapper := events.NewEventWrapper("Project Build", events.NewProjectBuildStartEvent, events.NewProjectBuildDoneEvent, events.NewProjectBuildFailedEvent)
	return p.perform(eventWrapper, services, wrapperAction(func(wrapper *serviceWrapper, wrappers map[string]*serviceWrapper) {
		serviceEventWrapper := events.NewEventWrapper("Service Build", events.NewServiceBuildStartEvent, events.NewServiceBuildDoneEvent, events.NewServiceBuildFailedEvent)
		wrapper.Do(wrappers, serviceEventWrapper, func(service Service) error {
			return service.Build(ctx, buildOptions)
		})
	}), nil)
}
