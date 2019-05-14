#!/usr/bin/env bash

docker run \
    --name frag \
    --network net_basic \
    -p 8080:8080 \
    -d \
    -v $(pwd)/test/configs/basic-docker:/app/configs \
    -v $(pwd)/_tmp/fabric-samples:/app/fabric-samples \
    -v $(pwd)/test/chaincode:/app/chaincode \
    frag:dev
