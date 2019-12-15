#!/bin/sh
web_path=$(pwd)/internal/app/handle/web
target=$(pwd)/web
if [ -e $web_path ]; then
    echo you cannot run test. path exists.: $web_path
    exit 1
fi
ln -s $target $web_path
go test ./... -cover
rm $web_path
