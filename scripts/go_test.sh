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

##### TESTING #####
go test ./... -cover

##### TearDown #####

# html template path
rm $web_path

# signature
rm $signature_path
