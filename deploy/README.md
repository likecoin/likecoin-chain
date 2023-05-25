# Setting up likecoin node on cloud services using Pulumi

This document describe how to use Pulumi to set up devnet that support development work. Although it may use at production environment, it is not recommended.

# Prerequisites

- [Pulumi](https://www.pulumi.com/docs/get-started/install/)
- Go 1.19
- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli-macos) (for azure nodes)
- [Gcloud CLI](https://cloud.google.com/sdk/docs/install-sdk) (for gcp nodes)

# Setup

## Account setup

Login to Azure CLI/ GCloud CLI and Pulumi CLI with the following commands

```
# For Azure
az login
# For GCP
gcloud auth application-default login
gcloud config set project $GCP_PROJECT_ID

pulumi login
```

Alternatively you can explore local setup via , `pulumi login --local`. However, we are not going to cover in this document.

Pulumi will use the currently logged in session of `az` command to perform following action. Please use `az account --set` to correctly set the default account.

## Environment and secret preparation

Prepare the variable for ease of setup, we may want to change according to the current cloud/testnet situation.

```
# For Azure
export RESOURCE_GROUP=likecoin-skynet
export REGION=southeastasia
export CLOUD=azure

# For GCP
export PROJECT_ID=likecoin-skynet
export CLOUD=gcp

# source cidr/ip addresses for ssh access, comma separated, default current ip address
export SSH_WHITELIST=$(curl checkip.amazonaws.com)

export PASSWORD=$(openssl rand -base64 48 | sed -e 's/[\/|=|+]//g')
export STACK=validator
echo $PASSWORD
```

Create a resource group on Azure with the following command

```
# For Azure
az group create --location $REGION --resource-group $RESOURCE_GROUP
```

Run the following command to setup a pulumi stack.

```
make setup-pulumi STACK=$STACK CLOUD=$CLOUD
```

This creates a file `Pulumi.$STACK.yaml` for our stack in which we will modify for configuring the deployment to work.

We will create a SSH keypair with the following command, we will be using RSA as per instructions provided by [Azure](https://docs.microsoft.com/en-us/azure/virtual-machines/linux/ssh-from-windows#create-an-ssh-key-pair).

**Note: Azure Pulumi does not seem to work with keys that has a passphrase due to the prompt being missed out during deployment hence we will be using one with empty passphrase.**

Create deployment SSH key and inflate the respective value in `Pulmi.$STACK.yaml`

```
make ssh-key
```

We should now see the values being assigned in `Pulumi.$STACK.yaml`

**NOTE: It is highly recommended to set a value for the config `vm-ssh-allow-list` to prevent global ssh access to the virtual machine. If any assigned ip address is dynamically allocated by the ISP, newly rotated ip address shall be updated on Azure Portal/Google Cloud Platform to maintain said ssh access.**

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

### Common

| Configuration                     | Description                                  | Mandatory |
| --------------------------------- | -------------------------------------------- | --------- |
| likecoin-skynet:cloud-provider    | Cloud provider (azure, gcp)                  | ✅        |
| likecoin-skynet:node-genesis      | URL to the genesis.json file                 | ✅        |
| likecoin-skynet:node-moniker      | Moniker identifier of the node               | ✅        |
| likecoin-skynet:node-seeds        | Comma separated P2P Seed nodes               | ✅        |
| likecoin-skynet:vm-username       | Admin username to the Virtual Machine        | ✅        |
| likecoin-skynet:vm-private-key    | SSH private key to the Virtual Machine       | ✅        |
| likecoin-skynet:vm-public-key     | SSH public key to the Virtual Machine        | ✅        |
| likecoin-skynet:vm-ssh-allow-list | Comma separated CIDR list for SSH white list | ❌        |
| likecoin-skynet:vm-disk-size      | Disk size for the Virtual Machine            | ✅        |
| likecoin-skynet:vm-hardware       | Hardware type for the Virtual Machine        | ✅        |

### Azure Specific

| Configuration                       | Description                           | Mandatory |
| ----------------------------------- | ------------------------------------- | --------- |
| likecoin-skynet:resource-group-name | Resource group name for Azure         | ✅        |
| likecoin-skynet:vm-password         | Admin password to the Virtual Machine | ✅        |

### GCP Specific

| Configuration                  | Description                          | Mandatory |
| ------------------------------ | ------------------------------------ | --------- |
| likecoin-skynet:project-id     | Project ID for GCP                   | ✅        |
| likecoin-skynet:vm-zone        | Zone for Virtual Machine             | ✅        |
| likecoin-skynet:network-region | Region for Network related resources | ✅        |
