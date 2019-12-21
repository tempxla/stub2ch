#!/bin/sh

##### Setup #####

# html template path
web_path=$(pwd)/internal/app/handle/web
target=$(pwd)/web
if [ -e $web_path ]; then
    echo you cannot run test. path exists.: $web_path
    exit 1
fi
ln -s $target $web_path

# signature
signature_path=/tmp/sig_stub2ch.txt
openssl dgst \
        -sha256 \
        -sign ~/.openssl/stub2ch_private.pem \
        tools/auth/passphrase.txt \
    | base64 -w 0 > $signature_path

# pass
pass_path=/tmp/pass_stub2ch.txt
cp tools/auth/passphrase.txt $pass_path

##### TESTING #####
cover_out=/tmp/stub2ch_cover.out
go test ./... -coverprofile=$cover_out
if [ $# -eq 1 ]; then
    if [ $1 = "html" ]; then
        go tool cover -html=$cover_out
    elif [ $1 = "func" ]; then
        go tool cover -func=$cover_out
    fi
fi

##### TearDown #####

# html template path
rm $web_path

# signature
rm $signature_path

# pass
rm $pass_path
