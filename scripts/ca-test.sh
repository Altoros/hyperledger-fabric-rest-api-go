#!/usr/bin/env bash

docker run \
    --rm \
    --network ca-network_ca \
    -v $(pwd)/test:/etc/newman \
    -t postman/newman:alpine \
    run CA.postman_collection.json \
    --env-var host=frag:8080

docker run \
    --rm \
    --network ca-network_ca \
    -v $(pwd)/test:/etc/newman \
    -t postman/newman:alpine \
    run CA.postman_collection.json \
    --env-var host=frag-tls:8080
