LIKED_VERSION ?= "1.2.0-rc3"
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

.PHONY: setup-node
setup-node:	
ifeq ($(MONIKER), "")
	@echo "Missing MONIKER env"
	@echo "Usage: MONIKER=<your moniker> make setup-node"
else
	LIKED_VERSION=$(LIKED_VERSION) COSMOVISOR_VERSION=$(COSMOVISOR_VERSION) ./node-setup.sh $(MONIKER) "$(GENESIS_URL)" $(LIKED_WORKDIR) $(LIKED_HOME) $(LIKED_USER) "$(LIKED_SEED_NODES)"
endif 
	
.PHONY: initialize-systemctl
initialize-systemctl: 
	sudo cp ./liked.service /etc/systemd/system/liked.service
	sudo systemctl daemon-reload

.PHONY: start-node
start-node:
	sudo systemctl start liked

.PHONY: stop-node
stop-node:
	sudo systemctl stop liked
