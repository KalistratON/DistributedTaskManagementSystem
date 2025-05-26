FROM docker.io/golang:1.22 AS builder

WORKDIR /app
COPY . .
RUN go build -o /bin/task cmd/task/main.go

ENTRYPOINT [ "/bin/task" ]
