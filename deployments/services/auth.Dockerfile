FROM docker.io/golang:1.22 AS builder

WORKDIR /app
COPY . .
RUN go build -o /bin/auth cmd/auth/main.go

ENTRYPOINT [ "/bin/auth" ]
