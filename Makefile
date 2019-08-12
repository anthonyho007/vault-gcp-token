BIN_NAME=vault-gcp-token
VAULT_CONFIG=$(PWD)/plugins.hcl
TMP_DIR="/private/tmp/vault-plugin"
PLUGIN_PATH="gcp-token"

export GO111MODULE=on
export VAULT_DEV_ROOT_TOKEN_ID="root"


build:
	go build -o $(BIN_NAME)

server:
	vault server -dev -dev-root-token-id=$(VAULT_DEV_ROOT_TOKEN_ID) -config=$(VAULT_CONFIG)

install-plugin:
	mkdir -p $(TMP_DIR)
	cp $(BIN_NAME) $(TMP_DIR)
	$(eval SHASUM=$(shell shasum -a 256 "$(TMP_DIR)/$(BIN_NAME)" | cut -d " " -f1))
	vault write sys/plugins/catalog/$(BIN_NAME) sha_256=$(SHASUM) command=$(BIN_NAME)
	vault secrets enable --plugin-name=$(BIN_NAME) --path=$(PLUGIN_PATH) plugin

test:
	go test ./...
