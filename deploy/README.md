# Setting up likecoin node on cloud services using Pulumi

This document describe how to use Pulumi to set up devnet that support development work. Although it may use at production environment, it is not recommended.

# Prerequisites

- [Pulumi](https://www.pulumi.com/docs/get-started/install/)
- Go 1.17+
- Azure CLI (for azure nodes)

# Setup

## Account setup

Login to Azure CLI and Pulumi CLI with the following commands

```
az login
pulumi login
```

Alternatively you can explore local setup via , `pulumi login --local`. However, we are not going to cover in this document.

Pulumi will use the currently logged in session of `az` command to perform following action. Please use `az account --set` to correctly set the default account.

## Environment and secret preparation

Prepare the variable for ease of setup, we may want to change according to the current cloud/testnet situation.

```
export RESOURCE_GROUP=likecoin-skynet
export REGION=southeastasia
export PASSWORD=$(openssl rand -hex 10)
export STACK=validator
echo $PASSWORD
```

Create a resource group on Azure with the following command

```
az group create --location $REGION --resource-group $RESOURCE_GROUP
```

Run the following command to setup a pulumi stack.

```
make setup-pulumi STACK=$STACK
```

This creates a file `Pulumi.$STACK.yaml` for our stack in which we will modify for configuring the deployment to work.

We will create a SSH keypair with the following command, we will be using RSA as per instructions provided by [Azure](https://docs.microsoft.com/en-us/azure/virtual-machines/linux/ssh-from-windows#create-an-ssh-key-pair).

**Note: Azure Pulumi does not seem to work with keys that has a passphrase due to the prompt being missed out during deployment hence we will be using one with empty passphrase.**

Create deployment SSH key and inflate the respective value in `Pulmi.$STACK.yaml`

```
make ssh-key
```

We should now see the values being assigned in `Pulumi.$STACK.yaml`

## Provision of resources

After Run the following command to execute the deployment

```
make deploy
```

Pulumi will execute a dry-run deployment to validate the deployment script, review and confirm the deployment.

After a successful deployment, the public ID of the deploy vm should be printed on console.

We can connect to the virtual machine via SSH as follow.

```
ssh -i id_rsa likecoin@<vm ip address>
```

Simply run the following command to start the service

```
make start-node
```

## Configurations

Pulumi stack configurations that is used by the deployment script

| Configuration                       | Description                                  | Mandatory |
| ----------------------------------- | -------------------------------------------- | --------- |
| likecoin-skynet:node-genesis        | URL to the genesis.json file                 | ❌        |
| likecoin-skynet:node-moniker        | Moniker identifier of the node               | ✅        |
| likecoin-skynet:node-seeds          | Comma separated P2P Seed nodes               | ❌        |
| likecoin-skynet:resource-group-name | Resource group name for Azure                | ✅        |
| likecoin-skynet:vm-username         | Admin username to the Virtual Machine        | ✅        |
| likecoin-skynet:vm-password         | Admin password to the Virtual Machine        | ✅        |
| likecoin-skynet:vm-private-key      | SSH private key to the Virtual Machine       | ✅        |
| likecoin-skynet:vm-public-key       | SSH public key to the Virtual Machine        | ✅        |
| likecoin-skynet:vm-ssh-allow-list   | Comma separated CIDR list for SSH white list | ❌        |
