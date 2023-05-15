package deploy

import (
	"context"

	"github.com/aaraney/deployx/convert"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	apiclient "github.com/docker/docker/client"
)

func getStackFilter(namespace string) filters.Args {
	filter := filters.NewArgs()
	filter.Add("label", convert.LabelNamespace+"="+namespace)
	return filter
}

func getStackServices(ctx context.Context, apiclient apiclient.APIClient, namespace string) ([]swarm.Service, error) {
	return apiclient.ServiceList(ctx, types.ServiceListOptions{Filters: getStackFilter(namespace)})
}

func getStackNetworks(ctx context.Context, apiclient apiclient.APIClient, namespace string) ([]types.NetworkResource, error) {
	return apiclient.NetworkList(ctx, types.NetworkListOptions{Filters: getStackFilter(namespace)})
}

func getStackSecrets(ctx context.Context, apiclient apiclient.APIClient, namespace string) ([]swarm.Secret, error) {
	return apiclient.SecretList(ctx, types.SecretListOptions{Filters: getStackFilter(namespace)})
}

func getStackConfigs(ctx context.Context, apiclient apiclient.APIClient, namespace string) ([]swarm.Config, error) {
	return apiclient.ConfigList(ctx, types.ConfigListOptions{Filters: getStackFilter(namespace)})
}
