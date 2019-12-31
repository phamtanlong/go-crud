#!/usr/bin/env sh

# Generate Go code for this proto file

protoc -I pb/ pb/proto.proto --go_out=plugins=grpc:pb