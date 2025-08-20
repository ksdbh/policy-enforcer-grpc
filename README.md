# Policy Enforcer (Go + gRPC + mTLS)
[![CI](https://github.com/ksdbh/policy-enforcer-grpc/actions/workflows/ci.yml/badge.svg)](https://github.com/ksdbh/policy-enforcer-grpc/actions/workflows/ci.yml)

A minimal **policy evaluation** service over **gRPC** with **mutual TLS**.
Includes client example, Dockerfile, Helm chart, CI, and a script to generate dev certificates.

## Features
- `.proto` contract with `Evaluate` RPC
- gRPC server with mTLS; simple client
- Basic rule engine (allow/deny + reasons)
- Dockerfile & Helm chart
- GitHub Actions CI (protoc, build, test)

## Quickstart (local)
```bash
# Generate dev certificates
bash certs/generate-certs.sh

# Build & run server
go build ./server && TLS_CA=certs/ca.pem TLS_CERT=certs/server.pem TLS_KEY=certs/server.key ./server

# In another shell, run client
TLS_CA=certs/ca.pem CLIENT_CERT=certs/client.pem CLIENT_KEY=certs/client.key go run ./client
```

## Helm (kind/minikube)
```bash
helm upgrade --install policy-enforcer charts/policy-enforcer
# Create the TLS secret named 'policy-enforcer-certs' with your PEM files before installing
```

## Docker
```bash
docker build -t policy-enforcer:dev .
docker run -p 50051:50051 -e TLS_CA=/certs/ca.pem -e TLS_CERT=/certs/server.pem -e TLS_KEY=/certs/server.key policy-enforcer:dev
```

## License
This project is licensed under the MIT License — see `LICENSE`.
