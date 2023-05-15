package deploy

import (
	"context"
	"fmt"
	"sort"

	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/api/types/swarm"
)

func sortServiceByName(services []swarm.Service) func(i, j int) bool {
	return func(i, j int) bool {
		return services[i].Spec.Name < services[j].Spec.Name
	}
}

func removeServices(
	ctx context.Context,
	dockerCli command.Cli,
	services []swarm.Service,
) bool {
	var hasError bool
	sort.Slice(services, sortServiceByName(services))
	for _, service := range services {
		fmt.Fprintf(dockerCli.Out(), "Removing service %s\n", service.Spec.Name)
		if err := dockerCli.Client().ServiceRemove(ctx, service.ID); err != nil {
			hasError = true
			fmt.Fprintf(dockerCli.Err(), "Failed to remove service %s: %s", service.ID, err)
		}
	}
	return hasError
}
