package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

const (
	LOAN_TABLE   = "loan_table"
	LENDER_TABLE = "lender_table"
)

var loanColumnTypes = []ColDef{
	{"FARM", "string"},
	{"LOAN_ID", "string"},
	{"LEND_DATE", "string"},
	{"LOAN_OFFICER", "string"},
	{"AMOUNT", "int64"},
	{"REPAY_DATE", "string"},
	{"TRACE", "bytes"},
}

var lenderColumnsTypes = []ColDef{
	{"LOAN_OFFICER", "string"},
	{"FARM", "string"},
	{"LOAN_ID", "string"},
}

var loanColumnsKeys = []bool{true, true, true, false, false, false, false}
var lenderColumnsKeys = []bool{true, true, true}

func createLoanTables(stub *shim.ChaincodeStub) error {
	if err := createTable(stub, LOAN_TABLE, loanColumnTypes, loanColumnsKeys); err != nil {
		return err
	}
	return createTable(stub, LENDER_TABLE, lenderColumnsTypes, lenderColumnsKeys)
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
	ccLogger.Debug(strconv.Itoa(rows) + " loan(s) in total for farm:" + farmId)
	returnVal, _ := json.Marshal(loans)
	return returnVal, nil
}

func getAllLoanIdByLender(stub *shim.ChaincodeStub, lenderId string) ([]byte, error) {
	rowsChan, err := stub.GetRows(LENDER_TABLE, []shim.Column{{Value: &shim.Column_String_{String_: lenderId}}})
	if err != nil {
		return nil, fmt.Errorf("getAllLoanIdByLender query failed. %s", err)
	}
	id := []string{}
	rows := 0
	for {
		select {
		case row, ok := <-rowsChan:
			if !ok {
				rowsChan = nil
			} else {
				rows++
				id = append(id, row.Columns[2].GetString_())
			}
		}
		if rowsChan == nil {
			break
		}
	}
	ccLogger.Debug(strconv.Itoa(rows) + " loan(s) in total for loan officer:" + lenderId)
	returnVal, _ := json.Marshal(id)
	return returnVal, nil
}

func generateLoanRow(loan *Loan) []*shim.Column {
	var loanColumns []*shim.Column
	loanColumns = append(loanColumns, &shim.Column{Value: &shim.Column_String_{String_: loan.Farm}})
	loanColumns = append(loanColumns, &shim.Column{Value: &shim.Column_String_{String_: loan.LoanId}})
	loanColumns = append(loanColumns, &shim.Column{Value: &shim.Column_String_{String_: loan.LendDate}})
	loanColumns = append(loanColumns, &shim.Column{Value: &shim.Column_String_{String_: loan.LoanOfficer}})
	loanColumns = append(loanColumns, &shim.Column{Value: &shim.Column_Int64{Int64: loan.Amount}})
	loanColumns = append(loanColumns, &shim.Column{Value: &shim.Column_String_{String_: loan.RepayDate}})
	traceMarshal, _ := json.Marshal(loan.Trace)
	loanColumns = append(loanColumns, &shim.Column{Value: &shim.Column_Bytes{Bytes: traceMarshal}})
	return loanColumns
}

func generateLenderRow(loanOfficer, farm, loanId string) []*shim.Column {
	var lenderColumns []*shim.Column
	lenderColumns = append(lenderColumns, &shim.Column{Value: &shim.Column_String_{String_: loanOfficer}})
	lenderColumns = append(lenderColumns, &shim.Column{Value: &shim.Column_String_{String_: farm}})
	lenderColumns = append(lenderColumns, &shim.Column{Value: &shim.Column_String_{String_: loanId}})
	return lenderColumns
}

func formatLoan(queryOutput shim.Row) *Loan {
	loan := new(Loan)
	loan.Farm = queryOutput.Columns[0].GetString_()
	loan.LoanId = queryOutput.Columns[1].GetString_()
	loan.LendDate = queryOutput.Columns[2].GetString_()
	loan.LoanOfficer = queryOutput.Columns[3].GetString_()
	loan.Amount = queryOutput.Columns[4].GetInt64()
	loan.RepayDate = queryOutput.Columns[5].GetString_()
	json.Unmarshal(queryOutput.Columns[6].GetBytes(), &loan.Trace)
	return loan
}

func populateSampleLoanRows(stub *shim.ChaincodeStub) {
	loan := Loan{}
	loan.Farm = "1234567"
	loan.LoanId = "TXK12SK9A1"
	loan.LendDate = "2016-03-18"
	loan.LoanOfficer = "ALPHA232"
	loan.Amount = 5000
	loan.RepayDate = "2017-03-18"
	trace0 := Loan_Trace{Date: "2016-03-18", Event: "loan release"}
	trace1 := Loan_Trace{Date: "2016-05-20", Event: "farm repaid 1000RMB"}
	loan.Trace = []*Loan_Trace{&trace0, &trace1}
	stub.InsertRow(LOAN_TABLE, shim.Row{Columns: generateLoanRow(&loan)})
	stub.InsertRow(LENDER_TABLE, shim.Row{Columns: generateLenderRow(loan.LoanOfficer, loan.Farm, loan.LoanId)})

	loan.Farm = "1234568"
	loan.LoanId = "I2KS12SJJS"
	loan.LoanOfficer = "BETA1290"
	loan.Trace = append(loan.Trace, &Loan_Trace{Date: "2016-07-18", Event: "farm repaid 1500RMB"})
	stub.InsertRow(LOAN_TABLE, shim.Row{Columns: generateLoanRow(&loan)})
	stub.InsertRow(LENDER_TABLE, shim.Row{Columns: generateLenderRow(loan.LoanOfficer, loan.Farm, loan.LoanId)})

}
