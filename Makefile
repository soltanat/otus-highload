gen:
	oapi-codegen --config ./oapi-server.cfg.yaml ./_oas/openapi.json
	oapi-codegen --config ./oapi-types.cfg.yaml ./_oas/openapi.json

lint:
	golangci-lint run --fix

up-async:
	docker-compose -f ./postgres-async-docker-compose.yaml -f ./docker-compose.yaml up -d --build
	docker-compose -f ./monitoring/monitoring.yaml up -d

down-async:
	docker-compose -f ./postgres-async-docker-compose.yaml -f ./docker-compose.yaml down -v
	docker-compose -f ./monitoring/monitoring.yaml down -v


up-sync:
	docker-compose -f ./postgres-sync-docker-compose.yaml up -d --build

down-sync:
	docker-compose -f ./postgres-sync-docker-compose.yaml down -v


run-gentransactions:
	go run ./cmd/gentransactions

up-postgres:
	docker-compose -f ./postgres-docker-compose.yaml up -d postgres

up-app:
	docker-compose -f ./postgres-docker-compose.yaml up -d app

down:
	docker-compose -f ./postgres-docker-compose.yaml down -v


migrate-up:
	 goose -dir ./migrations postgres "host=localhost user=postgres password=postgres dbname=social sslmode=disable" up

migrate-down:
	 goose -dir ./migrations postgres "host=localhost user=postgres password=postgres dbname=social sslmode=disable" down

upload:
	go run ./cmd/upload -d postgres://postgres:postgres@localhost:5432/social?sslmode=disable