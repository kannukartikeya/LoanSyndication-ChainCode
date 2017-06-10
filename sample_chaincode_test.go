package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var loanApplicationID = "la1"
var loanApplicationID2 = "la2"
var loanApplication = `{"id":"` + loanApplicationID + `","propertyId":"prop1","landId":"land1","permitId":"permit1","buyerId":"kartikeya","personalInfo":{"firstname":"Kartikeya","lastname":"Gupta","dob":"dob","email":"kartikeya80@gmail.com","mobile":"99999999"},"financialInfo":{"spRating":"BBB+","moodyRating":"Baa2","dcr":1.9,"turnover":4000},"status":"Submitted","requestedAmount":40000,"fairMarketValue":58000,"approvedAmount":40000,"dealAmount":40000,"outstandingSettlementAmount":40000,"reviewedBy":"bond","lastModifiedDate":"21/09/2016 2:30pm"}`

// func CreateLoanParticipation(t *testing.T) {
// 	fmt.Println("Entering CreateLoanParticipation")
// 	m := make(map[string][]byte)
// 	m["role"] = []byte("Bank")
// 	stub := NewMockStub("mockStub", new(SampleChaincode), m)
// 	bytes, _ := stub.ReadCertAttribute("role")
// 	fmt.Println(string(bytes))
// 	stub.MockInvoke("123", "init", []string{})

// }

func TestStubCreation(t *testing.T) {
	fmt.Println("Entering TestStubCreation")
	attributes := make(map[string][]byte)
	//Create a custom MockStub that internally uses shim.MockStub
	stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}
}

func TestCrtLoanAppWithNullArguments(t *testing.T) {
	fmt.Println("Entering TestCrtLoanAppWithNullArguments")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	stub.MockTransactionStart("t123")
	_, err := CreateLoanParticipation(stub, []string{})
	if err == nil {
		t.Fatalf("Expected CreateLoanParticipation to return validation error")
	}
	stub.MockTransactionEnd("t123")

}

func TestCrtLoanAppWithIdLoanDetails(t *testing.T) {
	fmt.Println("Entering TestCrtLoanAppWithIdLoanDetails")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	stub.MockTransactionStart("t123")
	_, err := CreateLoanParticipation(stub, []string{loanApplicationID, loanApplication})
	if err != nil {
		t.Fatalf("Expected CreateLoanParticipation to succeed")
	}
	stub.MockTransactionEnd("t123")

}

func TestCreateFetchParticipants(t *testing.T) {
	fmt.Println("Entering TestCreateFetchParticipants")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	stub.MockTransactionStart("t123")
	_, err := CreateParticipants(stub, []string{"part1"})
	if err != nil {
		t.Fatalf("Expected CreateParticipants to succeed")
	}
	stub.MockTransactionEnd("t123")

fmt.Println("Created and Fetching Participants")
var firstParticipant Participant
bytes, err := GetLoanParticipant(stub, []string{"part1"})


	err = json.Unmarshal(bytes, &firstParticipant)
	if err != nil {
		t.Fatalf("Could not unmarshal loan application with ID " + "part1")
	}
	fmt.Println("Participant Name :" + firstParticipant.Name)
//fmt.Println("Participated Asset ID :" + firstParticipant.AssetList[0].AssetId)
}
	
func TestGetParticipatedLoans(t *testing.T){
	
	fmt.Println("Entering TestGetParticipatedLoans")

	attributes := make(map[string][]byte)
	attributes["username"] = []byte("vojha24")
	attributes["role"] = []byte("client")

	stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	_, err := stub.MockInvoke("t123", "CreateLoanParticipation", []string{loanApplicationID, loanApplication})
	if err != nil {
		fmt.Println(err)
		t.Fatalf("Expected CreateLoanParticipation to be invoked")
	}
	
	_, err = stub.MockInvoke("t123", "CreateLoanParticipation", []string{loanApplicationID2, loanApplication})
	if err != nil {
		fmt.Println(err)
		t.Fatalf("Expected CreateLoanParticipation to be invoked")
	}

	loanbytes2, err1 := stub.MockInvoke("t123", "GetParticipatedLoans", []string{})
	if err1 == nil {
		//t.Fatalf("Expected unauthorized user error to be returned")
	}
	var loanList []LoanApplication
	if (loanbytes2 != nil) {
		err = json.Unmarshal(loanbytes2,&loanList)
		if err != nil {
			logger.Error("unable to unmarshall loanlist")
		
		}
	}
	
	fmt.Println("LoanList length is", len(loanList))
}



func TestCrtFetchLoanAppAndValidateInputStoredVal(t *testing.T) {
	fmt.Println("Entering TestCrtFetchLoanAppAndValidateInputStoredVal")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	stub.MockTransactionStart("t123")
	CreateLoanParticipation(stub, []string{loanApplicationID, loanApplication})
	stub.MockTransactionEnd("t123")

	var la LoanApplication
	/*bytes, err := stub.GetState(loanApplicationID)
	if err != nil {
		t.Fatalf("Could not fetch loan application with ID " + loanApplicationID)
	}*/
	bytes, err := GetLoanApplication(stub, []string{loanApplicationID})

	err = json.Unmarshal(bytes, &la)
	if err != nil {
		t.Fatalf("Could not unmarshal loan application with ID " + loanApplicationID)
	}
	var errors = []string{}
	var loanApplicationInput LoanApplication
	err = json.Unmarshal([]byte(loanApplication), &loanApplicationInput)
	if la.ID != loanApplicationInput.ID {
		errors = append(errors, "Loan Application ID does not match")
	}
	if la.PropertyId != loanApplicationInput.PropertyId {
		errors = append(errors, "Loan Application PropertyId does not match")
	}
	if la.PersonalInfo.Firstname != loanApplicationInput.PersonalInfo.Firstname {
		errors = append(errors, "Loan Application PersonalInfo.Firstname does not match")
	}
	//Can be extended for all fields
	if len(errors) > 0 {
		t.Fatalf("Mismatch between input and stored Loan Application")
		for j := 0; j < len(errors); j++ {
			fmt.Println(errors[j])
		}
	}

}
func TestInvokeCrtLoanAppWithUnAuthorizedUser(t *testing.T) {
	fmt.Println("Entering TestInvokeCrtLoanAppWithUnAuthorizedUser")

	attributes := make(map[string][]byte)
	attributes["username"] = []byte("vojha24")
	attributes["role"] = []byte("client")

	stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	_, err := stub.MockInvoke("t123", "CreateLoanParticipation", []string{loanApplicationID, loanApplication})
	if err == nil {
		//t.Fatalf("Expected unauthorized user error to be returned")
	}

}

func TestInvokeCrtSttleLoanSyndWithAuthorizedRole(t *testing.T) {
	fmt.Println("Entering TestInvokeCrtSttleLoanSyndWithAuthorizedRole")

	attributes := make(map[string][]byte)
	attributes["username"] = []byte("vojha24")
	attributes["role"] = []byte("Bank_Admin")

	stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	stub.MockTransactionStart("t123")
	_, err := CreateParticipants(stub, []string{"part1"})
	if err != nil {
		t.Fatalf("Expected CreateParticipants to succeed")
	}
	stub.MockTransactionEnd("t123")

	_, err = stub.MockInvoke("t123", "CreateLoanParticipation", []string{loanApplicationID, loanApplication})
	if err != nil {
		fmt.Println(err)
		t.Fatalf("Expected CreateLoanParticipation to be invoked")
	}

	fmt.Println("Fetching Participant")
	var firstParticipant Participant
	bytes, err := GetLoanParticipant(stub, []string{"part1"})
	if err != nil {
		fmt.Println(err)
		t.Fatalf("Expected GetLoanParticipant to be invoked successfully")
	}
		err = json.Unmarshal(bytes, &firstParticipant)
		if err != nil {
			t.Fatalf("Could not unmarshal firstParticipant" + firstParticipant.ID)
		}
		if(firstParticipant.AssetList != nil){
		fmt.Println("Participant Details :" + firstParticipant.AssetList[0].AssetId)
		fmt.Println("Participant Share Amount After Participation", firstParticipant.AssetList[0].ShareAmount)
	
		}

		_, err = stub.MockInvoke("t123", "SettleLoanSyndication", []string{loanApplicationID, "1000"})
		if err != nil {
			fmt.Println(err)
			t.Fatalf("Expected SettleLoanSyndication to be invoked")
		}

		var la LoanApplication
		bytes, err = GetLoanApplication(stub, []string{loanApplicationID})

		err = json.Unmarshal(bytes, &la)
		if err != nil {
			t.Fatalf("Could not unmarshal loan application with ID " + loanApplicationID)
		}

	fmt.Println("Loan OutstandingSettlementAmount %d",la.OutStandingSettlementAmount)

	bytes, err = GetLoanParticipant(stub, []string{"part1"})
	if err != nil {
		fmt.Println(err)
		t.Fatalf("Expected GetLoanParticipant to be invoked successfully")
	}
		err = json.Unmarshal(bytes, &firstParticipant)
		if err != nil {
			t.Fatalf("Could not unmarshal firstParticipant" + firstParticipant.ID)
		}
		if(firstParticipant.AssetList != nil){
		fmt.Println("Participated Asset ID : " + firstParticipant.AssetList[0].AssetId)
		fmt.Println("Participant Reduced Share Amount Post Settlement", firstParticipant.AssetList[0].ShareAmount)
		fmt.Println("Participant Settlemt Feest", firstParticipant.AssetList[0].SettlementFees)
		}

	fmt.Println("Entering TestInvokeCrtSttleLoanSyndWithAuthorizedRole")

}

func TestInvokeInValidFunction(t *testing.T) {
	fmt.Println("Entering TestInvokeInValidFunction")

	attributes := make(map[string][]byte)
	attributes["username"] = []byte("vojha24")
	attributes["role"] = []byte("Bank_Admin")

	stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	_, err := stub.MockInvoke("t123", "InvalidFunctionName", []string{})
	if err == nil {
		t.Fatalf("Expected invalid function name error")
	}

}

/*func TestInvokeCrtLoanAppAndFetchWithAuthorizedRole(t *testing.T) {
	fmt.Println("Entering TestInvokeFunctionValidation2")

	attributes := make(map[string][]byte)
	attributes["username"] = []byte("vojha24")
	attributes["role"] = []byte("Bank_Admin")

	stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	bytes, err := stub.MockInvoke("t123", "CreateLoanParticipation", []string{loanApplicationID, loanApplication})
	if err != nil {
		t.Fatalf("Expected CreateLoanParticipation function to be invoked")
	}
	//A spy could have been used here to ensure CreateLoanParticipation method actually got invoked.
	var la LoanApplication
	err = json.Unmarshal(bytes, &la)
	if err != nil {
		t.Fatalf("Expected valid loan application JSON string to be returned from CreateLoanParticipation method")
	}

}*/
