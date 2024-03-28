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