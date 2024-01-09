# Pull the base image
FROM bufbuild/buf:1.28.1

# Install protoc
RUN apt-get update && apt-get -y install protobuf-compiler

# Set environment variables for Go
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPATH=/go

# Install buf and other tools
RUN go install github.com/bufbuild/buf/cmd/buf@v1.28.1 && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.0 && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.0

# Set the working directory in the docker container
WORKDIR /workspace

COPY ./proto proto

# Set the entrypoint using the full path
ENTRYPOINT ["/go/bin/buf"]