PKG ?= github.com/PLAB-IO/aws-cloudformation-debugger
BIN ?= cfdbg
VERSION ?= master
ARCH ?= $(shell go env GOOS)-$(shell go env GOARCH)

platform_temp = $(subst -, ,$(ARCH))
GOOS = $(word 1, $(platform_temp))
GOARCH = $(word 2, $(platform_temp))

# Code Style
fmt:
	go fmt **/*.go

lint:
	golint ...

sanity: fmt lint

build: build-dirs
	GOOS=$(GOOS) \
	GOARCH=$(GOARCH) \
	VERSION=$(VERSION) \
	PKG=$(PKG) \
	BIN=$(BIN) \
	go build \
        -o bin/$(GOOS)/$(GOARCH)/$(BIN)

build-dirs:
	@mkdir -p bin/$(GOOS)/$(GOARCH)

test: build
	./bin/$(GOOS)/$(GOARCH)/cfdbg --profile $(profile) --stack-name cfdbg-demo

###################
###################

deploy-demo:
	aws --profile $(profile) \
	  cloudformation package \
		--template-file demo/main.yml \
		--s3-bucket $(bucket) \
		--output-template-file .cf-main-output.yml

	aws --profile $(profile) \
	  cloudformation deploy \
		--template-file .cf-main-output.yml \
		--capabilities CAPABILITY_NAMED_IAM CAPABILITY_AUTO_EXPAND \
		--stack-name cfdbg-demo
