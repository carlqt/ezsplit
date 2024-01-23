.PHONY: up
up:
	docker compose up

down:
	docker compose down

run:
	go run ./cmd/ezsplit/main.go

atlas_apply:
	atlas schema apply --env local

graphql_playground:
	go run ./cmd/graphql_playground/main.go

schema_generate:
	go run github.com/99designs/gqlgen generate