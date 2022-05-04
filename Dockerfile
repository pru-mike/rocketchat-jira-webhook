FROM golang:1.16 AS builder

WORKDIR /build
ADD . .
RUN make build-docker

FROM alpine

ARG CONFIG="examples/minimal.toml"
WORKDIR /rocketchat-jira-webhook
COPY --from=builder /build/bin .
COPY $CONFIG /etc/rocketchat-jira-webhook/config.toml
CMD ["./rocketchat-jira-webhook"]