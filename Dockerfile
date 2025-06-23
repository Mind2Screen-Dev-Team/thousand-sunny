# Stage 1: Build
FROM golang:1.24-alpine3.21 AS builder

ARG APP_NAME=core
ARG APP_PORT=8080

# Install required dependencies for building
RUN apk add --no-cache git

# Set environment variables for Go build
ENV CGO_ENABLED=0 GOOS=linux

# Create and set the working directory
WORKDIR /app

# Copy Go module manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main ./cmd/$APP_NAME

# Stage 2: Final runtime
FROM alpine:3.21

RUN apk add --no-cache tzdata curl busybox-extras bash

WORKDIR /app

ENV PATH=/app/bin:$PATH

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main /app/main
COPY --from=builder /app/config.yaml /app/cfg-backup/config.yaml

EXPOSE $APP_PORT

# Set the entrypoint
ENTRYPOINT ["/app/main", "-cfg=config.yaml"]
