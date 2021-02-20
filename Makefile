COMMIT=$(shell git rev-parse --short HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
TAG=$(shell git describe --tags |cut -d- -f1)

LDFLAGS = -ldflags "-X github.com/aint/CloudElephant/cmd.gitTag=${TAG} -X github.com/aint/CloudElephant/cmd.gitCommit=${COMMIT} -X github.com/aint/CloudElephant/cmd.gitBranch=${BRANCH}"

build:
	go build ${LDFLAGS} -v -o ce .

test:
	go test -v ./...
