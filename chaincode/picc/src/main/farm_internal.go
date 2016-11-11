package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

const (
	FARM_TABLE                = "farm_table"
	FARM_LOCATION_INDEX_TABLE = "farm_location_index_table"
)

var farmColumnTypes = []ColDef{
	{"ID", "string"},
	{"BASICINFO", "string"},
	{"FUNDINGINFO", "string"},
	{"INVENTORY", "string"},
	{"FEED", "string"},
}

var farmLocationIndexColumnTypes = []ColDef{
	{"PROVINCE", "string"},
	{"CITY", "string"},
	{"FARM_NAME", "string"}, //同一城市中的养殖场不能重名
	{"FARM_ID", "string"},
}

var farmColumnsKeys = []bool{true, false, false, false, false}
var farmLocationIndexColumnKeys = []bool{true, true, true, false}

func createFarmTables(stub *shim.ChaincodeStub) error {
	if err := createTable(stub, FARM_TABLE, farmColumnTypes, farmColumnsKeys); err != nil {
		return err
	}
	return createTable(stub, FARM_LOCATION_INDEX_TABLE, farmLocationIndexColumnTypes, farmLocationIndexColumnKeys)
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

func getAllFarmIdsByCity(stub *shim.ChaincodeStub, location []string) ([]byte, error) {
	if len(location) != 2 {
		return nil, errors.New("args length mismatch in getAllFarmsIdsByCity")
	}
	columns := []shim.Column{}
	col0 := shim.Column{Value: &shim.Column_String_{String_: location[0]}}
	col1 := shim.Column{Value: &shim.Column_String_{String_: location[1]}}
	columns = append(columns, col0)
	columns = append(columns, col1)

	rowsChan, err := stub.GetRows(FARM_LOCATION_INDEX_TABLE, columns)
	if err != nil {
		ccLogger.Error(err)
		return nil, err
	}
	rows := 0
	returnStr := []string{}
	for {
		select {
		case row, ok := <-rowsChan:
			if !ok {
				rowsChan = nil
			} else {
				rows++
				returnStr = append(returnStr, row.Columns[3].GetString_())
			}
		}
		if rowsChan == nil {
			break
		}
	}
	ccLogger.Debug(strconv.Itoa(rows) + " farms in total in " + location[0] + ", " + location[1])
	return json.Marshal(returnStr)
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

func generateFarmLocationIndexRow(province, city, name, id string) []*shim.Column {
	var farmLocationIndexVal []string
	var farmLocationIndexColumns []*shim.Column
	farmLocationIndexVal = append(farmLocationIndexVal, province)
	farmLocationIndexVal = append(farmLocationIndexVal, city)
	farmLocationIndexVal = append(farmLocationIndexVal, name)
	farmLocationIndexVal = append(farmLocationIndexVal, id)

	for i := 0; i < len(farmLocationIndexVal); i++ {
		farmLocationIndexColumns = append(farmLocationIndexColumns, &shim.Column{Value: &shim.Column_String_{String_: farmLocationIndexVal[i]}})
	}
	return farmLocationIndexColumns
}

func populateSampleFarmRows(stub *shim.ChaincodeStub) {
	//new a farm and insert for testing...
	farm := new(Farm)
	farm.ID = "1234567"
	basicInfo := new(Farm_BasicInfo)
	basicInfo.Name = "承德第一肉牛养殖场"
	basicInfo.Province = "HEBEI"
	basicInfo.City = "CHENGDE"
	basicInfo.Addr = "xxx street ###, GPS: {41.231, 117.234}"
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
	inserted, ok = stub.InsertRow(FARM_LOCATION_INDEX_TABLE, shim.Row{Columns: generateFarmLocationIndexRow(farm.BasicInfo.Province, farm.BasicInfo.City, farm.BasicInfo.Name, farm.ID)})
	if inserted {
		ccLogger.Debug("a new row inserted in farm location index table")
	} else if ok != nil {
		ccLogger.Error("error inserting new row in farm location index table")
	}

	farm.ID = "1234568"
	farm.BasicInfo.Name = "承德第二肉牛养殖场"
	inserted, ok = stub.InsertRow(FARM_TABLE, shim.Row{Columns: generateFarmRow(farm)})
	if inserted {
		ccLogger.Debug("another new farm object inserted in farm table")
	} else if ok != nil {
		ccLogger.Error("error inserting new row in farm table")
	}
	inserted, ok = stub.InsertRow(FARM_LOCATION_INDEX_TABLE, shim.Row{Columns: generateFarmLocationIndexRow(farm.BasicInfo.Province, farm.BasicInfo.City, farm.BasicInfo.Name, farm.ID)})
	if inserted {
		ccLogger.Debug("a new row inserted in farm location index table")
	} else if ok != nil {
		ccLogger.Error("error inserting new row in farm location index table")
	}
}
