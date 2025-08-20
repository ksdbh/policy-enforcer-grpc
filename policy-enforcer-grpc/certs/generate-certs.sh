#!/usr/bin/env bash
set -euo pipefail
mkdir -p certs
pushd certs >/dev/null
openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes -keyout ca.key -out ca.pem -subj "/CN=PolicyCA"
openssl req -newkey rsa:4096 -nodes -keyout server.key -out server.csr -subj "/CN=policy-enforcer"
openssl x509 -req -in server.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out server.pem -days 365 -sha256
openssl req -newkey rsa:4096 -nodes -keyout client.key -out client.csr -subj "/CN=policy-client"
openssl x509 -req -in client.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out client.pem -days 365 -sha256
ls -1
popd >/dev/null
