.PHONY: help

help: ## Print this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

gen-webhooks: ## Generate Webhooks Client
	oapi-codegen -generate types,client -package webhooks ./epilot-webhooks-client/openapi.yml > epilot-webhooks-client/gen_client.go

build: ## Build the binary
	go build -o provider

download-spec: ## Download the OpenAPI Spec
	wget https://docs.api.epilot.io/webhooks.yaml -O ./epilot-webhooks-client/openapi.yml