ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=postgres password=postgres dbname=postgres host=localhost port=5432 sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)/internal/pkg
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/db/migrations

.PHONY: test-migration-up
test-migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: test-migration-down
test-migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

.PHONY: generate-proto
generate-proto:
	protoc --proto_path=proto/ --go_out=internal/app/homework/pb --go-grpc_out=internal/app/homework/pb proto/*.proto

.PHONY: docker-compose
docker-up:
	docker-compose up

docker-cleanup:
	docker-compose down
