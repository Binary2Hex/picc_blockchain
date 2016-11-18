package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

const (
	INSURANCE_TABLE = "insurance_table"
)

var insuranceColumnTypes = []ColDef{
	{"FARM", "string"},
	{"NUMBER", "string"},
	{"START_DATE", "string"},
	{"END_DATE", "string"},
	{"CHECKED", "bool"},
	{"DETAILS", "string"},
	{"AMOUNT_INSURED", "int64"},
	{"AMOUNT_LOSS", "int64"},
}

var insuranceColumnsKeys = []bool{true, true, false, false, false, false, false, false}

func createInsuranceTable(stub *shim.ChaincodeStub) error {
	return createTable(stub, INSURANCE_TABLE, insuranceColumnTypes, insuranceColumnsKeys)
}

func getAllInsurancesByFarm(stub *shim.ChaincodeStub, farmId string) ([]byte, error) {
	rowsChan, err := stub.GetRows(INSURANCE_TABLE, []shim.Column{{Value: &shim.Column_String_{String_: farmId}}})
	if err != nil {
		return nil, fmt.Errorf("getAllInsurancesByFarm query failed. %s", err)
	}
	var insurances []*Insurance
	rows := 0
	for {
		select {
		case row, ok := <-rowsChan:
			if !ok {
				rowsChan = nil
			} else {
				rows++
				insurance := formatInsurance(row)
				insurances = append(insurances, insurance)
			}
		}
		if rowsChan == nil {
			break
		}
	}
	ccLogger.Debug(strconv.Itoa(rows) + " insurance(s) in total for farm:" + farmId)
	returnVal, _ := json.Marshal(insurances)
	return returnVal, nil
}

func generateInsuranceRow(insurance *Insurance) []*shim.Column {
	var insuranceColumns []*shim.Column
	insuranceColumns = append(insuranceColumns, &shim.Column{Value: &shim.Column_String_{String_: insurance.Farm}})
	insuranceColumns = append(insuranceColumns, &shim.Column{Value: &shim.Column_String_{String_: insurance.Number}})
	insuranceColumns = append(insuranceColumns, &shim.Column{Value: &shim.Column_String_{String_: insurance.StartDate}})
	insuranceColumns = append(insuranceColumns, &shim.Column{Value: &shim.Column_String_{String_: insurance.EndDate}})
	insuranceColumns = append(insuranceColumns, &shim.Column{Value: &shim.Column_Bool{Bool: insurance.Checked}})
	insuranceColumns = append(insuranceColumns, &shim.Column{Value: &shim.Column_String_{String_: insurance.Details}})
	insuranceColumns = append(insuranceColumns, &shim.Column{Value: &shim.Column_Int64{Int64: insurance.AmountInsured}})
	insuranceColumns = append(insuranceColumns, &shim.Column{Value: &shim.Column_Int64{Int64: insurance.AmountLoss}})
	return insuranceColumns
}

func formatInsurance(queryOutput shim.Row) *Insurance {
	insurance := new(Insurance)
	insurance.Farm = queryOutput.Columns[0].GetString_()
	insurance.Number = queryOutput.Columns[1].GetString_()
	insurance.StartDate = queryOutput.Columns[2].GetString_()
	insurance.EndDate = queryOutput.Columns[3].GetString_()
	insurance.Checked = queryOutput.Columns[4].GetBool()
	insurance.Details = queryOutput.Columns[5].GetString_()
	insurance.AmountInsured = queryOutput.Columns[6].GetInt64()
	insurance.AmountLoss = queryOutput.Columns[7].GetInt64()
	return insurance
}

func populateSampleInsuranceRows(stub *shim.ChaincodeStub) {
	insurance := Insurance{}
	insurance.Farm = "1234567"
	insurance.Number = "TKPX12KD1SKS"
	insurance.StartDate = "2016-03-12"
	insurance.EndDate = "2016-09-12"
	insurance.Checked = true
	insurance.Details = "all checked"
	insurance.AmountInsured = 500
	insurance.AmountLoss = 0
	stub.InsertRow(INSURANCE_TABLE, shim.Row{Columns: generateInsuranceRow(&insurance)})

	insurance.Number = "PXAL12L4KS2K"
	stub.InsertRow(INSURANCE_TABLE, shim.Row{Columns: generateInsuranceRow(&insurance)})

	insurance.Farm = "1234568"
	insurance.Number = "XK2SKS91AS5K"
	insurance.StartDate = "2016-05-23"
	insurance.EndDate = "2016-11-23"
	insurance.AmountInsured = 700
	insurance.AmountLoss = 8
	stub.InsertRow(INSURANCE_TABLE, shim.Row{Columns: generateInsuranceRow(&insurance)})

}
