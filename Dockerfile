FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o main cmd/main.go

FROM alpine:3.13
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .env


EXPOSE 8000
CMD [ "/app/main" ]