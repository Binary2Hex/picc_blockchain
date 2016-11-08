package main

import "fmt"
import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type MainCC struct {
}

func (t *MainCC) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	ccLogger.Debug("Init called!")
	t.createFarmTable(stub)
	// for testing
	populateSampleFarmRows(stub)

	return nil, nil
}

func (t *MainCC) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	ccLogger.Debug("Invoke called!")
	return nil, nil
}

func (t *MainCC) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	ccLogger.Debug("Query called!")

	if function == "getFarmById" {
		var farm *Farm
		if len(args) != 1 {
			return nil, errors.New("args not match, need 1 arg as farm id")
		}
		farm = t.getFarmById(stub, args[0])
		returnVal, _ := json.Marshal(farm)
		return returnVal, nil
	} else if function == "getFarmAmount" {
		return t.getFarmAmount(stub)
	}

	ccLogger.Debug("function " + function + " not supported!")
	return nil, errors.New("function " + function + " not supported!")
}

func main() {
	err := shim.Start(new(MainCC))
	if err != nil {
		fmt.Printf("Error starting MainCC : %s", err)
	}
}
