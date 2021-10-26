#!/bin/sh
echo "creating proto file"
protoc --go_out=plugins=grpc,paths=source_relative:../internal/items/. item.proto