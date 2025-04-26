.DELETE_ON_ERROR:
.SUFFIXES:

# TODO: Flag changes don't apply immediately in GNU make 4.3 as included in
# Debian Bookworm (https://unix.stackexchange.com/a/784954).
MAKEFLAGS += --no-builtin-rules --warn-undefined-variables

.DEFAULT_GOAL := all

SHELL = /bin/sh
.SHELLFLAGS += -e

CURL := curl
DOCKER := docker
GO := go
INSTALL := install
TAR := tar
