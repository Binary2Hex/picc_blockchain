package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

const (
	BEEF_TABLE = "beef_table"
)

var beefColumnTypes = []ColDef{
	{"FARM", "string"},
	{"EARLABEL", "string"},
	{"FARMER", "string"},
	{"SUBSIDY", "int64"},
	{"INVEST_FROM_FARMER", "int64"},
	{"INVEST_FROM_FARM", "int64"},
	{"BIRTHDAY", "string"},
	{"STATE", "string"},
	{"INSURANCE_STATE", "string"},
	{"CHECKED", "bool"},
	{"TRACE", "bytes"},
}

var beefColumnsKeys = []bool{true, true, false, false, false, false, false, false, false, false, false}

func createBeefTable(stub *shim.ChaincodeStub) error {
	return createTable(stub, BEEF_TABLE, beefColumnTypes, beefColumnsKeys)
}

func getAllBeevesByFarm(stub *shim.ChaincodeStub, farmId string) ([]byte, error) {
	rowsChan, err := stub.GetRows(BEEF_TABLE, []shim.Column{{Value: &shim.Column_String_{String_: farmId}}})
	if err != nil {
		return nil, fmt.Errorf("getAllBeevesByFarm query failed. %s", err)
	}
	var beeves []*Beef
	rows := 0
	for {
		select {
		case row, ok := <-rowsChan:
			if !ok {
				rowsChan = nil
			} else {
				rows++
				beef := formatBeef(row)
				beeves = append(beeves, beef)
			}
		}
		if rowsChan == nil {
			break
		}
	}
	ccLogger.Debug(strconv.Itoa(rows) + " beefs in total for farm:" + farmId)
	returnVal, _ := json.Marshal(beeves)
	return returnVal, nil
}

func generateBeefRow(beef *Beef) []*shim.Column {

	var beefColumns []*shim.Column
	beefColumns = append(beefColumns, &shim.Column{Value: &shim.Column_String_{String_: beef.Farm}})
	beefColumns = append(beefColumns, &shim.Column{Value: &shim.Column_String_{String_: beef.EarLabel}})
	beefColumns = append(beefColumns, &shim.Column{Value: &shim.Column_String_{String_: beef.Farmer}})
	beefColumns = append(beefColumns, &shim.Column{Value: &shim.Column_Int64{Int64: beef.Subsidy}})
	beefColumns = append(beefColumns, &shim.Column{Value: &shim.Column_Int64{Int64: beef.InvestFromFarmer}})
	beefColumns = append(beefColumns, &shim.Column{Value: &shim.Column_Int64{Int64: beef.InvestFromFarm}})
	beefColumns = append(beefColumns, &shim.Column{Value: &shim.Column_String_{String_: beef.Birthday}})
	beefColumns = append(beefColumns, &shim.Column{Value: &shim.Column_String_{String_: beef.State}})
	beefColumns = append(beefColumns, &shim.Column{Value: &shim.Column_String_{String_: beef.InsuranceState}})
	beefColumns = append(beefColumns, &shim.Column{Value: &shim.Column_Bool{Bool: beef.Checked}})
	traceMarshal, _ := json.Marshal(beef.Trace)
	beefColumns = append(beefColumns, &shim.Column{Value: &shim.Column_Bytes{Bytes: traceMarshal}})

	return beefColumns
}

func formatBeef(queryOutput shim.Row) *Beef {
	beef := new(Beef)
	beef.Farm = queryOutput.Columns[0].GetString_()
	beef.EarLabel = queryOutput.Columns[1].GetString_()
	beef.Farmer = queryOutput.Columns[2].GetString_()
	beef.Subsidy = queryOutput.Columns[3].GetInt64()
	beef.InvestFromFarmer = queryOutput.Columns[4].GetInt64()
	beef.InvestFromFarm = queryOutput.Columns[5].GetInt64()
	beef.Birthday = queryOutput.Columns[6].GetString_()
	beef.State = queryOutput.Columns[7].GetString_()
	beef.InsuranceState = queryOutput.Columns[8].GetString_()
	beef.Checked = queryOutput.Columns[9].GetBool()
	err := json.Unmarshal(queryOutput.Columns[10].GetBytes(), &beef.Trace)
	if err != nil {
		ccLogger.Errorf("Cannot un-marshal [%x]", queryOutput)
	}

	return beef
}

func populateSampleBeefRows(stub *shim.ChaincodeStub) {
	beef1 := Beef{}
	beef1.Birthday = "2015-02-11"
	beef1.Checked = true
	beef1.EarLabel = "Z5TC923U81"
	beef1.Farm = "1234567"
	beef1.Farmer = "ZHANGSAN"
	beef1.InsuranceState = "insured"
	beef1.InvestFromFarm = 2000
	beef1.InvestFromFarmer = 2500
	beef1.Subsidy = 2000
	beef1.State = "in"

	trace10 := new(Beef_Trace)
	trace10.Date = "2016-02-15 09:10:23"
	trace10.Event = "imported from Finland"
	trace11 := new(Beef_Trace)
	trace11.Date = "2016-02-16 13:11:54"
	trace11.Event = "transfered to farm"

	beef1.Trace = append(beef1.Trace, trace10)
	beef1.Trace = append(beef1.Trace, trace11)

	beef2 := beef1
	beef2.EarLabel = "CB29AL2D91"
	beef2.Farm = "1234568"
	beef2.Farmer = "LISI"

	inserted, ok := stub.InsertRow(BEEF_TABLE, shim.Row{Columns: generateBeefRow(&beef1)})
	if inserted {
		ccLogger.Debug("a new beef object inserted in beef table")
	} else if ok != nil {
		ccLogger.Error("error inserting new row in beef table")
		ccLogger.Error(ok)
	}

	inserted, ok = stub.InsertRow(BEEF_TABLE, shim.Row{Columns: generateBeefRow(&beef2)})
	if inserted {
		ccLogger.Debug("another new beef object inserted in beef table")
	} else if ok != nil {
		ccLogger.Error("error inserting new row in beef table")
		ccLogger.Error(ok)
	}

}
