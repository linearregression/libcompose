package project

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/codeship/libcompose/project/events"
	"github.com/codeship/libcompose/project/options"
)

// Create creates the specified services (like docker create).
func (p *Project) Create(ctx context.Context, options options.Create, services ...string) error {
	if options.NoRecreate && options.ForceRecreate {
		return fmt.Errorf("no-recreate and force-recreate cannot be combined")
	}

	if err := p.initialize(ctx); err != nil {
		return err
	}

	eventWrapper := events.NewEventWrapper("Project Create", events.NewProjectCreateStartEvent, events.NewProjectCreateDoneEvent, events.NewProjectCreateFailedEvent)
	return p.perform(eventWrapper, services, wrapperAction(func(wrapper *serviceWrapper, wrappers map[string]*serviceWrapper) {
		serviceEventWrapper := events.NewEventWrapper("Service Create", events.NewServiceCreateStartEvent, events.NewServiceCreateDoneEvent, events.NewServiceCreateFailedEvent)
		wrapper.Do(wrappers, serviceEventWrapper, func(service Service) error {
			return service.Create(ctx, options)
		})
	}), nil)
}
