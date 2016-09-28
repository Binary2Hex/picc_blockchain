/*
Copyright IBM Corp. 2016 All Rights Reserved.
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
	"regexp"
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
	Id            string `json:"id"`
	Vaccinated    bool   `json:"vaccinated"`
	InsuranceID   string `json:"insuranceID"`
	InsuranceCorp string `json:"insuranceCorp"`
	Loan          int    `json:"loan"`
	LoanID        string `json:"loanID"`
	LoanCorp      string `json:"loanCorp"`
	Origin        string `json:"origin"`
	Trader        string `json:"trader"`
	Status        int    `json:"status"`
	Owner         string `json:"owner"`
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

type CattleChaincode struct {
}

func (t *CattleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	if function == "init" {
		// TODO, THIS IS JUST FOR TESTING
		cattle0 := Cattle{Id: "0000000000"}
		cattle1 := Cattle{Id: "0000000001"}
		cattle0Bytes, err := json.Marshal(cattle0)
		if err != nil {
			return nil, err
		}
		stub.PutState(cattle0.Id, cattle0Bytes)
		cattle1Bytes, err := json.Marshal(cattle1)
		if err != nil {
			return nil, err
		}
		stub.PutState(cattle1.Id, cattle1Bytes)
		return nil, nil
	}

	return nil, errors.New("function: " + function + " not supported!")

}

func (t *CattleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	username, _, affiliationRole, err := t.getCaller(stub)
	if err != nil {
		return nil, err
	}

	if function == "createCattle" {
		if affiliationRole != GOV {
			return nil, errors.New("Unauthorized to create new cattle")
		}
		if len(args) != 2 {
			return nil, errors.New("2 args are required for createCattle")
		}
		return t.createCattle(stub, args[0], args[1])
	}
	if function == "traderToFarm" {
		if len(args) != 3 {
			return nil, errors.New("3 args are required for traderToFarm")
		}
		return t.traderToFarm(stub, username, args[0], args[1], args[2])
	}

	if function == "vaccinate" {
		if affiliationRole != GOV {
			return nil, errors.New("Unauthorized to create new cattle")
		}
		if len(args) != 1 {
			return nil, errors.New("1 arg is required for vaccinate")
		}
		return t.vaccinate(stub, args[0])
	}

	if function == "applyForInsurance" {
		if affiliationRole != FARM {
			return nil, errors.New("only farm can apply for insurance")
		}
		// TODO, to be implemented!!!
	}

	if function == "insure" {
		if affiliationRole != INSURANCE {
			return nil, errors.New("only insurance company can insure")
		}
		// TODO, to be implemented!!!
	}

	if function == "applyForLoan" {
		if affiliationRole != FARM {
			return nil, errors.New("only farm can apply for loan")
		}
		// TODO, to be implemented!!!
	}

	if function == "loan" {
		if affiliationRole != BANK {
			return nil, errors.New("only Bank can serve a loan")
		}
		// TODO, to be implemented!!!
	}

	fmt.Println("function: " + function + " not supported!")
	return nil, errors.New("function: " + function + " not supported!")
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
		if affiliationRole != GOV {
			return nil, errors.New("only GOV can do a full query")
		}
		return t.getAllCattles(stub)
	}

	if function == "getAllMyCattles" {
		// TODO, to be implemented
		return nil, nil
	}

	if function == "getCattleByID" {
		if len(args) != 1 {
			return nil, errors.New("need 1 arg for getCattleByID")
		}
		// TODO, to be implemented
		return nil, nil
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
	iter, err := stub.RangeQueryState("0", ":")
	if err != nil {
		return nil, err
	}

	var cattles []Cattle
	var cattle Cattle
	for iter.HasNext() {
		_, bytes, err := iter.Next()
		if err == nil {
			json.Unmarshal(bytes, &cattle)
			cattles = append(cattles, cattle)
		}
	}

	res, err := json.Marshal(cattles)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (t *CattleChaincode) createCattle(stub *shim.ChaincodeStub, id string, trader string) ([]byte, error) {
	fmt.Println("createCattle called")
	matched, err := regexp.Match("^[0-9]{10}$", []byte(id))
	if matched == false {
		fmt.Println("id format error")
		return nil, err
	}
	bytes, err := stub.GetState(id)
	if err != nil {
		return nil, err
	}
	if bytes != nil {
		fmt.Println("id: " + id + " already exists")
		return nil, errors.New("id: " + id + " already exists")
	}
	cattle := Cattle{Id: id, Trader: trader, Owner: trader}
	cattleBytes, err := json.Marshal(cattle)
	if err != nil {
		return nil, err
	}
	stub.PutState(id, cattleBytes)
	fmt.Println("successfully created with id: " + id)
	return []byte("successfully created with id: " + id), nil
}

func (t *CattleChaincode) traderToFarm(stub *shim.ChaincodeStub, username, id string, trader, farm string) ([]byte, error) {
	bytes, err := stub.GetState(id)
	if err != nil {
		return nil, err
	}
	if bytes == nil {
		return nil, errors.New("no cattle found with id: " + id)
	}

	var cattle Cattle
	err = json.Unmarshal(bytes, &cattle)
	if err != nil {
		return nil, err
	}

	if cattle.Owner != username {
		return nil, errors.New("unauthorized to transfer from trader to farm")
	}
	cattle.Owner = farm
	return t.rawSave(stub, cattle)
}

func (t CattleChaincode) vaccinate(stub *shim.ChaincodeStub, id string) ([]byte, error) {
	bytes, err := stub.GetState(id)
	if err != nil {
		return nil, err
	}
	if bytes == nil {
		return nil, errors.New("no cattle found with id: " + id)
	}

	var cattle Cattle
	err = json.Unmarshal(bytes, &cattle)
	if err != nil {
		return nil, err
	}
	cattle.Vaccinated = true
	return t.rawSave(stub, cattle)
}

func (t CattleChaincode) rawSave(stub *shim.ChaincodeStub, cattle Cattle) ([]byte, error) {
	bytes, err := json.Marshal(cattle)
	if err != nil {
		return nil, err
	}
	stub.PutState(cattle.Id, bytes)
	fmt.Println("save cattle with id:" + cattle.Id + "completed")
	return []byte("save cattle with id:" + cattle.Id + "completed"), nil
}

func main() {
	err := shim.Start(new(CattleChaincode))
	if err != nil {
		fmt.Printf("Error starting CattleChaincode : %s", err)
	}
}
