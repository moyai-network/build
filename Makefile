.PHONY: all build run raw-run test tidy migrate sqlc submodules sync-flyway lint push cloc

submodules:
	git submodule update --init --recursive

sqlc:
	$(MAKE) -C external/flyway sqlc

migrate:
	$(MAKE) -C external/flyway local-migrate

raw-run:
	go run ./cmd

run: submodules sqlc
	go run ./cmd

test:
	go test ./cmd ./internal/...

tidy:
	go mod tidy

sync-flyway:
	@if [ -f scripts/sync-flyway.sh ]; then bash scripts/sync-flyway.sh; else echo "scripts/sync-flyway.sh not found"; fi

lint:
	@if [ -f scripts/lint.sh ]; then bash scripts/lint.sh; else go test ./cmd ./internal/...; fi

push: sync-flyway
	@if [ -f scripts/push.sh ]; then bash scripts/push.sh; else echo "scripts/push.sh not found"; fi

cloc:
	@if [ -f scripts/cloc.sh ]; then bash scripts/cloc.sh; else echo "scripts/cloc.sh not found"; fi
