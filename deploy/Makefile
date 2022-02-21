LIKED_VERSION ?= "1.2.0-rc3"
COSMOVISOR_VERSION ?= "1.1.0"
MONIKER ?= ""
GENESIS_URL ?= ""

LIKED_WORKDIR ?= $(HOME)
LIKED_HOME ?= "$(LIKED_WORKDIR)/.liked"
LIKED_USER ?= "$(USER)"

.PHONY: setup-pulumi
setup-pulumi:
	pulumi stack init $(PULUMI_STACK)
	cp ./Pulumi.config.yaml.template "./Pulumi.$(PULUMI_STACK).yaml"

.PHONY: deploy
deploy: 
	pulumi up

.PHONY: destroy
destroy:
	pulumi destroy 

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
