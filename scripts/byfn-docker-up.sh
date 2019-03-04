#!/usr/bin/env bash

docker run \
    --name frag \
    --network net_byfn \
    -p 8080:8080 \
    -d \
    -v $(pwd)/test/byfn-docker:/app/configs \
    -v $(pwd)/test:/app/test \
    frag:dev
