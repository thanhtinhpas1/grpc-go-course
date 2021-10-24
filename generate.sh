#!/bin/bash

protoc greet/greetpb/greet.proto --go_out=. --go-grpc_out=.

protoc calculator/pb/calculator.proto --go_out=. --go-grpc_out=.
