
LIKED_VERSION ?= "1.2.0"
COSMOVISOR_VERSION ?= "1.1.0"
MONIKER ?= ""
GENESIS_URL ?= ""
LIKED_SEED_NODES ?= ""

LIKED_WORKDIR ?= $(HOME)
LIKED_HOME ?= "$(LIKED_WORKDIR)/.liked"
LIKED_USER ?= "$(USER)"
LIKED_REPO_SOURCE ?= "likecoin"

.PHONY: setup-node
setup-node:
ifeq ($(MONIKER), "")
	@echo "Missing MONIKER env"
	@echo "Usage: make setup-node MONIKER=<your moniker> GENESIS_URL=<url> LIKED_SEED_NODES=<seed nodes>"
else ifeq ($(GENESIS_URL), "")
	@echo "Missing GENESIS_URL env"
	@echo "Usage: make setup-node MONIKER=<your moniker> GENESIS_URL=<url> LIKED_SEED_NODES=<seed nodes>"
else ifeq ($(LIKED_SEED_NODES), "")
	@echo "Missing SEED_NODES env"
	@echo "Usage: make setup-node MONIKER=<your moniker> GENESIS_URL=<url> LIKED_SEED_NODES=<seed nodes>"
else
	LIKED_VERSION=$(LIKED_VERSION) COSMOVISOR_VERSION=$(COSMOVISOR_VERSION) LIKED_REPO_SOURCE=$(LIKED_REPO_SOURCE) ./node-setup.sh "$(MONIKER)" "$(GENESIS_URL)" "$(LIKED_SEED_NODES)" "$(LIKED_WORKDIR)" "$(LIKED_HOME)" "$(LIKED_USER)" 
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
