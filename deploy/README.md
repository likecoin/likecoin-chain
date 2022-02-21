# Setting up likecoin node on cloud services using Pulumi

# Prerequisites

- [Pulumi](https://www.pulumi.com/docs/get-started/install/)
- Go 1.17+
- Azure CLI (for azure nodes)

# Setup

## Azure

Login to Azure CLI and Pulumi CLI with the following commands

```
az login

pulumi login
```

or if you wish to use Pulumi without an account in which you will need to enter a passphrase for secret storing

```
az login

pulumi login --local
export PULUMI_CONFIG_PASSPHRASE=<passphrase>
```

Create a resource group on Azure with the following command

```
az group create --location <your location> --resource-group <your resource group name>
```

Pulumi should be able to capture your login session and perform deployments on your behalf.

Run the following command to setup a pulumi stack.

```
PULUMI_STACK=<stack name> make setup-pulumi
```

This creates a config file `Pulumi.<stack name>.yaml` for your stack in which you will have to modify for the deployment to work.

First, you can create a SSH keypair with the following command, we will be using RSA as per instructions provided by [Azure](https://docs.microsoft.com/en-us/azure/virtual-machines/linux/ssh-from-windows#create-an-ssh-key-pair).

**Note: Azure Pulumi does not seem to work with keys that has a passphrase due to the prompt being missed out during deployment hence we will be using one with empty passphrase.**

```
ssh-keygen -t rsa -f rsa -m PEM
```

Run the following commands to setup configurations for your stack

```
pulumi config set node-deployment:resource-group-name <your resource group name>
pulumi config set node-deployment:vm-password --secret <your vm password>

cat rsa.pub | pulumi config set node-deployment:vm-public-key --
cat rsa | pulumi config set node-deployment:vm-private-key --secret --
```

You should now see the values being assigned in `Pulumi.<stack name>.yaml`

You may now run the following command to execute the deployment

```
make deploy
```

Pulumi will execute a dry-run deployment to validate the deployment script, You may select `Yes` to confirm the deployment.

After a successful deployment, connect to the virtual machine via SSH to the IP address output displayed on screen

```
ssh -i rsa <vm username>@<vm ip address>
```

Simply run the following command to start the service

```
make start-node
```

## Configurations

Pulumi stack configurations that is used by the deployment script

| Configuration                       | Description                                  | Mandatory |
| ----------------------------------- | -------------------------------------------- | --------- |
| node-deployment:node-genesis        | URL to the genesis.json file                 | ❌        |
| node-deployment:node-moniker        | Moniker identifier of the node               | ✅        |
| node-deployment:node-seeds          | Comma separated P2P Seed nodes               | ❌        |
| node-deployment:resource-group-name | Resource group name for Azure                | ✅        |
| node-deployment:vm-username         | Admin username to the Virtual Machine        | ✅        |
| node-deployment:vm-password         | Admin password to the Virtual Machine        | ✅        |
| node-deployment:vm-private-key      | SSH private key to the Virtual Machine       | ✅        |
| node-deployment:vm-public-key       | SSH public key to the Virtual Machine        | ✅        |
| node-deployment:vm-ssh-allow-list   | Comma separated CIDR list for SSH white list | ❌        |
