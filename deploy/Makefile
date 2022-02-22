COSMOVISOR_VERSION ?= "1.1.0"
MONIKER ?= ""
GENESIS_URL ?= ""

LIKED_WORKDIR ?= $(HOME)
LIKED_HOME ?= "$(LIKED_WORKDIR)/.liked"
LIKED_USER ?= "$(USER)"
STACK ?= validator

.PHONY: setup-pulumi
setup-pulumi:
	pulumi stack init $(STACK)
	cp ./Pulumi.config.yaml.template "./Pulumi.$(STACK).yaml"

.PHONY: ssh-key
ssh-key: 
	ssh-keygen -t rsa -f id_rsa -m PEM
	pulumi config set -s $(STACK) \
		likecoin-skynet:resource-group-name $(RESOURCE_GROUP)
	pulumi config set -s $(STACK) \
		like-skynet:vm-password --secret $(PASSWORD)
	cat id_rsa.pub | pulumi config set likecoin-skynet:vm-public-key --
	cat id_rsa | pulumi config set likecoin-skynet:vm-private-key --secret --

.PHONY: deploy
deploy: 
	pulumi up -s $(STACK)

.PHONY: destroy
destroy:
	pulumi destroy -s $(STACK)
