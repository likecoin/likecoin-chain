package main

import (
	hashUtils "deploy/pkg/utils/hash"
	"strings"

	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/compute"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		resourceGroupName := cfg.Require("resource-group-name")
		vmUsername := cfg.Require("vm-username")
		vmPassword := cfg.RequireSecret("vm-password")
		publicKey := cfg.Require("vm-public-key")
		privateKey := cfg.RequireSecret("vm-private-key")
		sshAllowList := cfg.Require("vm-ssh-allow-list")

		genesisUrl := cfg.Require("node-genesis")
		seeds := cfg.Require("node-seeds")
		moniker := cfg.Require("node-moniker")

		stackName := ctx.Stack()

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
			return err
		}

		// Create network for VMs
		virtualNetwork, err := network.NewVirtualNetwork(
			ctx,
			"Create Node Network",
			&network.VirtualNetworkArgs{
				VirtualNetworkName: pulumi.Sprintf("vn-%s", stackName),
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
			return err
		}

		// Create public ip address
		pubIp, err := network.NewPublicIPAddress(ctx, "Create Public IP Address", &network.PublicIPAddressArgs{
			PublicIpAddressName:      pulumi.Sprintf("node-ip-%s", stackName),
			ResourceGroupName:        pulumi.String(resourceGroup.Name),
			PublicIPAllocationMethod: pulumi.String(network.IPAllocationMethodDynamic),
		})
		if err != nil {
			return err
		}

		sshSecurityRule := network.SecurityRuleTypeArgs{
			Access:                   pulumi.String(network.AccessAllow),
			Protocol:                 pulumi.String("*"),
			SourcePortRange:          pulumi.String("*"),
			DestinationAddressPrefix: pulumi.String("*"),
			DestinationPortRange:     pulumi.String("22"),
			Direction:                pulumi.String(network.SecurityRuleDirectionInbound),
			Name:                     pulumi.Sprintf("ssh-inbound-%s", stackName),
			Priority:                 pulumi.Int(102),
		}

		if len(sshAllowCIDRs) > 0 {
			sshSecurityRule.SourceAddressPrefixes = sshAllowCIDRs
		} else {
			sshSecurityRule.SourceAddressPrefix = pulumi.String("*")
		}
		// Create network security group
		networkSg, err := network.NewNetworkSecurityGroup(ctx, "Create Network Security Group", &network.NetworkSecurityGroupArgs{
			NetworkSecurityGroupName: pulumi.Sprintf("network-sg-%s", stackName),
			ResourceGroupName:        pulumi.String(resourceGroup.Name),
			SecurityRules: network.SecurityRuleTypeArray{
				// Inbound rule for rpc and peer connections
				network.SecurityRuleTypeArgs{
					Access:                   pulumi.String(network.AccessAllow),
					Protocol:                 pulumi.String("*"),
					SourceAddressPrefix:      pulumi.String("*"),
					SourcePortRange:          pulumi.String("*"),
					DestinationAddressPrefix: pulumi.String("*"),
					DestinationPortRange:     pulumi.String("26656-26657"),
					Direction:                pulumi.String(network.SecurityRuleDirectionInbound),
					Name:                     pulumi.Sprintf("network-inbound-%s", stackName),
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
					Name:                     pulumi.Sprintf("network-outbound-%s", stackName),
					Priority:                 pulumi.Int(101),
				},
				sshSecurityRule,
			},
		})
		if err != nil {
			return err
		}

		// Create network interface with previously created ip address and security group
		networkIf, err := network.NewNetworkInterface(ctx, "Create Network Interface", &network.NetworkInterfaceArgs{
			NetworkInterfaceName: pulumi.Sprintf("network-if-%s", stackName),
			ResourceGroupName:    pulumi.String(resourceGroup.Name),
			IpConfigurations: network.NetworkInterfaceIPConfigurationArray{
				network.NetworkInterfaceIPConfigurationArgs{
					Name:                      pulumi.Sprintf("node-ipcfg-%s", stackName),
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
			return err
		}

		// Create virtual machine for node to run on
		vm, err := compute.NewVirtualMachine(ctx, "Create Virtual Machine", &compute.VirtualMachineArgs{
			VmName:            pulumi.Sprintf("node-vm-%s", stackName),
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
				VmSize: pulumi.String("Standard_B2s"),
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
					Name:         pulumi.Sprintf("node-vm-os-disk-%s", stackName),
					DeleteOption: pulumi.String(compute.DeleteOptionsDelete),
				},
			},
		}, pulumi.DependsOn([]pulumi.Resource{
			networkIf,
		}),
		)
		if err != nil {
			return err
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

		// TODO: Extract this part as it is common on all cloud services
		// Send files and run commands to setup node
		connection := remote.ConnectionArgs{
			Host:       ipAddress,
			User:       pulumi.String(vmUsername),
			PrivateKey: privateKey,
		}

		scriptChange, err := hashUtils.GetFileHash("../../node-setup.sh")
		if err != nil {
			return err
		}

		cpScript, err := remote.NewCopyFile(ctx, "Copy deployment script to VM", &remote.CopyFileArgs{
			Connection: connection,
			LocalPath:  pulumi.String("../../node-setup.sh"),
			RemotePath: pulumi.Sprintf("/home/%s/node-setup.sh", vmUsername),
			Triggers: pulumi.Array{
				pulumi.String(scriptChange),
			},
		}, pulumi.DependsOn([]pulumi.Resource{
			vm,
			pubIp,
		}))
		if err != nil {
			return err
		}

		chmodScript, err := remote.NewCommand(ctx, "Set deployment script permission", &remote.CommandArgs{
			Connection: connection,
			Create:     pulumi.Sprintf("chmod u+x /home/%s/node-setup.sh", vmUsername),
			Triggers: pulumi.Array{
				pulumi.String(scriptChange),
			},
		}, pulumi.DependsOn([]pulumi.Resource{
			cpScript,
		}))
		if err != nil {
			return err
		}

		serviceTemplateChange, err := hashUtils.GetFileHash("../../liked.service.template")
		if err != nil {
			return err
		}

		cpService, err := remote.NewCopyFile(ctx, "Copy node service template to VM", &remote.CopyFileArgs{
			Connection: connection,
			LocalPath:  pulumi.String("../../liked.service.template"),
			RemotePath: pulumi.Sprintf("/home/%s/liked.service.template", vmUsername),
			Triggers: pulumi.Array{
				pulumi.String(serviceTemplateChange),
			},
		}, pulumi.DependsOn([]pulumi.Resource{
			vm,
			pubIp,
		}))
		if err != nil {
			return err
		}

		makeFileChange, err := hashUtils.GetFileHash("../../Makefile")
		if err != nil {
			return err
		}

		cpMakefile, err := remote.NewCopyFile(ctx, "Copy Makefile to VM", &remote.CopyFileArgs{
			Connection: connection,
			LocalPath:  pulumi.String("../../Makefile"),
			RemotePath: pulumi.Sprintf("/home/%s/Makefile", vmUsername),
			Triggers: pulumi.Array{
				pulumi.String(makeFileChange),
			},
		}, pulumi.DependsOn([]pulumi.Resource{
			vm,
			pubIp,
		}))
		if err != nil {
			return err
		}

		installBuildEssentialsCommand, err := remote.NewCommand(ctx, "Install make on VM", &remote.CommandArgs{
			Connection: connection,
			Create:     pulumi.String("sudo apt install -y make"),
		}, pulumi.DependsOn([]pulumi.Resource{
			cpScript,
			cpService,
			chmodScript,
			cpMakefile,
		}))
		if err != nil {
			return err
		}

		setupNodeCommand, err := remote.NewCommand(ctx, "Execute setup node command", &remote.CommandArgs{
			Connection: connection,
			Create:     pulumi.Sprintf("MONIKER=%s LIKED_USER=%s LIKED_WORKDIR=/home/%s GENESIS_URL=%s LIKED_SEED_NODES=%s make setup-node", moniker, vmUsername, vmUsername, genesisUrl, seeds),
			Triggers: pulumi.Array{
				pulumi.String(scriptChange),
				pulumi.String(serviceTemplateChange),
				pulumi.String(makeFileChange),
			},
		}, pulumi.DependsOn([]pulumi.Resource{
			installBuildEssentialsCommand,
			cpScript,
			cpService,
			chmodScript,
			cpMakefile,
		}))
		if err != nil {
			return err
		}

		_, err = remote.NewCommand(ctx, "Execute initialize node service command", &remote.CommandArgs{
			Connection: connection,
			Create:     pulumi.String("make initialize-systemctl"),
			Triggers: pulumi.Array{
				pulumi.String(scriptChange),
				pulumi.String(serviceTemplateChange),
				pulumi.String(makeFileChange),
			},
		}, pulumi.DependsOn([]pulumi.Resource{
			installBuildEssentialsCommand,
			setupNodeCommand,
			cpScript,
			cpService,
			chmodScript,
			cpMakefile,
		}))
		if err != nil {
			return err
		}

		ctx.Export("Stack Name", pulumi.String(stackName))
		ctx.Export("Virtual Machine ID", vm.ID())
		ctx.Export("Virtual Machine IP", ipAddress)
		ctx.Export("SSH Endpoint", pulumi.Sprintf("%s@%s", vmUsername, ipAddress))

		return nil
	})
}
