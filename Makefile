TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=github.com
NAMESPACE=spring-media
NAME=terraform-provider-awsmt
BINARY=${NAME}
VERSION=0.0.1
OS_ARCH=$(shell go env GOOS)_$(shell go env GOHOSTARCH)
BUILD_DIR=build
SWEEP_DIR?=./awsmt
SWEEP?=eu-central-1

default: install

clean:
	rm -rf ./${BUILD_DIR}/

build-app:
	go build -o ./${BUILD_DIR}/${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./build/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./build/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./build/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./build/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./build/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./build/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./build/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./build/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./build/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./build/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./build/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./build/${BINARY}_${VERSION}_windows_amd64

install: build-app
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ./build/${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

$(BUILD_DIR):
	mkdir $(BUILD_DIR)

%/coverage.profile: $(BUILD_DIR)
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -coverprofile $@

%/coverage_func.txt: %/coverage.profile
	go tool cover -func=$< -o $@

%/coverage.html: %/coverage.profile
	go tool cover -html=$< -o $@

test: $(clean) $(BUILD_DIR)/coverage.html $(BUILD_DIR)/coverage_func.txt $(BUILD_DIR)/coverage.profile
	@echo "finished test"

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(SWEEP_DIR) -v -sweep=$(SWEEP) $(SWEEPARGS) -timeout 60m