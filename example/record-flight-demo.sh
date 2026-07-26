#!/bin/sh
set -eu

cd "$(dirname "$0")"
go run . -record-demo=flight
