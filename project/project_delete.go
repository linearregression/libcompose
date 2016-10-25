package project

import (
	"golang.org/x/net/context"

	"github.com/codeship/libcompose/project/events"
	"github.com/codeship/libcompose/project/options"
)

// Delete removes the specified services (like docker rm).
// Delete removes the specified services (like docker rm).
func (p *Project) Delete(ctx context.Context, options options.Delete, services ...string) error {
	eventWrapper := events.NewEventWrapper("Project Delete", events.NewProjectDeleteStartEvent, events.NewProjectDeleteDoneEvent, events.NewProjectDeleteFailedEvent)
	return p.perform(eventWrapper, services, wrapperAction(func(wrapper *serviceWrapper, wrappers map[string]*serviceWrapper) {
		serviceEventWrapper := events.NewEventWrapper("Service Delete", events.NewServiceDeleteStartEvent, events.NewServiceDeleteDoneEvent, events.NewServiceDeleteFailedEvent)
		wrapper.Do(nil, serviceEventWrapper, func(service Service) error {
			return service.Delete(ctx, options)
		})
	}), nil)
}
