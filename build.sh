#!/usr/bin/env bash
GOOS=linux GOARCH=amd64 go build -o main main.go
GOOS=linux GOARCH=amd64 go build -o jobs jobs.go
