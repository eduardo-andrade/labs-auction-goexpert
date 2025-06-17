FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /auction cmd/auction/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata curl

WORKDIR /app

COPY --from=builder /auction .
COPY .env .

ENV TZ=America/Sao_Paulo

EXPOSE 8080
ENTRYPOINT ["./auction"]