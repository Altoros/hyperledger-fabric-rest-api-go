# Fabric REST API Go

REST control system for [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric) network.

Features
* Invoke and query chaincode
* Chaincode detailed info
* Channels list
* Channel peers
* ...

## Local test instructions

Requirements:
* https://hyperledger-fabric.readthedocs.io/en/release-1.4/install.html
* https://github.com/hyperledger/fabric-sdk-go
* Go 1.12 

Clone repo outside GOPATH in order to go module work correctly.

Build and run
```
make run
```

*NOTE: only one test network could be up at the same time*

Stop and clear all test networks
```
make clear
```

### Basic network test

Start local basic test network
```
make basic_up
```

To create and install/instantiate chaincode do POST request to 
*localhost:8080/init_test_fixtures*
 

Postman collection

*/test/FabricApiBasic.postman_collection.json* 


### BYFN test

Additional requirements:
* [qj](https://stedolan.github.io/jq/) command-line JSON processor

Run E2E tests against BYFN
```
make byfn
```

Postman collection

*/test/FabricApiBYFN.postman_collection.json*

## Work progress 

### Close plans

* Create Dockerfile
* Test docker container inside testnet 

### Strategic plans

* Automate postman collections tests with Newman
* Test with BYFN (partially done)
* Create full integration test with BYFN, makefile and shell scripts (partially done)
* Cover code with unit test (partially done)
* Move configuration to ENV variables
* Channels creation
* Chaincode installation & instantiation 
* Organisations and users management 

