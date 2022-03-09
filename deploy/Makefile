STACK ?= validator
CLOUD ?= azure

.PHONY: setup-pulumi
setup-pulumi:
	pulumi stack init $(STACK)
	cp ./Pulumi.$(CLOUD).yaml.template "./Pulumi.$(STACK).yaml"

.PHONY: ssh-key
ssh-key: 
	ssh-keygen -t rsa -f id_rsa -m PEM
	pulumi config set -s $(STACK) \
		likecoin-skynet:resource-group-name $(RESOURCE_GROUP)
	pulumi config set -s $(STACK) \
		likecoin-skynet:vm-ssh-allow-list $(SSH_WHITELIST)
	pulumi config set -s $(STACK) \
		likecoin-skynet:vm-password --secret $(PASSWORD)
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
	make -C cmd/deploy/assets setup-node

.PHONY: initialize-systemctl
initialize-systemctl:
	make -C cmd/deploy/assets initialize-systemctl

.PHONY: start-node
start-node:
	make -C cmd/deploy/assets start-node

.PHONY: stop-node
stop-node:
	make -C cmd/deploy/assets stop-node
