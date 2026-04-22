-include .env
export

up:
	docker compose up -d --build

down:
	docker compose down

rebuild:
	docker compose down
	docker compose up -d --build

migrate-up:
	docker compose run --rm migrate \
		-path /migrations \
		-database "postgres://$${POSTGRES_USER}:$${POSTGRES_PASSWORD}@postgres:5432/$${POSTGRES_DB}?sslmode=disable" \
		up

migrate-down:
	docker compose run --rm migrate \
		-path /migrations \
		-database "postgres://$${POSTGRES_USER}:$${POSTGRES_PASSWORD}@postgres:5432/$${POSTGRES_DB}?sslmode=disable" \
		down 1

db-clean:
	docker compose exec postgres psql -U $${POSTGRES_USER} -d $${POSTGRES_DB} -c \
		"TRUNCATE users, events, participants, assignments, invitations, wishlists, wishlist_items, messages, email_verification_codes RESTART IDENTITY CASCADE;"

db-reset:
	docker compose down -v
	docker compose up -d --build
