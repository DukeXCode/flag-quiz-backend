FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache gcc musl-dev sqlite-dev
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
EXPOSE 8080
RUN apk add --no-cache sqlite-libs sqlite
RUN mkdir -p /app/data
COPY migration/ /app/migration/
COPY --from=builder /app/main .
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh
ENTRYPOINT ["entrypoint.sh"]

