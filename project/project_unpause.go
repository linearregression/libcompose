package project

import (
	"golang.org/x/net/context"

	"github.com/codeship/libcompose/project/events"
)

// Unpause pauses the specified services containers (like docker pause).
func (p *Project) Unpause(ctx context.Context, services ...string) error {
	eventWrapper := events.NewEventWrapper("Project Pause", events.NewProjectUnpauseStartEvent, events.NewProjectUnpauseDoneEvent, events.NewProjectUnpauseFailedEvent)
	return p.perform(eventWrapper, services, wrapperAction(func(wrapper *serviceWrapper, wrappers map[string]*serviceWrapper) {
		serviceEventWrapper := events.NewEventWrapper("Service Pause", events.NewServiceUnpauseStartEvent, events.NewServiceUnpauseDoneEvent, events.NewServiceUnpauseFailedEvent)
		wrapper.Do(nil, serviceEventWrapper, func(service Service) error {
			return service.Unpause(ctx)
		})
	}), nil)
}
