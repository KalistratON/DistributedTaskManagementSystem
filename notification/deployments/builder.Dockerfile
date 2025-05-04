FROM docker.io/golang:1.22 AS builder

WORKDIR /app
COPY . .
RUN go build -mod=readonly -o /bin/notification notification/cmd/main.go

ENTRYPOINT [ "/bin/notification" ]
