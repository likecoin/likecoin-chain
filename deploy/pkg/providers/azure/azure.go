package providers

import (
	"strings"

	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/compute"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func CreateNodeVM(deploymentId string, ctx *pulumi.Context, cfg *config.Config) (*remote.ConnectionArgs, []pulumi.Resource, error) {
	resourceGroupName := cfg.Require("resource-group-name")
	vmUsername := cfg.Require("vm-username")
	vmPassword := cfg.RequireSecret("vm-password")
	publicKey := cfg.Require("vm-public-key")
	privateKey := cfg.RequireSecret("vm-private-key")
	sshAllowList := cfg.Require("vm-ssh-allow-list")
	vmHardware := cfg.Require("vm-hardware")
	vmDiskSize := cfg.GetInt("vm-disk-size")
	if vmDiskSize == 0 {
		vmDiskSize = int(50)
	}

	var sshAllowCIDRs pulumi.StringArray
	for _, cidr := range strings.Split(sshAllowList, ",") {
		_cidr := pulumi.String(strings.TrimSpace(cidr))
		if len(_cidr) > 0 {
			sshAllowCIDRs = append(sshAllowCIDRs, _cidr)
		}
	}

	// Locate an Azure Resource Group
	resourceGroup, err := resources.LookupResourceGroup(ctx, &resources.LookupResourceGroupArgs{
		ResourceGroupName: resourceGroupName,
	})
	if err != nil {
		return nil, []pulumi.Resource{}, err
	}

	// Create network for VMs
	virtualNetwork, err := network.NewVirtualNetwork(
		ctx,
		"Create Node Network",
		&network.VirtualNetworkArgs{
			VirtualNetworkName: pulumi.Sprintf("vn-%s", deploymentId),
			ResourceGroupName:  pulumi.String(resourceGroup.Name),
			AddressSpace: &network.AddressSpaceArgs{
				AddressPrefixes: pulumi.StringArray{
					pulumi.String("10.0.0.0/16"),
				},
			},
			Subnets: network.SubnetTypeArray{
				network.SubnetTypeArgs{
					Name:          pulumi.String("default"),
					AddressPrefix: pulumi.String("10.0.1.0/24"),
				},
			},
		},
	)
	if err != nil {
		return nil, []pulumi.Resource{}, err
	}

	// Create public ip address
	pubIp, err := network.NewPublicIPAddress(ctx, "Create Public IP Address", &network.PublicIPAddressArgs{
		PublicIpAddressName:      pulumi.Sprintf("node-ip-%s", deploymentId),
		ResourceGroupName:        pulumi.String(resourceGroup.Name),
		PublicIPAllocationMethod: pulumi.String(network.IPAllocationMethodDynamic),
	})
	if err != nil {
		return nil, []pulumi.Resource{}, err
	}

	sshSecurityRule := network.SecurityRuleTypeArgs{
		Access:                   pulumi.String(network.AccessAllow),
		Protocol:                 pulumi.String("*"),
		SourcePortRange:          pulumi.String("*"),
		DestinationAddressPrefix: pulumi.String("*"),
		DestinationPortRange:     pulumi.String("22"),
		Direction:                pulumi.String(network.SecurityRuleDirectionInbound),
		Name:                     pulumi.Sprintf("ssh-inbound-%s", deploymentId),
		Priority:                 pulumi.Int(102),
	}

	if len(sshAllowCIDRs) > 0 {
		sshSecurityRule.SourceAddressPrefixes = sshAllowCIDRs
	} else {
		sshSecurityRule.SourceAddressPrefix = pulumi.String("*")
	}
	// Create network security group
	networkSg, err := network.NewNetworkSecurityGroup(ctx, "Create Network Security Group", &network.NetworkSecurityGroupArgs{
		NetworkSecurityGroupName: pulumi.Sprintf("network-sg-%s", deploymentId),
		ResourceGroupName:        pulumi.String(resourceGroup.Name),
		SecurityRules: network.SecurityRuleTypeArray{
			// Inbound rule for rpc and peer connections
			network.SecurityRuleTypeArgs{
				Access:                   pulumi.String(network.AccessAllow),
				Protocol:                 pulumi.String("*"),
				SourceAddressPrefix:      pulumi.String("*"),
				SourcePortRange:          pulumi.String("*"),
				DestinationAddressPrefix: pulumi.String("*"),
				DestinationPortRange:     pulumi.String("26656"),
				Direction:                pulumi.String(network.SecurityRuleDirectionInbound),
				Name:                     pulumi.Sprintf("network-inbound-%s", deploymentId),
				Priority:                 pulumi.Int(100),
			},
			// Outbound rule without port limit
			network.SecurityRuleTypeArgs{
				Access:                   pulumi.String(network.AccessAllow),
				Protocol:                 pulumi.String("*"),
				SourceAddressPrefix:      pulumi.String("*"),
				SourcePortRange:          pulumi.String("*"),
				DestinationAddressPrefix: pulumi.String("*"),
				DestinationPortRange:     pulumi.String("*"),
				Direction:                pulumi.String(network.SecurityRuleDirectionOutbound),
				Name:                     pulumi.Sprintf("network-outbound-%s", deploymentId),
				Priority:                 pulumi.Int(101),
			},
			sshSecurityRule,
		},
	})
	if err != nil {
		return nil, []pulumi.Resource{}, err
	}

	// Create network interface with previously created ip address and security group
	networkIf, err := network.NewNetworkInterface(ctx, "Create Network Interface", &network.NetworkInterfaceArgs{
		NetworkInterfaceName: pulumi.Sprintf("network-if-%s", deploymentId),
		ResourceGroupName:    pulumi.String(resourceGroup.Name),
		IpConfigurations: network.NetworkInterfaceIPConfigurationArray{
			network.NetworkInterfaceIPConfigurationArgs{
				Name:                      pulumi.Sprintf("node-ipcfg-%s", deploymentId),
				PrivateIPAllocationMethod: pulumi.String(network.IPAllocationMethodDynamic),
				Subnet: network.SubnetTypeArgs{
					Id: virtualNetwork.Subnets.Index(pulumi.Int(0)).Id(),
				},
				PublicIPAddress: network.PublicIPAddressTypeArgs{
					Id: pubIp.ID(),
				},
			},
		},
		NetworkSecurityGroup: network.NetworkSecurityGroupTypeArgs{
			Id: networkSg.ID(),
		},
	}, pulumi.DependsOn([]pulumi.Resource{
		networkSg,
	}))
	if err != nil {
		return nil, []pulumi.Resource{}, err
	}

	// Create virtual machine for node to run on
	vm, err := compute.NewVirtualMachine(ctx, "Create Virtual Machine", &compute.VirtualMachineArgs{
		VmName:            pulumi.Sprintf("node-vm-%s", deploymentId),
		ResourceGroupName: pulumi.String(resourceGroup.Name),
		NetworkProfile: compute.NetworkProfileArgs{
			NetworkInterfaces: compute.NetworkInterfaceReferenceArray{
				compute.NetworkInterfaceReferenceArgs{
					Id:      networkIf.ID(),
					Primary: pulumi.Bool(true),
				},
			},
		},
		HardwareProfile: &compute.HardwareProfileArgs{
			VmSize: pulumi.String(vmHardware),
		},
		Location: pulumi.String(resourceGroup.Location),
		OsProfile: &compute.OSProfileArgs{
			AdminUsername: pulumi.String(vmUsername),
			AdminPassword: vmPassword,
			ComputerName:  pulumi.String(vmUsername),
			LinuxConfiguration: &compute.LinuxConfigurationArgs{
				DisablePasswordAuthentication: pulumi.Bool(true),
				PatchSettings: &compute.LinuxPatchSettingsArgs{
					AssessmentMode: pulumi.String(compute.LinuxPatchAssessmentModeImageDefault),
				},
				ProvisionVMAgent: pulumi.Bool(true),
				Ssh: &compute.SshConfigurationArgs{
					PublicKeys: &compute.SshPublicKeyTypeArray{
						compute.SshPublicKeyTypeArgs{
							KeyData: pulumi.String(publicKey),
							Path:    pulumi.Sprintf("/home/%s/.ssh/authorized_keys", vmUsername),
						},
					},
				},
			},
		},
		StorageProfile: &compute.StorageProfileArgs{
			ImageReference: &compute.ImageReferenceArgs{
				Offer:     pulumi.String("0001-com-ubuntu-server-focal"),
				Publisher: pulumi.String("Canonical"),
				Sku:       pulumi.String("20_04-lts-gen2"),
				Version:   pulumi.String("latest"),
			},
			OsDisk: &compute.OSDiskArgs{
				Caching:      compute.CachingTypesReadWrite,
				CreateOption: pulumi.String(compute.DiskCreateOptionFromImage),
				ManagedDisk: &compute.ManagedDiskParametersArgs{
					StorageAccountType: pulumi.String(compute.StorageAccountType_Premium_LRS),
				},
				DiskSizeGB:   pulumi.Int(vmDiskSize),
				Name:         pulumi.Sprintf("node-vm-os-disk-%s", deploymentId),
				DeleteOption: pulumi.String(compute.DeleteOptionsDelete),
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{
		networkIf,
	}),
	)
	if err != nil {
		return nil, []pulumi.Resource{}, err
	}

	// Ensure the resources are created first
	ready := pulumi.All(vm.ID(), pubIp.Name, pulumi.String(resourceGroup.Name))

	ipAddress := ready.ApplyT(func(args []interface{}) (string, error) {
		name := args[1].(string)
		resourceGroupName := args[2].(string)
		ip, err := network.LookupPublicIPAddress(ctx, &network.LookupPublicIPAddressArgs{
			ResourceGroupName:   resourceGroupName,
			PublicIpAddressName: name,
		})
		if err != nil {
			return "", err
		}
		return *ip.IpAddress, nil
	}).(pulumi.StringOutput)

	// Send files and run commands to setup node
	connection := remote.ConnectionArgs{
		Host:       ipAddress,
		User:       pulumi.String(vmUsername),
		PrivateKey: privateKey,
	}

	ctx.Export("Virtual Machine ID", vm.ID())
	ctx.Export("Virtual Machine IP", ipAddress)
	ctx.Export("SSH Endpoint", pulumi.Sprintf("%s@%s", vmUsername, ipAddress))

	return &connection, []pulumi.Resource{vm, pubIp}, nil
}
