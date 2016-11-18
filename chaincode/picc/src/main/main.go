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

	createFarmTables(stub)
	// for testing
	populateSampleFarmRows(stub)

	createBeefTable(stub)
	populateSampleBeefRows(stub)

	createInsuranceTable(stub)
	populateSampleInsuranceRows(stub)

	createLoanTables(stub)
	populateSampleLoanRows(stub)
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
	} else if function == "getBeefByFarmAndLabel" {
		if len(args) != 2 {
			return nil, errors.New("args not match for getBeefByFarmAndLabel, need 2 arg as farm id and ear label")
		}
		return getBeefByFarmAndLabel(stub, args)
	} else if function == "getAllInsurancesByFarm" {
		if len(args) != 1 {
			return nil, errors.New("args not match for getAllInsurancesByFarm, need 1 arg as farm id")
		}
		return getAllInsurancesByFarm(stub, args[0])
	} else if function == "getAllLoansByFarm" {
		if len(args) != 1 {
			return nil, errors.New("args not match for getAllLoansByFarm, need 1 arg as farm id")
		}
		return getAllLoansByFarm(stub, args[0])
	} else if function == "getAllFarmIdsByCity" {
		if len(args) != 2 {
			return nil, errors.New("args not match for getAllFarmIdsByCity, need 2 args as province and city")
		}
		return getAllFarmIdsByCity(stub, args)
	} else if function == "getAllFarmIdsByProvince" {
		if len(args) != 1 {
			return nil, errors.New("args not match for getAllFarmIdsByProvince, need 1 arg as province")
		}
		return getAllFarmIdsByProvince(stub, args[0])
	} else if function == "getAllFarmIdsByName" {
		if len(args) != 3 {
			return nil, errors.New("args not match for getAllFarmIdsByName, need 3 args as province, city and farm name")
		}
		return getAllFarmIdsByName(stub, args)
	} else if function == "getAllLoanIdByLender" {
		if len(args) != 1 {
			return nil, errors.New("args not match for getAllLoanIdByLender, need 1 arg as loan officer id")
		}
		return getAllLoanIdByLender(stub, args[0])
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
