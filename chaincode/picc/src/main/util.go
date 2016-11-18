package main

import (
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type ColDef struct {
	Name string
	Type string
}

func generateColumns(colTypes []ColDef, colKeys []bool) ([]*shim.ColumnDefinition, error) {
	var tableColumns []*shim.ColumnDefinition
	for i, k := range colTypes {
		if k.Type == "string" {
			tableColumns = append(tableColumns, &shim.ColumnDefinition{k.Name, shim.ColumnDefinition_STRING, colKeys[i]})
		} else if k.Type == "bytes" {
			tableColumns = append(tableColumns, &shim.ColumnDefinition{k.Name, shim.ColumnDefinition_BYTES, colKeys[i]})
		} else if k.Type == "bool" {
			tableColumns = append(tableColumns, &shim.ColumnDefinition{k.Name, shim.ColumnDefinition_BOOL, colKeys[i]})
		} else if k.Type == "int32" {
			tableColumns = append(tableColumns, &shim.ColumnDefinition{k.Name, shim.ColumnDefinition_INT32, colKeys[i]})
		} else if k.Type == "int64" {
			tableColumns = append(tableColumns, &shim.ColumnDefinition{k.Name, shim.ColumnDefinition_INT64, colKeys[i]})
		} else {
			ccLogger.Error("column type:" + k.Type + " NOT supported!")
		}
	}
	return tableColumns, nil
}

func createTable(stub *shim.ChaincodeStub, tableName string, columnTypes []ColDef, columnKeys []bool) error {
	if stub == nil || tableName == "" || columnTypes == nil || columnKeys == nil {
		ccLogger.Error("none of the parameters for createTable should be nil")
		return errors.New("none of the parameters for createTable should be nil")
	}

	if len(columnTypes) != len(columnKeys) {
		ccLogger.Error("length of columnTypes and columnKeys should be equal!")
		return errors.New("length of columnTypes and columnKeys should be equal!")
	}

	columns, err := generateColumns(columnTypes, columnKeys)
	if err != nil {
		ccLogger.Error("Failed to generate columns for table:" + tableName)
		return errors.New("Failed to generate columns for table:" + tableName)
	}

	err = stub.CreateTable(tableName, columns)
	if err != nil {
		ccLogger.Error(err)
		//it's ok if the table already existed, recreate it
		if table, err := stub.GetTable(tableName); err == nil && table != nil {
			stub.DeleteTable(tableName)
			err = stub.CreateTable(tableName, columns)
			if err != nil {
				ccLogger.Error(err)
				return errors.New("failed to create table:" + tableName)
			}
		} else {
			return errors.New("Unknown error creating beef table")
		}
	}
	ccLogger.Info("Successfully created table:" + tableName)
	return nil
}
