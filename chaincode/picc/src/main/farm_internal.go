package main

import "github.com/hyperledger/fabric/core/chaincode/shim"
import (
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	FARM_TABLE = "farm_table"
)

var farmColumnTypes = []ColDef{
	{"ID", "string"},
	{"BASICINFO", "string"},
	{"FUNDINGINFO", "string"},
	{"INVENTORY", "string"},
	{"FEED", "string"},
}

var farmColumnsKeys = []bool{true, false, false, false, false}

func createFarmTable(stub *shim.ChaincodeStub) error {
	return createTable(stub, FARM_TABLE, farmColumnTypes, farmColumnsKeys)
}

func getFarmById(stub *shim.ChaincodeStub, id string) *Farm {
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

func getFarmAmount(stub *shim.ChaincodeStub) ([]byte, error) {
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

	err = json.Unmarshal([]byte(queryOutput.Columns[2].GetString_()), &farm.FundingInfo)
	if err != nil {
		ccLogger.Info("Cannot un-marshal [%x]", queryOutput)
	}

	err = json.Unmarshal([]byte(queryOutput.Columns[3].GetString_()), &farm.Inventory)
	if err != nil {
		ccLogger.Info("Cannot un-marshal [%x]", queryOutput)
	}

	err = json.Unmarshal([]byte(queryOutput.Columns[4].GetString_()), &farm.Feed)
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
	fundingInfo, _ := json.Marshal(farm.FundingInfo)
	inventory, _ := json.Marshal(farm.Inventory)
	feed, _ := json.Marshal(farm.Feed)
	farmVal = append(farmVal, string(basicInfo))
	farmVal = append(farmVal, string(fundingInfo))
	farmVal = append(farmVal, string(inventory))
	farmVal = append(farmVal, string(feed))

	for i := 0; i < len(farmVal); i++ {
		farmColumns = append(farmColumns, &shim.Column{Value: &shim.Column_String_{String_: farmVal[i]}})
	}
	return farmColumns
}

func populateSampleFarmRows(stub *shim.ChaincodeStub) {
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

	fundingInfo := new(Farm_FundingInfo)
	fundingInfo.Outlay = "100"
	fundingInfo.PaidIn = "1000"
	fundingInfo.PovertyRelief = "12"
	farm.FundingInfo = fundingInfo

	inventory := new(Farm_Inventory)
	inventory.AboveOne = 120
	inventory.UnderOne = 230
	inventory.Born = 12
	inventory.Butchery = 40
	inventory.Dead = 2
	inventory.Import = 40
	inventory.Init = 88
	inventory.Insurance = 45
	inventory.Year = 2016
	inventory.Sell = 212
	farm.Inventory = append(farm.Inventory, inventory)

	feed2016 := new(Farm_Feed)
	feed2016.Year = 2016
	feed2016.Type1 = 120
	feed2016.Type2 = 200
	farm.Feed = append(farm.Feed, feed2016)
	feed2015 := new(Farm_Feed)
	feed2015.Year = 2015
	feed2015.Type1 = 130
	feed2015.Type2 = 210
	farm.Feed = append(farm.Feed, feed2015)

	inserted, ok := stub.InsertRow(FARM_TABLE, shim.Row{Columns: generateFarmRow(farm)})
	if inserted {
		ccLogger.Debug("a new farm object inserted in farm table")
	} else if ok != nil {
		ccLogger.Error("error inserting new row in farm table")
	}
	farm.ID = "1234568"
	inserted, ok = stub.InsertRow(FARM_TABLE, shim.Row{Columns: generateFarmRow(farm)})
	if inserted {
		ccLogger.Debug("another new farm object inserted in farm table")
	} else if ok != nil {
		ccLogger.Error("error inserting new row in farm table")
	}
}
