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

	//new a farm and insert for testing...
	farm := new(Farm)
	farm.ID = "1234567"
	basicInfo := new(Farm_BasicInfo)
	basicInfo.Addr = "BEIJING"
	basicInfo.Owner = "ALICE"
	basicInfo.Area = "120"
	basicInfo.Quantity = "2000"
	basicInfo.Species = "cattle"
	farm.BasicInfo = basicInfo
	farm.CattleId = []string{"1234567_eartag0", "1234567_eartag1"}
	farm.InsuranceId = []string{"insurance_id_xx0", "insurance_id_xx1"}
	farm.LoanId = []string{"loan_id_xx0", "loan_id_xx1"}
	stub.InsertRow(FARM_TABLE, shim.Row{Columns: generateFarmRow(farm)})
	ccLogger.Debug("a new farm object inserted in farm table")
	farm.ID = "1234568"
	stub.InsertRow(FARM_TABLE, shim.Row{Columns: generateFarmRow(farm)})
	ccLogger.Debug("another new farm object inserted in farm table")

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
