package project

import (
	"golang.org/x/net/context"

	"github.com/codeship/libcompose/project/events"
)

// Kill kills the specified services (like docker kill).
func (p *Project) Kill(ctx context.Context, signal string, services ...string) error {
	eventWrapper := events.NewEventWrapper("Project Kill", events.NewProjectKillStartEvent, events.NewProjectKillDoneEvent, events.NewProjectKillFailedEvent)
	return p.perform(eventWrapper, services, wrapperAction(func(wrapper *serviceWrapper, wrappers map[string]*serviceWrapper) {
		serviceEventWrapper := events.NewEventWrapper("Service Pull", events.NewServiceKillStartEvent, events.NewServiceKillDoneEvent, events.NewServiceKillFailedEvent)
		wrapper.Do(nil, serviceEventWrapper, func(service Service) error {
			return service.Kill(ctx, signal)
		})
	}), nil)
}
