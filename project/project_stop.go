package project

import (
	"golang.org/x/net/context"

	"github.com/codeship/libcompose/project/events"
)

// Stop stops the specified services (like docker stop).
func (p *Project) Stop(ctx context.Context, timeout int, services ...string) error {
	eventWrapper := events.NewEventWrapper("Project Stop", events.NewProjectStopStartEvent, events.NewProjectStopDoneEvent, events.NewProjectStopFailedEvent)
	return p.perform(eventWrapper, services, wrapperAction(func(wrapper *serviceWrapper, wrappers map[string]*serviceWrapper) {
		serviceEventWrapper := events.NewEventWrapper("Service Stop", events.NewServiceStopStartEvent, events.NewServiceStopDoneEvent, events.NewServiceStopFailedEvent)
		wrapper.Do(nil, serviceEventWrapper, func(service Service) error {
			return service.Stop(ctx, timeout)
		})
	}), nil)
}
