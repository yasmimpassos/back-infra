FROM golang:1.25 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o backend

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /src/backend .

EXPOSE 8080

CMD ["./backend"]