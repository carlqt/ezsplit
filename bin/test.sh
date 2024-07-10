#!/bin/sh

atlas schema apply --env local --auto-approve

jet -dsn=postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable -schema=public -path=./.gen

go test -p=1 -v ./...
