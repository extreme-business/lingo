# Pull the base image
FROM golang:1.21

ARG DST_DIR=/protobufs/gen/pb-go

# Install protoc
RUN apt-get update && apt-get -y install protobuf-compiler

# Install protoc-gen-go and protoc-gen-go-grpc plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

# Set environment variables for Go
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory in the docker container
RUN mkdir -p /app

# Copy the proto files into the docker container
COPY ./proto /proto
COPY ./provision/protoc/run.sh /app/run.sh
RUN chmod +x /app/run.sh

RUN mkdir -p /proto
RUN mkdir -p /proto/gen
RUN mkdir -p /proto/gen/pb-go

ENTRYPOINT ["sh", "/app/run.sh"]