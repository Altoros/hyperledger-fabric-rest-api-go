#!/usr/bin/env bash

CC_VERSION="1.0"
CC_PATH="./test/chaincode/events"

tmpName=$(mktemp)
tar -zcvf ${tmpName} -C ${CC_PATH} .
tarballCc=$(base64 -w 0 ${tmpName})

echo Install test chaincode

http POST localhost:8080/chaincodes/install \
    name="${CC_ID}" \
    version="${CC_VERSION}" \
    channel="${CHANNEL}" \
    data="${tarballCc}" \
    peers:='["org1/peer0"]'

echo Instantiate test chaincode

http POST localhost:8080/chaincodes/instantiate \
    name="${CC_ID}" \
    version="${CC_VERSION}" \
    channel="${CHANNEL}" \
    policy="AND('Org1MSP.member')" \
    args:='[]' \
    peers:='["org1/peer0"]'
