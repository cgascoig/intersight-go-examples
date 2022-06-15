SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
# .DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX = >

GO_MODULE := github.com/cgascoig/intersight-prometheus-exporter
DOCKER_IMAGE_ID := ghcr.io/cgascoig/intersight-prometheus-exporter

GO_CMD ?= go
GO_BUILD_CMD := $(GO_CMD) build -v 
GO_BUILD_FLAGS := -ldflags "-X main.commit=`git rev-parse HEAD`"
GO_PATH ?= $(shell go env GOPATH)

all: build/list-ntp-policies build/alarm-streamer build/workflow-runner
.PHONY: all

clean:
> rm -Rf build
.PHONY: clean

build/list-ntp-policies: go.mod go.sum Makefile $(shell find list-ntp-policies -name \*.go -type f)
> mkdir -p $(@D)
> go build -o "$@" ./list-ntp-policies

build/alarm-streamer: go.mod go.sum Makefile $(shell find alarm-streamer -name \*.go -type f)
> mkdir -p $(@D)
> go build -o "$@" ./alarm-streamer

build/workflow-runner: go.mod go.sum Makefile $(shell find workflow-runner -name \*.go -type f)
> mkdir -p $(@D)
> go build -o "$@" ./workflow-runner
