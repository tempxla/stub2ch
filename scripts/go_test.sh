#!/bin/sh
template_path=./internal/app/handle/web
if [ -e $template_path ]; then
    echo you cannot run test. path exists.: $template_path
    exit 1
fi
ln -s ./web $template_path
go test ./... -cover
rm $template_path
