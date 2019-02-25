#!/usr/bin/env bash

if [ -d "test/fabric-samples/first-network" ]; then
    cd test/fabric-samples/first-network
    ./byfn.sh down
fi

