.PHONY: 

run:
	go run cmd/restful/main.go

test:
	go test ./... -cover -v

dependencies:
	docker compose -f script/docker-compose.yml up -d

dependencies-down:
	docker compose -f script/docker-compose.yml down

migrate-create:
	@read -p  "Migration name (eg:create_users, alter_entities, ...): " NAME; \
	migrate create -ext sql -dir migrations -seq $$NAME

migrate-up:
	migrate -database ${POSTGRESQL_URL} -path migrations up

migrate-down:
	migrate -database ${POSTGRESQL_URL} -path migrations down -all
