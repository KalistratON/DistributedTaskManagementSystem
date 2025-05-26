#!/bin/bash

export MONGODB_CONNECTION_STRING=mongodb://root:password@localhost:27017; go test -v ./internal/repository/user
export MONGODB_CONNECTION_STRING=mongodb://root:password@localhost:27017; go test -bench=. -v ./internal/repository/user


