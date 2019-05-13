#!/usr/bin/env bash

docker run \
    --name frag \
    --network ca-network_ca \
    -p 8081:8080 \
    -d \
    -v $(pwd)/test/configs/ca-docker:/app/configs \
    -v $(pwd)/test:/app/test \
    frag:dev

docker run \
    --name frag-tls \
    --network ca-network_ca \
    -p 8082:8080 \
    -d \
    -v $(pwd)/test/ca-network/ca-tls:/app/ca-tls \
    -v $(pwd)/test/configs/ca-tls-docker:/app/configs \
    -v $(pwd)/test:/app/test \
    frag:dev
