# this dockerfile is used to build the docker image for the apps in the apps folder.
# syntax=docker/dockerfile:1.2

FROM golang:1.21-alpine AS base
ENV GO111MODULE="on"
ENV GOOS="linux"
ENV CGO_ENABLED=0
RUN apk update \
    && apk add --no-cache \
    ca-certificates \
    curl \
    tzdata \
    git \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*
COPY . /src/lingo
WORKDIR /src/lingo
RUN --mount=type=cache,target=/go/pkg/mod go mod download
RUN go generate ./...

FROM base AS debug
WORKDIR /src/apps/lingo
RUN go install github.com/go-delve/delve/cmd/dlv@v1.21.0
EXPOSE 8080
EXPOSE 2345
ENTRYPOINT ["dlv", "debug", "--continue", "--headless", "--listen=:2345", "--api-version=2", "--accept-multiclient", "--log"]

FROM base AS lint
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1
RUN golangci-lint run ./...

FROM base AS builder
COPY --from=base /src /src
RUN go build -o /src/bin/lingo .

FROM gcr.io/distroless/static-debian11 as prod
COPY --from=builder /src/bin/lingo ./
EXPOSE 8080
ENTRYPOINT ["./lingo"]