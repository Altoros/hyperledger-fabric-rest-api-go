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

Run all unit and end-2-end tests

```
make test
```

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
 
Run E2E tests against Basic network
```
make basic_e2e_test
```

Postman collection

*/test/FabricApiBasic.postman_collection.json* 


### BYFN test

Additional requirements:
* [qj](https://stedolan.github.io/jq/) command-line JSON processor

Start local basic test network
```
make byfn_up
```

Run E2E tests against BYFN
```
make byfn_e2e_test
```

Postman collection

*/test/FabricApiBYFN.postman_collection.json*

## Work progress 

### Close plans

* Channels creation
* Chaincode installation & instantiation 

### Strategic plans

* Organisations and users management 
* More tests with BYFN
* Cover code with unit test (partially done)
* Move configuration to ENV variables

