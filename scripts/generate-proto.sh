#!/usr/bin/env bash
set -euo pipefail

echo "Generating protobuf Go code into gen/proto/..."

# Requires: protoc, protoc-gen-go, protoc-gen-go-grpc
# Install planners locally:
#   go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0
#   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
OUT_DIR="$ROOT_DIR/gen/proto"

mkdir -p "$OUT_DIR"

protoc \
	--proto_path="$ROOT_DIR/proto" \
	--go_out="$OUT_DIR" --go_opt=paths=source_relative \
	--go-grpc_out="$OUT_DIR" --go-grpc_opt=paths=source_relative \
	$(find "$ROOT_DIR/proto" -name "*.proto") || true

echo "proto generation finished (if protoc and plugins are installed)"
