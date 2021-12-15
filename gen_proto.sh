#!/bin/bash 

COSMOS_SDK_VERSION="0.42.11"
SWAGGER_OUT="./swagger-gen"

mkdir -p ${SWAGGER_OUT}

protoc \
  -I "$GOPATH/pkg/mod/github.com/cosmos/cosmos-sdk@v${COSMOS_SDK_VERSION}/proto" \
  -I "$GOPATH/pkg/mod/github.com/cosmos/cosmos-sdk@v${COSMOS_SDK_VERSION}/third_party/proto" \
  --gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
  --grpc-gateway_out=logtostderr=true:. \
  --proto_path proto \
  ./proto/iscn/* \
  --swagger_out ${SWAGGER_OUT} \
  --swagger_opt logtostderr=true --swagger_opt fqn_for_swagger_name=true --swagger_opt simple_operation_ids=true

mv github.com/likecoin/likechain/x/iscn/types/* x/iscn/types/
