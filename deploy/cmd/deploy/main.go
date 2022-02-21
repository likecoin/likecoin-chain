package main

import (
	providers "deploy/pkg/providers"
	hashUtils "deploy/pkg/utils/hash"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")

		cloudProviderName := cfg.Require("cloud-provider")

		vmUsername := cfg.Require("vm-username")
		genesisUrl := cfg.Require("node-genesis")
		seeds := cfg.Require("node-seeds")
		moniker := cfg.Require("node-moniker")

		stackName := ctx.Stack()

		cloudProvider := providers.Provider{
			Name: cloudProviderName,
		}

		CreateNodeVM, err := cloudProvider.GetVMBuilder()
		if err != nil {
			return err
		}

		connectionArgs, dependencies, err := CreateNodeVM(stackName, ctx, cfg)
		if err != nil {
			return err
		}

		scriptChange, err := hashUtils.GetFileHash("./assets/node-setup.sh")
		if err != nil {
			return err
		}

		cpScript, err := remote.NewCopyFile(ctx, "Copy deployment script to VM", &remote.CopyFileArgs{
			Connection: connectionArgs,
			LocalPath:  pulumi.String("./assets/node-setup.sh"),
			RemotePath: pulumi.Sprintf("/home/%s/node-setup.sh", vmUsername),
			Triggers: pulumi.Array{
				pulumi.String(scriptChange),
			},
		}, pulumi.DependsOn(dependencies))
		if err != nil {
			return err
		}

		chmodScript, err := remote.NewCommand(ctx, "Set deployment script permission", &remote.CommandArgs{
			Connection: connectionArgs,
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

		serviceTemplateChange, err := hashUtils.GetFileHash("./assets/liked.service.template")
		if err != nil {
			return err
		}

		cpService, err := remote.NewCopyFile(ctx, "Copy node service template to VM", &remote.CopyFileArgs{
			Connection: connectionArgs,
			LocalPath:  pulumi.String("./assets/liked.service.template"),
			RemotePath: pulumi.Sprintf("/home/%s/liked.service.template", vmUsername),
			Triggers: pulumi.Array{
				pulumi.String(serviceTemplateChange),
			},
		}, pulumi.DependsOn(dependencies))
		if err != nil {
			return err
		}

		makeFileChange, err := hashUtils.GetFileHash("./assets/Makefile")
		if err != nil {
			return err
		}

		cpMakefile, err := remote.NewCopyFile(ctx, "Copy Makefile to VM", &remote.CopyFileArgs{
			Connection: connectionArgs,
			LocalPath:  pulumi.String("./assets/Makefile"),
			RemotePath: pulumi.Sprintf("/home/%s/Makefile", vmUsername),
			Triggers: pulumi.Array{
				pulumi.String(makeFileChange),
			},
		}, pulumi.DependsOn(dependencies))
		if err != nil {
			return err
		}

		installBuildEssentialsCommand, err := remote.NewCommand(ctx, "Install make on VM", &remote.CommandArgs{
			Connection: connectionArgs,
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
			Connection: connectionArgs,
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
			Connection: connectionArgs,
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

		return nil
	})
}
