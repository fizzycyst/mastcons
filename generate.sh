#!/bin/bash
# Set up the proto buffer..
protoc  amcrpc/amcpb/amc.proto --go_out=plugins=grpc:.