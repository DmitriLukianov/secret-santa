include .env
export 

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d postgres

env-down:
	@docker compose down postgres

env-cleanup:
	@read -p "очистить все volume файлы окружения? Опасность утери данных.[y/n]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down postgres && \
		rm -rf out/pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

migrate-create:
	@if  [ -z "$(seq)" ]; then \
	echo "Отсутствует необходимый параметр seq. Пример make migrate-create seq=init"; \
	exit 1; \
	fi;
	docker compose run --rm migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if  [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action. Пример make migrate-action action=up"; \
		exit 1; \
	fi;
	docker compose run --rm migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"