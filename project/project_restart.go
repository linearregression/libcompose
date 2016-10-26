package project

import (
	"golang.org/x/net/context"

	"github.com/codeship/libcompose/project/events"
)

// Restart restarts the specified services (like docker restart).
func (p *Project) Restart(ctx context.Context, timeout int, services ...string) error {
	eventWrapper := events.NewEventWrapper("Project Restart", events.NewProjectRestartStartEvent, events.NewProjectRestartDoneEvent, events.NewProjectRestartFailedEvent)
	return p.perform(eventWrapper, services, wrapperAction(func(wrapper *serviceWrapper, wrappers map[string]*serviceWrapper) {
		serviceEventWrapper := events.NewEventWrapper("Service Restart", events.NewServiceRestartStartEvent, events.NewServiceRestartDoneEvent, events.NewServiceRestartFailedEvent)
		wrapper.Do(wrappers, serviceEventWrapper, func(service Service) error {
			return service.Restart(ctx, timeout)
		})
	}), nil)
}
