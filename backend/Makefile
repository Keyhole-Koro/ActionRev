BIN := bin/server
MODULE := github.com/keyhole-koro/actionrev

.PHONY: generate build test lint run clean

## proto → Go コード生成 (buf が必要)
generate:
	buf generate

## ビルド
build:
	go build -o $(BIN) ./cmd/server

## テスト
test:
	go test ./... -race -count=1

## lint (golangci-lint が必要)
lint:
	golangci-lint run ./...

## ローカル起動 (docker compose 経由)
run:
	docker compose -f ../infra/docker-compose.yml up --build backend

## 生成物・バイナリ削除
clean:
	rm -rf bin/ gen/
