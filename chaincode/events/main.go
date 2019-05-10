package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type TestChaincode struct {
}

func (t *TestChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *TestChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// Get the function and arguments from the request
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "setA":
		err := stub.PutState("A", []byte(args[0]))
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte("setOk"))

	case "getA":
		a, err := stub.GetState("A")
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(a)

	case "shimPayload":
		return shim.Success([]byte("test payload"))
	case "shimError":
		return shim.Error("test error")
	case "event":
		err := stub.PutState("hello", []byte("event"))
		if err != nil {
			return shim.Error(err.Error())
		}

		err = stub.SetEvent("testEvent", []byte("event payload"))
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte("event emitted"))
	}

	return shim.Error("Unknown action")
}

func main() {
	err := shim.Start(new(TestChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
