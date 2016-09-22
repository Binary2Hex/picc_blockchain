/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"

	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//==============================================================================================================================
//	 参与者类型映射，参考membersrvc.yaml定义
//==============================================================================================================================

const GOV = 1
const FARM = 2
const FARMER = 3
const BANK = 4
const INSURANCE = 5

//==============================================================================================================================
//	 状态类型 - 如果业务逻辑需要，在此处定义状态类型
//==============================================================================================================================
const STATE_INITIALIZED = 0

//==============================================================================================================================
//	 网络中的一头/批肉牛
//==============================================================================================================================
type Cattle struct {
	Id          string `json:"id"`
	Vaccinated  string `json:"vaccinated"`
	InsuranceID string `json:"insuranceID"`
	Loan        int    `json:"loan"`
	LoanID      string `json:"loanID"`
	Origin      string `json:"origin"`
	Trader      string `json:"trader"`
	Status      int    `json:"status"`
}

//==============================================================================================================================
//	 记录网络中所有肉牛的id，可以用来查询所有的肉牛
//==============================================================================================================================
type CattleSet struct {
	Ids []string `json:"ids"`
}

//==============================================================================================================================
//	User_and_eCert - Struct for storing the JSON of a user and their ecert
//==============================================================================================================================

type User_and_eCert struct {
	Identity string `json:"identity"`
	eCert    string `json:"ecert"`
}
type ECertResponse struct {
	OK string
}

//==============================================================================================================================
//	 常量定义
//==============================================================================================================================
const CATTLE_TABLE_NAME = "beefCattleTable"

// SimpleChaincode example implementation
type CattleChaincode struct {
}

func (t *CattleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	columnDefinitios := []*shim.ColumnDefinition{
		{Name: "Id", Type: shim.ColumnDefinition_STRING, Key: true},
		{Name: "Vaccinated", Type: shim.ColumnDefinition_STRING, Key: false},
		{Name: "InsuranceID", Type: shim.ColumnDefinition_STRING, Key: true},
		{Name: "Loan", Type: shim.ColumnDefinition_INT64, Key: false},
		{Name: "LoanID", Type: shim.ColumnDefinition_STRING, Key: true},
		{Name: "Origin", Type: shim.ColumnDefinition_STRING, Key: true},
		{Name: "Trader", Type: shim.ColumnDefinition_STRING, Key: true},
		{Name: "Status", Type: shim.ColumnDefinition_INT32, Key: false},
	}
	table, err := stub.GetTable(CATTLE_TABLE_NAME)
	if table != nil {
		fmt.Printf("table already exists:%s\n", CATTLE_TABLE_NAME)
		return []byte("Table already exists"), nil
	}
	err = stub.CreateTable(CATTLE_TABLE_NAME, columnDefinitios)
	if err != nil {
		return nil, err
	}

	// TODO, THIS IS JUST FOR TESTING
	t.insertRow(stub, "0000000000", "Vaccinated@2016", "PICC_INSR_789", 100, "PICC_LOAN_123", "Switzerland", "TraderA", 0)
	t.insertRow(stub, "0000000001", "Vaccinated@2015", "PICC_INSR_711", 80, "", "Switzerland", "TraderA", 0)

	return nil, nil
}
func (t *CattleChaincode) insertRow(stub *shim.ChaincodeStub, id string, vaccinated string, insuranceID string, loan int64, loanID string, origin string, trader string, status int32) (bool, error) {
	row := shim.Row{Columns: []*shim.Column{
		{Value: &shim.Column_String_{String_: id}},
		{Value: &shim.Column_String_{String_: vaccinated}},
		{Value: &shim.Column_String_{String_: insuranceID}},
		{Value: &shim.Column_Int64{Int64: loan}},
		{Value: &shim.Column_String_{String_: loanID}},
		{Value: &shim.Column_String_{String_: origin}},
		{Value: &shim.Column_String_{String_: trader}},
		{Value: &shim.Column_Int32{Int32: status}},
	}}
	return stub.InsertRow(CATTLE_TABLE_NAME, row)
}

func (t *CattleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *CattleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	username, affiliation, affiliationRole, err := t.getCaller(stub)
	if err != nil {
		return nil, err
	}
	fmt.Printf("username: %s\n", username)
	fmt.Printf("affiliation: %s\n", affiliation)
	fmt.Printf("affiliationRole: %d\n", affiliationRole)
	if affiliation != "picc_poc" {
		return nil, errors.New("User '" + username + "' doesn't belong to picc group!!!")
	}

	if function == "getAllCattles" {
		return t.getAllCattles(stub)
	}

	return []byte(username), nil
}

/*******************************General Functions**************************/
//==============================================================================================================================
//	 get_ecert - Takes the name passed and calls out to the REST API for HyperLedger to retrieve the ecert
//				 for that user. Returns the ecert as retrived including html encoding.
//==============================================================================================================================
func (t *CattleChaincode) getEcert(stub *shim.ChaincodeStub, name string) (*x509.Certificate, error) {

	// TODO: Hardcoded!! to be removed!!!
	res, err := http.Get("http://localhost:5000/registrar/" + name + "/ecert")
	if err != nil {
		return nil, errors.New("Couldn't get ecert by name:" + name)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("Couldn't read ecert from response body")
	}
	var ecertResponse ECertResponse
	json.Unmarshal(body, &ecertResponse)
	cert := ecertResponse.OK
	if cert == "" {
		return nil, errors.New("Couldn't get ecert content from response body")
	}
	decodedCert, err := url.QueryUnescape(cert) // make % etc normal //
	if err != nil {
		return nil, errors.New("Could not decode certificate")
	}
	pem, _ := pem.Decode([]byte(decodedCert)) // Make Plain text   //
	x509Cert, err := x509.ParseCertificate(pem.Bytes)

	if err != nil {
		return nil, errors.New("Couldn't parse ecert for user " + name)
	}

	return x509Cert, nil
}

//==============================================================================================================================
//	 get_caller - Retrieves the username anf affiliation of the user who invoked the chaincode.
//				  Returns the username as a string, affiliation as an int
//==============================================================================================================================

func (t *CattleChaincode) getCaller(stub *shim.ChaincodeStub) (string, string, int, error) {
	bytes, err := stub.GetCallerCertificate()
	if err != nil {
		return "", "", -1, errors.New("Couldn't retrieve caller certificate")
	}
	x509Cert, err := x509.ParseCertificate(bytes) // Extract Certificate from result of GetCallerCertificate
	if err != nil {
		return "", "", -1, errors.New("Couldn't parse certificate")
	}
	commanName := x509Cert.Subject.CommonName
	x509Cert, err = t.getEcert(stub, commanName)
	if err != nil {
		return "", "", -1, errors.New("Couldn't get ecert for user: " + commanName)
	}
	cn := x509Cert.Subject.CommonName
	res := strings.Split(cn, "\\")
	affiliation := res[1]
	affiliationRole, _ := strconv.Atoi(res[2])
	return commanName, affiliation, affiliationRole, nil
}

func (t *CattleChaincode) getAllCattles(stub *shim.ChaincodeStub) ([]byte, error) {
	//rows, err := stub.GetRows(CATTLE_TABLE_NAME, []shim.Column{{Value: &shim.Column_String_{String_: "0000000000"}}})
	rows, err := stub.GetRows(CATTLE_TABLE_NAME, nil)
	var cattles []shim.Row
	for row := range rows {
		cattles = append(cattles, row)
	}
	if err != nil {
		return nil, err
	}
	response, err := json.Marshal(cattles)
	if err != nil {
		fmt.Printf("marshal row error: %s\n", err)
	}
	return response, nil
}

func main() {
	err := shim.Start(new(CattleChaincode))
	if err != nil {
		fmt.Printf("Error starting CattleChaincode : %s", err)
	}
}
