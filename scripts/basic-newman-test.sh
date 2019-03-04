#!/usr/bin/env bash

docker run \
    --rm \
    --network net_basic \
    -v $(pwd)/test:/etc/newman \
    -t postman/newman:alpine \
    run FabricApiBasic.postman_collection.json \
    --env-var host=frag:8080
