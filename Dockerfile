# 1단계: 빌드
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 2단계: 실행
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/set.sql .
# 필요하다면 static 파일 등 추가 복사

EXPOSE 8080

CMD ["./main"]