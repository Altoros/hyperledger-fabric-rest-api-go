#!/usr/bin/env bash

docker run \
    --name frag \
    --network net_byfn \
    -p 8080:8080 \
    -d \
    -v $(pwd)/test/configs/byfn-docker:/app/configs \
    -v $(pwd)/_tmp/fabric-samples:/app/fabric-samples \
    frag:dev
