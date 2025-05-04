FROM docker.io/golang:1.22 AS builder

WORKDIR /app
COPY . .
RUN go build -o /bin/auth auth/cmd/main.go

ENTRYPOINT [ "/bin/auth" ]
