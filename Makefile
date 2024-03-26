gen:
	oapi-codegen --config ./oapi-server.cfg.yaml ./_oas/openapi.json
	oapi-codegen --config ./oapi-types.cfg.yaml ./_oas/openapi.json

lint:
	golangci-lint run --fix

up:
	docker-compose up -d --build

down:
	docker-compose down -v