package project

import (
	"golang.org/x/net/context"

	"github.com/codeship/libcompose/project/events"
)

// Pull pulls the specified services (like docker pull).
func (p *Project) Pull(ctx context.Context, services ...string) error {
	return p.forEach(services, wrapperAction(func(wrapper *serviceWrapper, wrappers map[string]*serviceWrapper) {
		serviceEventWrapper := events.NewEventWrapper("Service Pull", events.NewServicePullStartEvent, events.NewServicePullDoneEvent, events.NewServicePullFailedEvent)
		wrapper.Do(nil, serviceEventWrapper, func(service Service) error {
			return service.Pull(ctx)
		})
	}), nil)
}
