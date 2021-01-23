setup:
	docker-compose exec server bundle exec rails db:create
	docker-compose exec server bundle exec rails db:migrate
	docker-compose exec server bundle exec rails db:seed

migrate:
	docker-compose exec server bundle exec rails db:migrate

rspec:
	docker-compose exec server rspec

console:
	docker-compose exec server /bin/bash -c "bundle exec rails console"

start:
	docker-compose up

start_daemon:
	docker-compose up -d

build_and_start:
	docker-compose up --build

stop:
	docker-compose down
