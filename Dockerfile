# syntax=docker/dockerfile:1

ARG GO_VERSION=1.20
ARG DOCKERD_VERSION=20.10.14

FROM golang:${GO_VERSION}-alpine as build-deployx
RUN mkdir /build
WORKDIR /build

COPY go.mod go.sum /build/
RUN go mod download && go mod verify

COPY version /build/version
COPY cmd /build/cmd
COPY commands /build/commands
COPY convert /build/convert
COPY deploy /build/deploy

RUN go build -o docker-deployx cmd/deployx/*

FROM docker:${DOCKERD_VERSION} AS dockerd-release

COPY --from=build-deployx /build/docker-deployx /usr/local/bin
RUN mkdir -p /usr/local/lib/docker/cli-plugins && ln -s /usr/local/bin/docker-deployx /usr/local/lib/docker/cli-plugins/docker-deployx
