package main

import "fmt"
import "github.com/hyperledger/fabric/core/chaincode/shim"

type FarmCC struct {
}

func (t *FarmCC) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("init called")
	return nil, nil
}

func (t *FarmCC) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke called")
	return nil, nil
}

func (t *FarmCC) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query called")
	return nil, nil
}

func main() {
	err := shim.Start(new(FarmCC))
	if err != nil {
		fmt.Printf("Error starting FarmCC : %s", err)
	}
}
