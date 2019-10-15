VERSION_MAJOR ?= 0
VERSION_MINOR ?= 1
VERSION_BUILD ?= 0
VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)

ORG := github.com
OWNER := cloud-native-taiwan
REPOPATH ?= $(ORG)/$(OWNER)/controller101

$(shell mkdir -p ./out)

############
# Building #
############

.PHONY: build
build: out/controller

.PHONY: out/controller
out/controller:
	go build -ldflags="-s -w -X $(REPOPATH)/pkg/version.version=$(VERSION)" \
	  -a -o $@ cmd/main.go

.PHONY: build_image
build_image:
	docker build -t kairen/controller101:$(VERSION) .

##############
# Generating #
##############

vendor:
	go mod vendor

.PHONY: verify-codegen
verify-codegen: vendor 
	./hack/k8s/verify-codegen.sh
  
.PHONY: codegen
codegen: vendor
	./hack/k8s/update-generated.sh
	
###########
# Testing #
###########

.PHONY: test
test:
	./hack/test-go.sh

.PHONY: clean
clean:
	rm -rf out/ vendor/