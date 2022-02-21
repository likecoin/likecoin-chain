package providers

import (
	AzureProvider "deploy/pkg/providers/azure"
	"errors"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const (
	GCP   string = "gcp"
	Azure string = "azure"
)

type Provider struct {
	Name string
}

func (provider Provider) GetVMBuilder() (func(deploymentId string, ctx *pulumi.Context, cfg *config.Config) (*remote.ConnectionArgs, []pulumi.Resource, error), error) {
	switch provider.Name {
	case Azure:
		return AzureProvider.CreateNodeVM, nil
	}

	return nil, errors.New("unknown provider")
}
