#!/bin/bash
protoc blog/blogproto/blog.proto --go_out=plugins=grpc:.
