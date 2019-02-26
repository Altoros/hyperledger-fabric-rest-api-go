#!/usr/bin/env bash

if [ -d "test/byfn/fabric-samples/first-network" ]; then
    cd test/byfn/fabric-samples/first-network
    ./byfn.sh down
fi

