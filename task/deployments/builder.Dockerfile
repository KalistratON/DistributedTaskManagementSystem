FROM docker.io/golang:1.22 AS builder

WORKDIR /app
COPY . .
RUN go build -o /bin/task task/cmd/main.go

ENTRYPOINT [ "/bin/task" ]
