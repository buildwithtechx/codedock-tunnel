.PHONY: dev build fmt format typecheck test clean docker-build

dev:
	npm run dev

build:
	npm run build

fmt:
	npm run fmt

format:
	npm run format

typecheck:
	npm run typecheck

test:
	npm run test

clean:
	npm run clean

docker-build:
	docker build -t codedock-tunnel:dev .
