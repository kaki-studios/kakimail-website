# syntax=docker/dockerfile:1

# Accept the Go version for the image to be set as a build argument.
# Default to Go 1.22.1
ARG GO_VERSION=1.22.1

# First stage: build the executable.
FROM golang:${GO_VERSION} AS builder

# Git is required for fetching the dependencies.
RUN apt-get update -y && apt-get install -y ca-certificates git libsqlite3-dev && update-ca-certificates

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.mod ./go.sum ./
RUN go mod download

# Import the code from the context.
COPY ./ ./

# Build the executable to `/app`
RUN go build \
    -installsuffix 'static' \
    -o /app .


# Expose the ports to our application
EXPOSE 8001
EXPOSE 8000

# Mount the certificate cache directory as a volume, so it remains even after
# we deploy a new version
VOLUME ["/cert-cache"]

# Run the compiled binary.
ENTRYPOINT ["/app"]
