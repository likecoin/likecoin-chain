package providers

import (
	"strings"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	serviceaccount "github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceAccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func CreateNodeVM(_deploymentId string, ctx *pulumi.Context, cfg *config.Config) (*remote.ConnectionArgs, []pulumi.Resource, error) {
	projectId := cfg.Require("project-id")
	vmUsername := cfg.Require("vm-username")
	publicKey := cfg.Require("vm-public-key")
	privateKey := cfg.RequireSecret("vm-private-key")
	sshAllowList := cfg.Require("vm-ssh-allow-list")

	vmHardware := cfg.Require("vm-hardware")
	vmDiskSize := cfg.GetInt("vm-disk-size")
	networkRegion := cfg.Require("network-region")
	vmZone := cfg.Require("vm-zone")

	deploymentId := strings.ReplaceAll(_deploymentId, "_", "-") // GCP naming convention

	var sshAllowCIDRs pulumi.StringArray
	for _, cidr := range strings.Split(sshAllowList, ",") {
		_cidr := pulumi.String(strings.TrimSpace(cidr))
		if len(_cidr) > 0 {
			sshAllowCIDRs = append(sshAllowCIDRs, _cidr)
		}
	}
	if len(sshAllowCIDRs) == 0 {
		sshAllowCIDRs = pulumi.StringArray{
			pulumi.String("0.0.0.0/0"),
		}
	}

	// Locate a GCP project
	project, err := organizations.LookupProject(ctx, &organizations.LookupProjectArgs{
		ProjectId: &projectId,
	})
	if err != nil {
		return nil, []pulumi.Resource{}, err
	}
	// Create service account
	serviceAccount, err := serviceaccount.NewAccount(ctx, "Create Service Account", &serviceaccount.AccountArgs{
		AccountId: pulumi.Sprintf("account-%s", strings.ReplaceAll(deploymentId, "_", "-")),
		Project:   pulumi.String(*project.ProjectId),
	})
	if err != nil {
		return nil, []pulumi.Resource{}, err
	}

	// Create network for VMs
	network, err := compute.NewNetwork(ctx, "Create Network", &compute.NetworkArgs{
		Name:                  pulumi.Sprintf("network-%s", deploymentId),
		Project:               pulumi.String(*project.ProjectId),
		AutoCreateSubnetworks: pulumi.Bool(false),
	})
	if err != nil {
		return nil, []pulumi.Resource{}, err
	}

	subnetwork, err := compute.NewSubnetwork(ctx, "Create Subnetwork", &compute.SubnetworkArgs{
		Name:        pulumi.Sprintf("subnetwork-%s", deploymentId),
		Project:     pulumi.String(*project.ProjectId),
		IpCidrRange: pulumi.String("10.0.1.0/24"),
		Region:      pulumi.String(networkRegion),
		Network:     network.ID(),
	}, pulumi.DependsOn(
		[]pulumi.Resource{
			network,
		},
	))
	if err != nil {
		return nil, []pulumi.Resource{}, err
	}

	// Create firewall rules for ssh and compute
	sshFirewall, err := compute.NewFirewall(ctx, "Create SSH Firewall",
		&compute.FirewallArgs{
			Name:         pulumi.Sprintf("network-ssh-firewall-%s", deploymentId),
			Project:      pulumi.String(*project.ProjectId),
			Network:      network.SelfLink,
			SourceRanges: sshAllowCIDRs,
			Direction:    pulumi.String("INGRESS"),
			Allows: &compute.FirewallAllowArray{
				&compute.FirewallAllowArgs{
					Protocol: pulumi.String("tcp"),
					Ports: pulumi.StringArray{
						pulumi.String("22"),
					},
				},
			},
		}, pulumi.DependsOn(
			[]pulumi.Resource{
				network, subnetwork,
			},
		),
	)
	if err != nil {
		return nil, []pulumi.Resource{}, err
	}

	computeFirewall, err := compute.NewFirewall(ctx, "Create Compute Firewall",
		&compute.FirewallArgs{
			Name:    pulumi.Sprintf("network-firewall-%s", deploymentId),
			Project: pulumi.String(*project.ProjectId),
			Network: network.SelfLink,
			SourceRanges: pulumi.StringArray{
				pulumi.String("0.0.0.0/0"),
			},
			Direction: pulumi.String("INGRESS"),
			Allows: &compute.FirewallAllowArray{
				&compute.FirewallAllowArgs{
					Protocol: pulumi.String("tcp"),
					Ports: pulumi.StringArray{
						pulumi.String("26656"),
					},
				},
			},
		}, pulumi.DependsOn(
			[]pulumi.Resource{
				network, subnetwork,
			},
		),
	)
	if err != nil {
		return nil, []pulumi.Resource{}, err
	}

	// Create public ip address
	pubIp, err := compute.NewAddress(ctx, "Create Public IP Address", &compute.AddressArgs{
		Name:    pulumi.Sprintf("node-ip-%s", deploymentId),
		Project: pulumi.String(*project.ProjectId),
		Region:  pulumi.String(networkRegion),
	})
	if err != nil {
		return nil, []pulumi.Resource{}, err
	}

	vm, err := compute.NewInstance(ctx, "Create VM Instance", &compute.InstanceArgs{
		Name:        pulumi.Sprintf("node-vm-%s", deploymentId),
		Project:     pulumi.String(*project.ProjectId),
		MachineType: pulumi.String(vmHardware),
		Zone:        pulumi.String(vmZone),
		BootDisk: &compute.InstanceBootDiskArgs{
			InitializeParams: &compute.InstanceBootDiskInitializeParamsArgs{
				Image: pulumi.String("ubuntu-2004-focal-v20220204"),
				Size:  pulumi.Int(vmDiskSize),
			},
		},
		Metadata: pulumi.StringMap{
			"enable-oslogin": pulumi.String("false"),
			"ssh-keys":       pulumi.Sprintf("%s:%s", vmUsername, publicKey),
		},
		NetworkInterfaces: compute.InstanceNetworkInterfaceArray{
			&compute.InstanceNetworkInterfaceArgs{
				Subnetwork: subnetwork.ID(),
				AccessConfigs: &compute.InstanceNetworkInterfaceAccessConfigArray{
					&compute.InstanceNetworkInterfaceAccessConfigArgs{
						NatIp: pubIp.Address,
					},
				},
			},
		},
		ServiceAccount: &compute.InstanceServiceAccountArgs{
			Email: serviceAccount.Email,
			Scopes: pulumi.StringArray{
				pulumi.String("cloud-platform"),
			},
		},
	}, pulumi.DependsOn(
		[]pulumi.Resource{
			sshFirewall, computeFirewall,
		},
	))
	if err != nil {
		return nil, []pulumi.Resource{}, err
	}

	// Send files and run commands to setup node
	connection := remote.ConnectionArgs{
		Host:       pubIp.Address,
		User:       pulumi.String(vmUsername),
		PrivateKey: privateKey,
	}

	ctx.Export("Virtual Machine ID", vm.ID())
	ctx.Export("Virtual Machine IP", pubIp.Address)
	ctx.Export("SSH Endpoint", pulumi.Sprintf("%s@%s", vmUsername, pubIp.Address))

	return &connection, []pulumi.Resource{vm, pubIp}, nil
}
