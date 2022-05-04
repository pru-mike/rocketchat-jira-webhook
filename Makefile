CONFIG_FILE ?= examples/minimal.toml
LDFLAGS="-w -s"

tidy:
	go mod tidy

vet:
	go vet ./...

test:
	go test ./... -v

clean:
	rm -f bin/rocketchat-jira-webhook

run:
	go run cmd/main.go -config ${CONFIG_FILE}

build-docker: clean vet
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -a -ldflags=${LDFLAGS} -o bin/rocketchat-jira-webhook cmd/main.go
