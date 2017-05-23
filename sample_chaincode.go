package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)


var logger = shim.NewLogger("mylogger")

type SampleChaincode struct {
}

type PersonalInfo struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	DOB       string `json:"DOB"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
}

type FinancialInfo struct {
	MonthlySalary      int `json:"monthlySalary"`
	MonthlyRent        int `json:"monthlyRent"`
	OtherExpenditure   int `json:"otherExpenditure"`
	MonthlyLoanPayment int `json:"monthlyLoanPayment"`
}

type LoanApplication struct {
	ID                     string        `json:"id"`
	PropertyId             string        `json:"propertyId"`
	LandId                 string        `json:"landId"`
	PermitId               string        `json:"permitId"`
	BuyerId                string        `json:"buyerId"`
	AppraisalApplicationId string        `json:"appraiserApplicationId"`
	SalesContractId        string        `json:"salesContractId"`
	PersonalInfo           PersonalInfo  `json:"personalInfo"`
	FinancialInfo          FinancialInfo `json:"financialInfo"`
	Status                 string        `json:"status"`
	RequestedAmount        int           `json:"requestedAmount"`
	FairMarketValue        int           `json:"fairMarketValue"`
	ApprovedAmount         int           `json:"approvedAmount"`
	ReviewerId             string        `json:"reviewerId"`
	LastModifiedDate       string        `json:"lastModifiedDate"`
}

type Participant struct {
	ID                     string        `json:"id"`
	Name                   string        `json:"name"`
	AssetId								 string        `json:"loanId"`
  SharePerCent           int           `json:"share"`
	ShareAmount            int 					 `json:"shareAmount"`
}

func GetLoanApplication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering GetLoanApplication")

	if len(args) < 1 {
		logger.Error("Invalid number of arguments")
		return nil, errors.New("Missing loan application ID")
	}

	var loanApplicationId = args[0]
	bytes, err := stub.GetState(loanApplicationId)
	if err != nil {
		logger.Error("Could not fetch loan application with id "+loanApplicationId+" from ledger", err)
		return nil, err
	}
	return bytes, nil
}

func CreateLoanApplication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering CreateLoanApplication")

	if len(args) < 2 {
		logger.Error("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for loan application creation")
	}

	var loanApplicationId = args[0]
	var loanApplicationInput LoanApplication
	loanApplicationInput = LoanApplication{ID:"ID1",PropertyId:"prop1",LandId:"land1"}
	bytes, err1 := json.Marshal (&loanApplicationInput)
	 if err1 != nil {
		         fmt.Println("Could not marshal personal info object", err1)
			         return nil, err1
				  }

	err := stub.PutState(loanApplicationId, bytes )
	if err != nil {
		logger.Error("Could not save loan application to ledger", err)
		return nil, err
	}

	var customEvent = "{eventType: 'loanApplicationCreation', description:" + loanApplicationId + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return nil, err
	}
	logger.Info("Successfully saved loan application")
	return bytes, nil

}


func NonDeterministicFunction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering NonDeterministicFunction")
	//Use random number generator to generate the ID
	var random = rand.New(rand.NewSource(time.Now().UnixNano()))
	var loanApplicationID = "la1" + strconv.Itoa(random.Intn(1000))
	var loanApplication = args[0]
	var la LoanApplication
	err := json.Unmarshal([]byte(loanApplication), &la)
	if err != nil {
		fmt.Println("Could not unmarshal loan application", err)
		return nil, err
	}
	la.ID = loanApplicationID
	laBytes, err := json.Marshal(&la)
	if err != nil {
		fmt.Println("Could not marshal loan application", err)
		return nil, err
	}
	err = stub.PutState(loanApplicationID, laBytes)
	if err != nil {
		fmt.Println("Could not save loan application to ledger", err)
		return nil, err
	}

	fmt.Println("Successfully saved loan application")
	return []byte(loanApplicationID), nil
}

func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SampleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}


func GetCertAttribute(stub shim.ChaincodeStubInterface, attributeName string) (string, error) {
	logger.Debug("Entering GetCertAttribute")
	attr, err := stub.ReadCertAttribute(attributeName)
	if err != nil {
		return "", errors.New("Couldn't get attribute " + attributeName + ". Error: " + err.Error())
	}
	attrString := string(attr)
	return attrString, nil
}


func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "CreateLoanApplication" {
		//username, _ := GetCertAttribute(stub, "username")
		//role, _ := GetCertAttribute(stub, "role")

		return CreateLoanApplication(stub, args)
	/*	if role == "Bank_Admin" {
		return CreateLoanApplication(stub, args)
		} else {
			return nil, errors.New(username + " with role " + role + " does not have access to create a loan application")
		}*/

	}
	return nil, errors.New("Invalid function name")
}



/*func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "CreateLoanApplication" {
		username, _ := GetCertAttribute(stub, "username")
		role, _ := GetCertAttribute(stub, "role")
		if role == "Bank_Home_Loan_Admin" {
			return CreateLoanApplication(stub, args)
		} else {
			return nil, errors.New(username + " with role " + role + " does not have access to create a loan application")
		}

	}
	return nil, nil
}*/



func main() {

	lld, _ := shim.LogLevel("DEBUG")
	fmt.Println(lld)

	logger.SetLevel(lld)
	fmt.Println(logger.IsEnabledFor(lld))

	err := shim.Start(new(SampleChaincode))
	if err != nil {
		logger.Error("Could not start SampleChaincode")
	} else {
		logger.Info("SampleChaincode successfully started")
	}

}
