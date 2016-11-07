package main

import "fmt"
import "github.com/hyperledger/fabric/core/chaincode/shim"
import "github.com/op/go-logging"

var ccLogger = logging.MustGetLogger("FarmCC")

type FarmCC struct {
}

func (t *FarmCC) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	ccLogger.Debug("Init called!")

	return nil, nil
}

func (t *FarmCC) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	ccLogger.Debug("Invoke called!")
	return nil, nil
}

func (t *FarmCC) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	ccLogger.Debug("Query called!")
	return nil, nil
}

func main() {
	err := shim.Start(new(FarmCC))
	if err != nil {
		fmt.Printf("Error starting FarmCC : %s", err)
	}
}

func (t *FarmCC) createFarmTable(stub *shim.ChaincodeStub) error {
	return nil
}
