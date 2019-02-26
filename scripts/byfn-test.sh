#!/usr/bin/env bash

# // TODO add jq app installation check

QUERY_RESULT_1=$(curl -s -X GET 'localhost:8080/channels/mychannel/chaincodes/mycc/query?fcn=query&args=a' | jq -r '.result')

PASSED=0
FAULTED=0

function assertEqual() {
    echo -n "$3 "
    if [[ $1 == $2 ]]; then
        echo "$(tput setaf 2)Passed$(tput sgr0)"
        let "PASSED++"
    else
        echo "$(tput setaf 1)Faulted!$(tput sgr0)"
        let "FAULTED++"
    fi
}

echo ""
echo "***************************************************************************"
echo "===================== Start API testing against BYFN ======================"
echo "Start API testing against BYFN"

./build/frag -config=./test/byfn/config.json > /dev/null &
REST_PID=$!

echo "Starting REST API, PID = {$REST_PID}"
sleep 2

echo "Test welcome message"
WELCOME_RESULT=$(curl -s -X GET 'localhost:8080')

echo "Test query"
QUERY_RESULT_1=$(curl -s -X GET 'localhost:8080/channels/mychannel/chaincodes/mycc/query?fcn=query&args=a' | jq -r '.result')
QUERY_RESULT_2=$(curl -s -X GET 'localhost:8080/channels/mychannel/chaincodes/mycc/query?fcn=query&args=b' | jq -r '.result')

assertEqual ${QUERY_RESULT_1} "90" "Balance A should be equal to 90"
assertEqual ${QUERY_RESULT_2} "210" "Balance B should be equal to 210"

echo "Invoke chaincode, transferring 5 from A to B"
curl -s -X POST -d 'fcn=invoke&args=a,b,5' 'localhost:8080/channels/mychannel/chaincodes/mycc/invoke' > /dev/null

QUERY_RESULT_1=$(curl -s -X GET 'localhost:8080/channels/mychannel/chaincodes/mycc/query?fcn=query&args=a' | jq -r '.result')
QUERY_RESULT_2=$(curl -s -X GET 'localhost:8080/channels/mychannel/chaincodes/mycc/query?fcn=query&args=b' | jq -r '.result')

assertEqual ${QUERY_RESULT_1} "85" "Balance A should be equal to 85"
assertEqual ${QUERY_RESULT_2} "215" "Balance B should be equal to 215"

echo "Invoke chaincode, transferring 5 back from B to A"
curl -s -X POST -d 'fcn=invoke&args=b,a,5' 'localhost:8080/channels/mychannel/chaincodes/mycc/invoke' > /dev/null

echo "Killing REST API process"
kill $REST_PID

echo "$(tput setaf 2) Tests passed: $PASSED$(tput sgr0)"
if (( FAULTED > 0 )); then
    echo "$(tput setaf 1)Tests faulted: $FAULTED$(tput sgr0)"
else
    echo "$(tput setaf 2)All tests passed!$(tput sgr0)"
fi


echo "===================== Finish API testing against BYFN ====================="
echo "***************************************************************************"
echo ""
