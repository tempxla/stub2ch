#!/bin/sh
ln -s ./web ./internal/app/web
go test ./...
rm ./internal/app/web
