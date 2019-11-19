#!/bin/sh

set -e

yamllint -f parsable ${TARGETPATH} | /usr/bin/yamllint-action
