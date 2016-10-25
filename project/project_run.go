package project

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/codeship/libcompose/project/events"
	"github.com/codeship/libcompose/project/options"
)

// Run executes a one off command (like `docker run image command`).
func (p *Project) Run(ctx context.Context, serviceName string, commandParts []string, opts options.Run) (int, error) {
	if !p.ServiceConfigs.Has(serviceName) {
		return 1, fmt.Errorf("%s is not defined in the template", serviceName)
	}

	if err := p.initialize(ctx); err != nil {
		return 1, err
	}
	var exitCode int
	err := p.forEach([]string{}, wrapperAction(func(wrapper *serviceWrapper, wrappers map[string]*serviceWrapper) {
		serviceEventWrapper := events.NewEventWrapper("Service Run", events.NewServiceRunStartEvent, events.NewServiceRunDoneEvent, events.NewServiceRunFailedEvent)
		wrapper.Do(wrappers, serviceEventWrapper, func(service Service) error {
			if service.Name() == serviceName {
				code, err := service.Run(ctx, commandParts, opts)
				exitCode = code
				return err
			}
			return nil
		})
	}), func(service Service) error {
		return service.Create(ctx, options.Create{})
	})
	return exitCode, err
}
