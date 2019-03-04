#!/usr/bin/env bash

docker run \
    --rm \
    --network net_byfn \
    -v $(pwd)/test:/etc/newman \
    -t postman/newman:alpine \
    run FabricApiBYFN.postman_collection.json \
    --env-var host=frag:8080
