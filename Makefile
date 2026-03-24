include .env
export

export PROJECT_DIR=$(shell pwd)

env_up:
	docker-compose up -d todo-postgres

env_down:
	docker-compose down todo-postgres

env_clean:
	@read -p "Are you sure you want to clean volumes, [y/N]: " answer; \
	if [ "$$answer" = "y" ]; then \
		docker-compose down --remove-orphans todo-postgres && \
		rm -rf out/pgdata && \
		echo "Done"; \
	else \
		echo "Canceled"; fi

migrate-create:
	@if [ -z "$(seq)"]; then \
			echo "'seq' required" \
			exit 1; \
  		else \
  		  docker-compose run --rm todo-postgres-migrate \
          		create \
          		-ext sql \
          		-dir /migrations \
          		-seq "$(seq)"; \
	fi; \

migrate-up:
	make migrate-action action=up

migrate-down:
	make migrate-action action=down

migrate-action:
	@if [ -z "$(action)"]; then \
  	echo "'action' required" \
    			exit 1; \
    			fi; \
	docker-compose run --rm todo-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todo-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

env-port-forward:
	docker-compose up -d port-forwarder

env-port-close:
	docker-compose down port-forwarder