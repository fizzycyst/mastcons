#!/bin/bash
# Set up the proto buffer..
protoc amcpb/amc.proto --go_out=plugins=grpc:.