FROM golang:1.16 AS builder

WORKDIR /build
ADD . .
RUN make build-docker && \
    cd bin && \
    mv *_docker rocketchat-jira-webhook

FROM alpine

WORKDIR /rocketchat-jira-webhook
COPY --from=builder /build/bin .
COPY examples/minimal.toml /etc/rocketchat-jira-webhook/config.toml
CMD ["./rocketchat-jira-webhook"]