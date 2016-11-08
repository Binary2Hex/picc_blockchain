package main

import "github.com/hyperledger/fabric/core/chaincode/shim"
import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/op/go-logging"
	"strconv"
)

const (
	FARM_TABLE    = "farm"
	FARMCC_LOGGER = "FarmCC"
)

var ccLogger = logging.MustGetLogger(FARMCC_LOGGER)

type colDef struct {
	Name string
	Type string
}

var farmColumnTypes = []colDef{
	{"ID", "string"},
	{"BASICINFO", "string"},
	{"CATTLEID", "string"},
	{"INSURANCEID", "string"},
	{"LOANID", "string"},
}

var farmColumnsKeys = []bool{true, false, false, false, false}

func (f *MainCC) createFarmTable(stub *shim.ChaincodeStub) error {
	farmColumns, err := generateColumns(farmColumnTypes, farmColumnsKeys)
	if err != nil {
		ccLogger.Error("Failed to generate columns for the farm table")
		return errors.New("Failed to generate columns for the farm table")
	}

	err = stub.CreateTable(FARM_TABLE, farmColumns)
	if err != nil {
		ccLogger.Error(err)
		//it's ok if the table already existed, recreate it
		if table, err := stub.GetTable(FARM_TABLE); err == nil && table != nil {
			stub.DeleteTable(FARM_TABLE)
			err = stub.CreateTable(FARM_TABLE, farmColumns)
			if err != nil {
				ccLogger.Error(err)
				return errors.New("failed to create farm table")
			}
		} else {
			return errors.New("Unknown error creating farm table")
		}
	}
	ccLogger.Info("Successfully created farm table")
	return nil
}

func (f *MainCC) getFarmById(stub *shim.ChaincodeStub, id string) *Farm {
	var columns []shim.Column
	col := shim.Column{Value: &shim.Column_String_{String_: id}}
	columns = append(columns, col)
	queryOutput, err := stub.GetRow(FARM_TABLE, columns)
	if err != nil {
		ccLogger.Error(err)
		return nil
	} else if len(queryOutput.Columns) == 0 { // no farm found with the id provided
		ccLogger.Debug("No farm found with id:" + id)
		return nil
	}
	farm := formatFarm(queryOutput)
	ccLogger.Debug("farm retrieved...")
	ccLogger.Debug(farm)
	return farm
}

func (f *MainCC) getFarmAmount(stub *shim.ChaincodeStub) ([]byte, error) {
	// 或者用一个key来保存当前农场的数量，目前方式较繁琐
	rowsChan, err := stub.GetRows(FARM_TABLE, []shim.Column{})
	if err != nil {
		return nil, fmt.Errorf("getFarmAmount query failed. %s", err)
	}
	rows := 0
	for {
		select {
		case _, ok := <-rowsChan:
			if !ok {
				rowsChan = nil
			} else {
				rows++
			}
		}
		if rowsChan == nil {
			break
		}
	}
	ccLogger.Debug(strconv.Itoa(rows) + " farms in total")
	return []byte(strconv.Itoa(rows)), nil
}

func formatFarm(queryOutput shim.Row) *Farm {
	farm := &Farm{}
	farm.ID = queryOutput.Columns[0].GetString_()
	basicInfo := &Farm_BasicInfo{}
	err := json.Unmarshal([]byte(queryOutput.Columns[1].GetString_()), &basicInfo)
	if err != nil {
		ccLogger.Info("Cannot un-marshal [%x]", queryOutput)
	}
	farm.BasicInfo = basicInfo

	err = json.Unmarshal([]byte(queryOutput.Columns[2].GetString_()), &farm.CattleId)
	if err != nil {
		ccLogger.Info("Cannot un-marshal [%x]", queryOutput)
	}

	err = json.Unmarshal([]byte(queryOutput.Columns[3].GetString_()), &farm.InsuranceId)
	if err != nil {
		ccLogger.Info("Cannot un-marshal [%x]", queryOutput)
	}

	err = json.Unmarshal([]byte(queryOutput.Columns[4].GetString_()), &farm.LoanId)
	if err != nil {
		ccLogger.Info("Cannot un-marshal [%x]", queryOutput)
	}
	return farm
}

func generateFarmRow(farm *Farm) []*shim.Column {
	var farmVal []string
	var farmColumns []*shim.Column
	farmVal = append(farmVal, farm.ID)
	basicInfo, _ := json.Marshal(farm.BasicInfo)
	cattleId, _ := json.Marshal(farm.CattleId)
	insuranceId, _ := json.Marshal(farm.InsuranceId)
	loanId, _ := json.Marshal(farm.LoanId)
	farmVal = append(farmVal, string(basicInfo))
	farmVal = append(farmVal, string(cattleId))
	farmVal = append(farmVal, string(insuranceId))
	farmVal = append(farmVal, string(loanId))

	for i := 0; i < len(farmVal); i++ {
		farmColumns = append(farmColumns, &shim.Column{Value: &shim.Column_String_{String_: farmVal[i]}})
	}
	return farmColumns
}

func generateColumns(colTypes []colDef, colKeys []bool) ([]*shim.ColumnDefinition, error) {
	var tableColumns []*shim.ColumnDefinition
	for i, k := range colTypes {
		if k.Type == "string" {
			tableColumns = append(tableColumns, &shim.ColumnDefinition{k.Name, shim.ColumnDefinition_STRING, colKeys[i]})
		} else if k.Type == "bytes" {
			tableColumns = append(tableColumns, &shim.ColumnDefinition{k.Name, shim.ColumnDefinition_BYTES, colKeys[i]})
		} else if k.Type == "bool" {
			tableColumns = append(tableColumns, &shim.ColumnDefinition{k.Name, shim.ColumnDefinition_BOOL, colKeys[i]})

		}
	}
	return tableColumns, nil
}
