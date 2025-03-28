SHELL = /usr/bin/env bash -o pipefail

default: help

.PHONY: help
help:
	# Usage:
	@sed -n '/^\([a-z][^:]*\).*/s//    make \1/p' $(MAKEFILE_LIST)

.PHONY: build
build:
	# Generate the themes
	@echo "Generating Themes..."
	@go run main.go

.PHONY: deploy
deploy:
	# Deploy the themes
	@make build
	@git add *; \
	read -p "Enter commit message: " msg; \
	git commit -m "$$msg"
	@echo "Deploying Themes..."
	@npm version minor
	@git push
	@vsce package
	@vsce publish
	@echo "Themes deployed successfully!"
