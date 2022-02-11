#!/bin/bash 

# get protoc executions
go install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos
export PATH="$PATH:$(go env GOPATH)/bin"

pushd $(dirname $0) > /dev/null 2>&1

SWAGGER_DIR="swagger-gen"
COSMOS_SDK_DIR=$(go list -f '{{ .Dir }}' -m github.com/cosmos/cosmos-sdk)

pushd x > /dev/null 2>&1
MODULES=(*)
popd > /dev/null 2>&1

mkdir -p ${SWAGGER_DIR}

for module in "${MODULES[@]}"; do
    protoc \
      -I "$COSMOS_SDK_DIR/proto" \
      -I "$COSMOS_SDK_DIR/third_party/proto" \
      --gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
      --grpc-gateway_out=logtostderr=true:. \
      --proto_path proto \
      ./proto/${module}/* \
      --swagger_out ${SWAGGER_DIR} \
      --swagger_opt logtostderr=true --swagger_opt fqn_for_swagger_name=true --swagger_opt simple_operation_ids=true

    mv github.com/likecoin/likechain/x/${module}/types/* x/${module}/types/
done

popd > /dev/null 2>&1