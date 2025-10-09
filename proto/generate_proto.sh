#!/bin/bash

cd $(dirname $0)

PROTO_FILE="messages.proto"
OUTPUT_DIR="."

# Generate Go code.
protoc --go_out=$OUTPUT_DIR --go_opt=paths=source_relative \
    $PROTO_FILE

echo "Protobuf code generated successfully."
