# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# モジュールキャッシュを先にコピー (レイヤーキャッシュ最適化)
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /server ./cmd/server

# Runtime stage (distroless: シェルなし・最小攻撃面)
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
