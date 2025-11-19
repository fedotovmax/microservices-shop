SERVICE =
DC = docker compose -f dev.docker-compose.yml --env-file .env

up: 
	$(DC) up -d

down: 
	$(DC) down

build: 
	$(DC) build

build-service:
	@test -n "$(SERVICE)" || (echo "Error: SERVICE is empty"; exit 1)
	$(DC) build $(SERVICE)

restart: down up 

stop:
	$(DC) stop

start:
	$(DC) start

clean: down ## Stop and remove containers, networks, volumes
	$(DC) down -v --remove-orphans



