NAME=gosaic
VERSION=$(shell cat version)
BUILD_TIME=$(shell date)
BUILD_USER=$(shell whoami)
BUILD_HASH=$(shell git rev-parse HEAD 2>/dev/null || echo "")
ARCH=amd64
OS=linux darwin

LDFLAGS=-ldflags "-X 'github.com/atongen/gosaic/environment.Version=$(VERSION)' \
				          -X 'github.com/atongen/gosaic/environment.BuildTime=$(BUILD_TIME)' \
									-X 'github.com/atongen/gosaic/environment.BuildUser=$(BUILD_USER)' \
									-X 'github.com/atongen/gosaic/environment.BuildHash=$(BUILD_HASH)'"

all: clean test build

clean:
	go clean
	@rm -f `which ${NAME}`

vet:
	go vet `go list ./... | grep -v /vendor/`

test:
	go test -cover `go list ./... | grep -v /vendor/`

build: test
	go install ${LDFLAGS}

distclean:
	@mkdir -p dist
	rm -rf dist/*

dist: test distclean
	for arch in ${ARCH}; do \
		for os in ${OS}; do \
			env GOOS=$${os} GOARCH=$${arch} go build -v ${LDFLAGS} -o dist/${NAME}-${VERSION}-$${os}-$${arch}; \
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
			tar czf dist/${NAME}-${VERSION}-$${os}-$${arch}.tar.gz -C dist ${NAME}-${VERSION}-$${os}-$${arch} ${NAME}-${VERSION}-$${os}-$${arch}.asc; \
		done; \
	done; \
	find dist/ -type f  ! -name "*.tar.gz" -delete

tag:
	scripts/tag.sh

upload:
	if [ ! -z "\${GITHUB_TOKEN}" ]; then \
		ghr -t "${GITHUB_TOKEN}" -u ${BUILD_USER} -r ${NAME} -replace ${VERSION} dist/; \
	fi

release: tag package upload

.PHONY: all clean vet test build distclean dist sign package tag release
