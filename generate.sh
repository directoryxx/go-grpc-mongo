#!/bin/bash
# apt install -y clang-format protobuf-compiler
# go get github.com/golang/protobuf/protoc-gen-go
protoc blog/blogproto/blog.proto --go_out=plugins=grpc:.
