package convert

import (
	"os"
	"strings"

	composego "github.com/compose-spec/compose-go/types"
	"github.com/docker/docker/api/types"
	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"
)

const (
	// LabelNamespace is the label used to track stack resources
	LabelNamespace = "com.docker.stack.namespace"
)

// Namespace mangles names by prepending the name
type Namespace struct {
	name string
}

// Scope prepends the namespace to a name
func (n Namespace) Scope(name string) string {
	return n.name + "_" + name
}

// Descope returns the name without the namespace prefix
func (n Namespace) Descope(name string) string {
	return strings.TrimPrefix(name, n.name+"_")
}

// Name returns the name of the namespace
func (n Namespace) Name() string {
	return n.name
}

// NewNamespace returns a new Namespace for scoping of names
func NewNamespace(name string) Namespace {
	return Namespace{name: name}
}

// AddStackLabel returns labels with the namespace label added
func AddStackLabel(namespace Namespace, labels map[string]string) map[string]string {
	if labels == nil {
		labels = make(map[string]string)
	}
	labels[LabelNamespace] = namespace.name
	return labels
}

// Networks from the compose-file type to the engine API type
func Networks(namespace Namespace, networks composego.Networks, servicesNetworks map[string]struct{}) (map[string]types.NetworkCreate, []string) {
	if networks == nil {
		networks = make(composego.Networks)
	}

	externalNetworks := []string{}
	result := make(map[string]types.NetworkCreate)
	for internalName := range servicesNetworks {
		network := networks[internalName]
		if network.External.External {
			externalNetworks = append(externalNetworks, network.Name)
			continue
		}

		createOpts := types.NetworkCreate{
			Labels:     AddStackLabel(namespace, network.Labels),
			Driver:     network.Driver,
			Options:    network.DriverOpts,
			Internal:   network.Internal,
			Attachable: network.Attachable,
		}

		if network.Ipam.Driver != "" || len(network.Ipam.Config) > 0 {
			createOpts.IPAM = &networktypes.IPAM{}
		}

		if network.Ipam.Driver != "" {
			createOpts.IPAM.Driver = network.Ipam.Driver
		}
		for _, ipamConfig := range network.Ipam.Config {
			config := networktypes.IPAMConfig{
				Subnet: ipamConfig.Subnet,
			}
			createOpts.IPAM.Config = append(createOpts.IPAM.Config, config)
		}

		networkName := namespace.Scope(internalName)
		if network.Name != "" {
			networkName = network.Name
		}
		result[networkName] = createOpts
	}

	return result, externalNetworks
}

// Secrets converts secrets from the Compose type to the engine API type
func Secrets(namespace Namespace, secrets composego.Secrets) ([]swarm.SecretSpec, error) {
	result := []swarm.SecretSpec{}
	for name, secret := range secrets {
		if secret.External.External {
			continue
		}

		var obj SwarmFileObject
		var err error
		if secret.Driver != "" {
			obj = driverObjectConfig(namespace, name, composego.FileObjectConfig(secret))
		} else {
			obj, err = fileObjectConfig(namespace, name, composego.FileObjectConfig(secret))
		}
		if err != nil {
			return nil, err
		}
		spec := swarm.SecretSpec{Annotations: obj.Annotations, Data: obj.Data}
		if secret.Driver != "" {
			spec.Driver = &swarm.Driver{
				Name:    secret.Driver,
				Options: secret.DriverOpts,
			}
		}
		if secret.TemplateDriver != "" {
			spec.Templating = &swarm.Driver{
				Name: secret.TemplateDriver,
			}
		}
		result = append(result, spec)
	}
	return result, nil
}

func Configs(namespace Namespace, configs map[string]composego.ConfigObjConfig) ([]swarm.ConfigSpec, error) {
	result := []swarm.ConfigSpec{}
	for name, config := range configs {
		if config.External.External {
			continue
		}

		obj, err := fileObjectConfig(namespace, name, composego.FileObjectConfig(config))
		if err != nil {
			return nil, err
		}
		spec := swarm.ConfigSpec{Annotations: obj.Annotations, Data: obj.Data}
		if config.TemplateDriver != "" {
			spec.Templating = &swarm.Driver{
				Name: config.TemplateDriver,
			}
		}
		result = append(result, spec)
	}
	return result, nil
}

type SwarmFileObject struct {
	Annotations swarm.Annotations
	Data        []byte
}

func driverObjectConfig(namespace Namespace, name string, obj composego.FileObjectConfig) SwarmFileObject {
	if obj.Name != "" {
		name = obj.Name
	} else {
		name = namespace.Scope(name)
	}

	return SwarmFileObject{
		Annotations: swarm.Annotations{
			Name:   name,
			Labels: AddStackLabel(namespace, obj.Labels),
		},
		Data: []byte{},
	}
}

func fileObjectConfig(namespace Namespace, name string, obj composego.FileObjectConfig) (SwarmFileObject, error) {
	data, err := os.ReadFile(obj.File)
	if err != nil {
		return SwarmFileObject{}, err
	}

	if obj.Name != "" {
		name = obj.Name
	} else {
		name = namespace.Scope(name)
	}

	return SwarmFileObject{
		Annotations: swarm.Annotations{
			Name:   name,
			Labels: AddStackLabel(namespace, obj.Labels),
		},
		Data: data,
	}, nil
}
