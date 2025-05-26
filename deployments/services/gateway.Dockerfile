FROM docker.io/golang:1.22 AS builder

WORKDIR /app
COPY . .

RUN go build -o /bin/gateway cmd/gateway/server/main.go

ENTRYPOINT [ "/bin/gateway" ]
