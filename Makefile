ENV ?= dev

DC = docker compose --env-file .$(ENV).env

up: 
	$(DC) up -d

down: 
	$(DC) down

build: 
	$(DC) build

restart: down up 

stop:
	$(DC) stop

clean: down ## Stop and remove containers, networks, volumes
	$(DC) down -v --remove-orphans

dev: ## Start development environment
	@$(MAKE) up ENV=dev


# shell: ## Access container shell
# 	$(DC) exec app sh

# # Database operations
# db-backup: ## Backup database
# 	$(DC) exec db pg_dump -U postgres mydb > backup_$(shell date +%Y%m%d).sql

# db-restore: ## Restore database
# 	$(DC) exec -T db psql -U postgres mydb < $(FILE)