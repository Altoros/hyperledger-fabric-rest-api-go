#!/usr/bin/env bash

docker run \
    --name frag \
    --network ca-network_ca \
    -p 8080:8080 \
    -d \
    -v $(pwd)/test/ca-docker:/app/configs \
    -v $(pwd)/test:/app/test \
    frag:dev
