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

pushgithub:
	@echo "Enter commit message:"; \
	read COMMIT_MSG; \
	if [ -z "$$COMMIT_MSG" ]; then \
		echo "Commit message cannot be empty"; \
		exit 1; \
	fi; \
	git add .; \
	git commit -m "$$COMMIT_MSG"; \
	git push; \
	echo "✅ Всё готово! Проект отправлен на github"


