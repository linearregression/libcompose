package project

import (
	"golang.org/x/net/context"

	"github.com/codeship/libcompose/project/events"
)

// Log aggregates and prints out the logs for the specified services.
func (p *Project) Log(ctx context.Context, follow bool, services ...string) error {
	return p.forEach(services, wrapperAction(func(wrapper *serviceWrapper, wrappers map[string]*serviceWrapper) {
		wrapper.Do(nil, events.NewDummyEventWrapper("Log"), func(service Service) error {
			return service.Log(ctx, follow)
		})
	}), nil)
}
