package project

import (
	"golang.org/x/net/context"

	"github.com/codeship/libcompose/project/events"
)

// Start starts the specified services (like docker start).
func (p *Project) Start(ctx context.Context, services ...string) error {
	eventWrapper := events.NewEventWrapper("Project Start", events.NewProjectStartStartEvent, events.NewProjectStartDoneEvent, events.NewProjectStartFailedEvent)
	return p.perform(eventWrapper, services, wrapperAction(func(wrapper *serviceWrapper, wrappers map[string]*serviceWrapper) {
		serviceEventWrapper := events.NewEventWrapper("Project Start", events.NewServiceStartStartEvent, events.NewServiceStartDoneEvent, events.NewServiceStartFailedEvent)
		wrapper.Do(wrappers, serviceEventWrapper, func(service Service) error {
			return service.Start(ctx)
		})
	}), nil)
}
