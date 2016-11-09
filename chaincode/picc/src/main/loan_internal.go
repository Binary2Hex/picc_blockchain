package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

const (
	LOAN_TABLE = "loan_table"
)

var loanColumnTypes = []ColDef{
	{"FARM", "string"},
	{"LEND_DATE", "string"},
	{"LOAN_OFFICER", "string"},
	{"AMOUNT", "int64"},
	{"REPAY_DATE", "string"},
	{"TRACE", "bytes"},
}

var loanColumnsKeys = []bool{true, true, false, false, false, false}

func createLoanTable(stub *shim.ChaincodeStub) error {
	return createTable(stub, LOAN_TABLE, loanColumnTypes, loanColumnsKeys)
}

func getAllLoansByFarm(stub *shim.ChaincodeStub, farmId string) ([]byte, error) {
	rowsChan, err := stub.GetRows(LOAN_TABLE, []shim.Column{{Value: &shim.Column_String_{String_: farmId}}})
	if err != nil {
		return nil, fmt.Errorf("getAllLoansByFarm query failed. %s", err)
	}
	var loans []*Loan
	rows := 0
	for {
		select {
		case row, ok := <-rowsChan:
			if !ok {
				rowsChan = nil
			} else {
				rows++
				loan := formatLoan(row)
				loans = append(loans, loan)
			}
		}
		if rowsChan == nil {
			break
		}
	}
	ccLogger.Debug(strconv.Itoa(rows) + " insurances in total for farm:" + farmId)
	returnVal, _ := json.Marshal(loans)
	return returnVal, nil
}

func generateLoanRow(loan *Loan) []*shim.Column {
	var loanColumns []*shim.Column
	loanColumns = append(loanColumns, &shim.Column{Value: &shim.Column_String_{String_: loan.Farm}})
	loanColumns = append(loanColumns, &shim.Column{Value: &shim.Column_String_{String_: loan.LendDate}})
	loanColumns = append(loanColumns, &shim.Column{Value: &shim.Column_String_{String_: loan.LoanOfficer}})
	loanColumns = append(loanColumns, &shim.Column{Value: &shim.Column_Int64{Int64: loan.Amount}})
	loanColumns = append(loanColumns, &shim.Column{Value: &shim.Column_String_{String_: loan.RepayDate}})
	traceMarshal, _ := json.Marshal(loan.Trace)
	loanColumns = append(loanColumns, &shim.Column{Value: &shim.Column_Bytes{Bytes: traceMarshal}})
	return loanColumns
}

func formatLoan(queryOutput shim.Row) *Loan {
	loan := new(Loan)
	loan.Farm = queryOutput.Columns[0].GetString_()
	loan.LendDate = queryOutput.Columns[1].GetString_()
	loan.LoanOfficer = queryOutput.Columns[2].GetString_()
	loan.Amount = queryOutput.Columns[3].GetInt64()
	loan.RepayDate = queryOutput.Columns[4].GetString_()
	json.Unmarshal(queryOutput.Columns[5].GetBytes(), &loan.Trace)
	return loan
}

func populateSampleLoanRows(stub *shim.ChaincodeStub) {
	loan := Loan{}
	loan.Farm = "1234567"
	loan.LendDate = "2016-03-18"
	loan.LoanOfficer = "ALPHA"
	loan.Amount = 5000
	loan.RepayDate = "2017-03-18"
	trace0 := Loan_Trace{Date: "2016-03-18", Event: "loan release"}
	trace1 := Loan_Trace{Date: "2016-05-20", Event: "farm repaid 1000RMB"}
	loan.Trace = []*Loan_Trace{&trace0, &trace1}
	stub.InsertRow(LOAN_TABLE, shim.Row{Columns: generateLoanRow(&loan)})

	loan.Farm = "1234568"
	loan.Trace = append(loan.Trace, &Loan_Trace{Date: "2016-07-18", Event: "farm repaid 1500RMB"})
	stub.InsertRow(LOAN_TABLE, shim.Row{Columns: generateLoanRow(&loan)})

}
