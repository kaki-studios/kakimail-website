# syntax=docker/dockerfile:1

FROM golang:1.22.1

WORKDIR /app

COPY . .

RUN apt-get update -y && apt-get install -y libsqlite3-dev
RUN go mod download
RUN go build

EXPOSE 8000
EXPOSE 8001

VOLUME ["/certs"]


ENTRYPOINT "./kakimail-website"
