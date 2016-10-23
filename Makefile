VERSION=$(shell cat version)
BUILD_TIME=$(shell date)
BUILD_USER=$(shell whoami)
BUILD_HASH=$(shell git rev-parse HEAD)
ARCH=amd64
OS=linux darwin

LDFLAGS=-ldflags "-X 'gosaic/environment.Version=$(VERSION)' -X 'gosaic/environment.BuildTime=$(BUILD_TIME)' -X 'gosaic/environment.BuildUser=$(BUILD_USER)' -X 'gosaic/environment.BuildHash=$(BUILD_HASH)'"

all: deps test build

clean:
	rm -rf bin/* pkg/*

deps:
	go get -u github.com/constabulary/gb/...

test:
	gb test all

build:
	gb build ${LDFLAGS} all

distclean:
	@mkdir -p dist
	rm -rf dist/*

dist: test distclean
	for arch in ${ARCH}; do \
		for os in ${OS}; do \
			env GOOS=$${os} GOARCH=$${arch} gb build ${LDFLAGS} all; \
			mv bin/gosaic-$${os}-$${arch} dist/gosaic-${VERSION}-$${os}-$${arch}; \
		done; \
	done

sign: dist
	$(eval key := $(shell git config --get user.signingkey))
	for file in dist/*; do \
		gpg2 --armor --local-user ${key} --detach-sign $${file}; \
	done

package: sign
	for arch in ${ARCH}; do \
		for os in ${OS}; do \
			tar czf dist/gosaic-${VERSION}-$${os}-$${arch}.tar.gz -C dist gosaic-${VERSION}-$${os}-$${arch} gosaic-${VERSION}-$${os}-$${arch}.asc; \
		done; \
	done

.PHONY: all clean deps test build distclean dist sign package
