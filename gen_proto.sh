#!/bin/bash 

protoc \
  -I "$GOPATH/pkg/mod/github.com/cosmos/cosmos-sdk@v0.42.5/proto" \
  -I "$GOPATH/pkg/mod/github.com/cosmos/cosmos-sdk@v0.42.5/third_party/proto" \
  --gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
  --grpc-gateway_out=logtostderr=true:. \
  --proto_path proto \
  ./proto/iscn/* \
  --swagger_out ./swagger-gen \
  --swagger_opt logtostderr=true --swagger_opt fqn_for_swagger_name=true --swagger_opt simple_operation_ids=true

mv github.com/likecoin/likechain/x/iscn/types/* x/iscn/types/
