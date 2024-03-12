.PHONY: proto

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

help: 
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_0-9-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

update-submodules: ## Update submodules (sync proto)
	git submodule update --remote --merge

binary: ## Build binary bin/hawkeye
		mkdir -p bin
		go build -o bin/hawkeye hawkeye/main.go

clean: ## Clean bin directory
		rm -rf bin

proto: ## Generate Go code from proto files
	protoc --go_out=. --go_opt=Mproto/intent.proto=pkg/api --go-grpc_out=. --go-grpc_opt=Mproto/intent.proto=pkg/api proto/*.proto --experimental_allow_proto3_optional