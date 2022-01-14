PROJECT_NAME=$(shell basename "$(PWD)")

tidy:
	go mod tidy

vet:
	go vet ./...

test:
	go test ./... -v

BINARIES_DIRECTORY=bin
MAIN_FILE=cmd/main.go
CONFIG_FILE=examples/minimal.toml
LDFLAGS="-w -s"

clean:
	rm -rf ${BINARIES_DIRECTORY}

run:
	go run ${MAIN_FILE} -config ${CONFIG_FILE}

build-docker: vet clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags=${LDFLAGS} -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}_docker ${MAIN_FILE}
