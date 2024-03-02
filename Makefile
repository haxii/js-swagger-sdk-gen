VERSION=$(shell git describe --abbrev=0 --tags)
BUILD=$(shell git rev-parse --short HEAD)
CURDIR=$(shell pwd)

# Inject the build version (commit hash) into the executable.
LDFLAGS := -ldflags "-s -w -X main.Build=$(BUILD) -X main.Version=$(VERSION)"

.PHONY: release
release: depdir
	gox -verbose $(LDFLAGS) -osarch="linux/amd64 windows/amd64 darwin/arm64" \
	 -output="./dist/js-swagger-sdk-gen-$(VERSION)-{{.OS}}-{{.Arch}}/js-swagger-sdk-gen" ./cmd/js-swagger-sdk-gen/
	cd dist; \
		ls | while read line; do tar czvf ../$${line}.tar.gz $${line}; done

.PHONY: clean
clean:
	rm -rf bin dist *.tar.gz

depdir:
	mkdir -p bin dist
