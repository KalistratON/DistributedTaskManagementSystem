FROM docker.io/golang:1.22 AS builder

WORKDIR /app
COPY . .
RUN go build -o /bin/notification cmd/notification/server/main.go

ENTRYPOINT [ "/bin/notification" ]
