FROM golang:1.21.4-bookworm AS build

WORKDIR /app

RUN apt-get update
RUN apt-get install -yq libhyperscan-dev libpcap-dev

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/ ./cmd

RUN CGO_ENABLED=1 GOOS=linux go build -x -o /ginder cmd/hs_cmd.go
#cmd/ginder_cli.go

FROM ubuntu:latest

WORKDIR /

COPY --from=build /ginder /ginder

RUN apt-get update
RUN apt-get install -yq libhyperscan-dev libpcap-dev

ENTRYPOINT ["/ginder"]

