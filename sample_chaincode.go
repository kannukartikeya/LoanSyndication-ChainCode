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



func GetLoanParticipant(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering GetLoanParticipant")

	if len(args) < 1 {
		logger.Error("Invalid number of arguments")
		return nil, errors.New("Missing participant ID")
	}

	var participantID = args[0]
	bytes, err := stub.GetState(participantID)
	if err != nil {
		logger.Error("Could not fetch participant with id "+participantID+" from ledger", err)
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
	
	
	/*var loanApplicationInput LoanApplication
	loanApplicationInput = LoanApplication{ID:loanApplicationId,PropertyId:"prop1",LandId:"land1",ApprovedAmount:1000}
	bytes, err1 := json.Marshal (&loanApplicationInput)
	 if err1 != nil {
		         fmt.Println("Could not marshal personal info object", err1)
			         return nil, err1
				  }

	err := stub.PutState(loanApplicationId, bytes )*/
	
	
	var loanApplicationInput = args[1]

	err := stub.PutState(loanApplicationId, []byte(loanApplicationInput))
	if err != nil {
		logger.Error("Could not save loan application to ledger", err)
		return nil, err
	}

	partbytes, err := stub.GetState("part1")
	if err != nil {
		logger.Error("Could not fetch firstParticipant with id part1 from ledger", err)
		return nil, err
	}

 var firstParticipant Participant
 err = json.Unmarshal(partbytes,&firstParticipant)
 fmt.Println(firstParticipant.Name)
 
 //firstParticipant.ShareAmount = loanApplicationInput.ApprovedAmount *  firstParticipant.SharePerCent
 
 firstParticipant.ShareAmount = 1000 *  firstParticipant.SharePerCent

 partbytes2, err := json.Marshal (&firstParticipant)
 if err != nil {
        fmt.Println("Could not marshal firstParticipant info object", err)
        return nil, err
 }
 err = stub.PutState("part1", partbytes2)

	var customEvent = "{eventType: 'loanApplicationCreation', description:" + loanApplicationId + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return nil, err
	}
	logger.Info("Successfully saved loan application")
	return partbytes2, nil

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
//resets all the things
func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args)!=1 {
				return nil,errors.New("Incorrect number of arguments. Expecting 1")
	}

	var firstParticipant Participant
	//var secondParticipant Participant

	firstParticipant = Participant{ID:"part1",Name:"DeucheBank",SharePerCent:80}
	//secondParticipant = Participant{ID:"part2",Name:"CITIBank",SharePerCent:20}

	bytes, err1 := json.Marshal (&firstParticipant)
	 if err1 != nil {
		         fmt.Println("Could not marshal firstParticipant object", err1)
			         return nil, err1
				  }

	err := stub.PutState("part1", bytes )
	if err != nil {
			logger.Error("Could not save firstParticipant to ledger", err)
			return nil, err
		}

	return nil, nil
}

func (t *SampleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "GetLoanApplication" {
		return GetLoanApplication(stub, args)
	}
	if function == "GetLoanParticipant" {
		return GetLoanParticipant(stub, args)
	}
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
