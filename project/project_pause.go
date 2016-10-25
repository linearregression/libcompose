package project

import (
	"golang.org/x/net/context"

	"github.com/codeship/libcompose/project/events"
)

// Pause pauses the specified services containers (like docker pause).
func (p *Project) Pause(ctx context.Context, services ...string) error {
	eventWrapper := events.NewEventWrapper("Project Pause", events.NewProjectPauseStartEvent, events.NewProjectPauseDoneEvent, events.NewProjectPauseFailedEvent)
	return p.perform(eventWrapper, services, wrapperAction(func(wrapper *serviceWrapper, wrappers map[string]*serviceWrapper) {
		serviceEventWrapper := events.NewEventWrapper("Service Pause", events.NewServicePauseStartEvent, events.NewServicePauseDoneEvent, events.NewServicePauseFailedEvent)
		wrapper.Do(nil, serviceEventWrapper, func(service Service) error {
			return service.Pause(ctx)
		})
	}), nil)
}
