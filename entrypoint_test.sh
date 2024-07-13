#!/bin/sh

atlas schema apply --env local --auto-approve
go run cmd/ezsplit_jet/*.go
gotestsum -- -p=1 ./...
