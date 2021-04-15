#!/bin/bash 

protoc \
  -I "$GOPATH/pkg/mod/github.com/cosmos/cosmos-sdk@v0.42.4/proto" \
  -I "$GOPATH/pkg/mod/github.com/cosmos/cosmos-sdk@v0.42.4/third_party/proto" \
  --gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
  --grpc-gateway_out=logtostderr=true:. \
  --proto_path proto \
  ./proto/iscn/*

mv github.com/likecoin/likechain/x/iscn/types/*  x/iscn/types/
