#!/bin/sh

openssl dgst \
        -sha256 \
        -sign ~/.openssl/stub2ch_private.pem \
        tools/auth/passphrase.txt \
    | base64 -w 0
