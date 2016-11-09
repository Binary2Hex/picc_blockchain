package main

import "fmt"
import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
)

type MainCC struct {
}

const (
	MAINCC_LOGGER = "MAIN_CHAINCODE_LOGGER"
)

var ccLogger = logging.MustGetLogger(MAINCC_LOGGER)

func (t *MainCC) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	ccLogger.Debug("Init called!")

	createFarmTable(stub)
	// for testing
	populateSampleFarmRows(stub)

	createBeefTable(stub)
	populateSampleBeefRows(stub)

	createInsuranceTable(stub)
	populateSampleInsuranceRows(stub)

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
			return nil, errors.New("args not match for getFarmById, need 1 arg as farm id")
		}
		farm = getFarmById(stub, args[0])
		returnVal, _ := json.Marshal(farm)
		return returnVal, nil
	} else if function == "getFarmAmount" {
		return getFarmAmount(stub)
	} else if function == "getAllBeevesByFarm" {
		if len(args) != 1 {
			return nil, errors.New("args not match for getAllBeevesByFarm, need 1 arg as farm id")
		}
		return getAllBeevesByFarm(stub, args[0])
	} else if function == "getAllInsurancesByFarm" {
		if len(args) != 1 {
			return nil, errors.New("args not match for getAllInsurancesByFarm, need 1 arg as farm id")
		}
		return getAllInsurancesByFarm(stub, args[0])
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
