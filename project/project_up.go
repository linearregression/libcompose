package project

import (
	"golang.org/x/net/context"

	"github.com/codeship/libcompose/project/events"
	"github.com/codeship/libcompose/project/options"
)

// Up creates and starts the specified services (kinda like docker run).
// Up creates and starts the specified services (kinda like docker run).
func (p *Project) Up(ctx context.Context, options options.Up, services ...string) error {
	if err := p.initialize(ctx); err != nil {
		return err
	}
	eventWrapper := events.NewEventWrapper("Project Up", events.NewProjectUpStartEvent, events.NewProjectUpDoneEvent, events.NewProjectUpFailedEvent)
	return p.perform(eventWrapper, services, wrapperAction(func(wrapper *serviceWrapper, wrappers map[string]*serviceWrapper) {
		serviceEventWrapper := events.NewEventWrapper("Service Up", events.NewServiceUpStartEvent, events.NewServiceUpDoneEvent, events.NewServiceUpFailedEvent)
		wrapper.Do(wrappers, serviceEventWrapper, func(service Service) error {
			return service.Up(ctx, options)
		})
	}), func(service Service) error {
		return service.Create(ctx, options.Create)
	})
}
