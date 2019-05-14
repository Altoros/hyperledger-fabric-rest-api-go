#!/usr/bin/env bash

. ./scripts/fabric-samples.sh

cloneFabricSamples

docker-compose -f ./test/ca-network/docker-compose.yaml up -d
