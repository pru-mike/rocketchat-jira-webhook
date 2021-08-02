FROM golang:latest
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o rocketchat-jira-webhook
RUN go install
VOLUME /app
WORKDIR /app
EXPOSE 4567
CMD ["rocketchat-jira-webhook", "-config", "config.toml"]