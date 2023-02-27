ALL_ARCH = amd64 arm64

.EXPORT_ALL_VARIABLES:

all: $(addprefix build-arch-,$(ALL_ARCH))

VERSION_MAJOR ?= 1
VERSION_MINOR ?= 26
VERSION_BUILD ?= 2
TAG?=v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)
FLAGS=
ENVVAR=CGO_ENABLED=0
LDFLAGS?=-s
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
REGISTRY?=fred78290
BUILD_DATE?=`date +%Y-%m-%dT%H:%M:%SZ`
VERSION_LDFLAGS=-X main.phVersion=$(TAG)
IMAGE=$(REGISTRY)/vmware-desktop-autoscaler

deps:
	go mod vendor
#	wget "https://raw.githubusercontent.com/Fred78290/autoscaler/master/cluster-autoscaler/cloudprovider/grpc/grpc.proto" -O grpc/grpc.proto
#	protoc -I . -I vendor grpc/grpc.proto --go_out=plugins=grpc:.

build: $(addprefix build-arch-,$(ALL_ARCH))

build-arch-%: deps clean-arch-%
	$(ENVVAR) GOOS=$(GOOS) GOARCH=$* go build -buildvcs=false -ldflags="-X main.phVersion=$(TAG) -X main.phBuildDate=$(BUILD_DATE) ${LDFLAGS}" -a -o out/$(GOOS)/$*/vmware-desktop-autoscaler

test-unit: clean build
	bash ./scripts/run-tests.sh

make-image: $(addprefix make-image-arch-,$(ALL_ARCH))

make-image-arch-%:
	docker build --pull -t ${IMAGE}-$*:${TAG} -f Dockerfile.$* .
	@echo "Image ${TAG}-$* completed"

push-image: $(addprefix push-image-arch-,$(ALL_ARCH))

push-image-arch-%:
	docker push ${IMAGE}-$*:${TAG}

push-manifest:
	docker buildx build --pull --platform linux/amd64,linux/arm64 --push -t ${IMAGE}:${TAG} .
	@echo "Image ${TAG}* completed"

container-push-manifest: container push-manifest

clean: $(addprefix clean-arch-,$(ALL_ARCH))

clean-arch-%:
	rm -f ./out/$(GOOS)/$*/vmware-desktop-autoscaler

docker-builder:
	docker build -t kubernetes-desktop-autoscaler-builder ./builder

build-in-docker: $(addprefix build-in-docker-arch-,$(ALL_ARCH))

build-in-docker-arch-%: clean-arch-% docker-builder
	docker run --rm -v `pwd`:/gopath/src/github.com/Fred78290/vmware-desktop-autoscaler/ kubernetes-desktop-autoscaler-builder:latest bash \
		-c 'cd /gopath/src/github.com/Fred78290/vmware-desktop-autoscaler \
		&& BUILD_TAGS=${BUILD_TAGS} make -e REGISTRY=${REGISTRY} -e TAG=${TAG} -e BUILD_DATE=`date +%Y-%m-%dT%H:%M:%SZ` build-arch-$*'

container: $(addprefix container-arch-,$(ALL_ARCH))

container-arch-%: build-in-docker-arch-%
	@echo "Full in-docker image ${TAG}-$* completed"

test-in-docker: docker-builder
	docker run --rm -v `pwd`:/gopath/src/github.com/Fred78290/kubernetes-desktop-autoscaler/ kubernetes-desktop-autoscaler-builder:latest bash \
		-c 'cd /gopath/src/github.com/Fred78290/kubernetes-desktop-autoscaler && bash ./scripts/run-tests.sh'

.PHONY: all build test-unit clean docker-builder build-in-docker push-image push-manifest
