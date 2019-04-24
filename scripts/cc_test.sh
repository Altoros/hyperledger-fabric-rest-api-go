#!/usr/bin/env bash

CC_ID="test_$(date +%s)"
CHANNEL="mychannel"

export CC_ID
export CHANNEL

./scripts/cc_init.sh

go test -count=1 -v ./test/chaincode/events_test.go
