#!/usr/bin/env bash

docker run --name frag --network net_basic -p 8080:8080 -d \
    -v $(pwd)/test/basic-docker:/app/configs \
    -v $(pwd)/test:/app/test \
    frag:dev
