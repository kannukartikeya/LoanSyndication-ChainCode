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
	OutStandingSettlementAmount      int `json:"outstandingSettlementAmount"`
	ReviewerId             string        `json:"reviewerId"`
	LastModifiedDate       string        `json:"lastModifiedDate"`
}

type LoanList struct {
	Loans []LoanApplication
}

type Participant struct {
	ID                     string        `json:"id"`
	Name                   string        `json:"name"`
	AssetList []Asset
}
type Asset struct{
	AssetId								 string        `json:"loanId"`
	SharePerCent           int           `json:"share"`
	ShareAmount            int 					 `json:"shareAmount"`
	SyndicatedAmount 			 int					 `json:"syndicatedAmount"`
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

func GetParticipatedLoans(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering GetParticipatedLoans")

	bytes, err := stub.GetState("loanlist")
	if err != nil {
		logger.Error("Could not fetch loanlist from ledger", err)
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

func CreateParticipants(stub shim.ChaincodeStubInterface,args []string)([]byte,error){
	logger.Debug("Entering CreateParticipants")
	

	if len(args) < 1 {
		logger.Error("Invalid number of arguments")
		return nil, errors.New("Missing participant ID")
	}
	
	var firstParticipant Participant
	
	var participantID = args[0]
	
	//var secondParticipant Participant

	/*secondParticipant = Participant{ID:"part1",Name:"DeucheBank",
								AssetList: []Asset{
									{AssetId:"la1",
									SharePerCent:80,
									ShareAmount:0,
									SyndicatedAmount:1000},
								},
						}*/
	
	firstParticipant = Participant{ID:participantID ,Name:"DeucheBank"}
	
	//secondParticipant = Participant{ID:"part2",Name:"CITIBank",SharePerCent:20}

	bytes, err1 := json.Marshal (&firstParticipant)
	 if err1 != nil {
		         fmt.Println("Could not marshal firstParticipant object", err1)
			         return nil, err1
				  }

	err := stub.PutState(participantID, bytes )
	if err != nil {
			logger.Error("Could not save firstParticipant to ledger", err)
			return nil, err
		}

		return nil,nil

}

func CreateLoanParticipation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering CreateLoanParticipation")

	if len(args) < 2 {
		logger.Error("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for loan application creation")
	}

	var loanApplicationId = args[0]
	var loanApplicationInput = args[1]

	err := stub.PutState(loanApplicationId, []byte(loanApplicationInput))
	if err != nil {
		logger.Error("Could not save loan application to ledger", err)
		return nil, err
	}

	var participatedLoan LoanApplication
	err = json.Unmarshal([]byte(loanApplicationInput),&participatedLoan)
	if err != nil {
		return nil, err
	}
    fmt.Println("participatedLoan ID and amount " + participatedLoan.ID, participatedLoan.ApprovedAmount)
    

	loanbytes2, err := AppendToLoanList(stub,participatedLoan)
	    
    //err = ParticipateLoan(stub, "part1",loanApplicationInput, participatedLoan.ApprovedAmount)

	var customEvent = "{eventType: 'loanApplicationCreation', description:" + loanApplicationId + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return nil, err
	}
	logger.Info("Successfully saved loan application")
	return loanbytes2,err
}

func AppendToLoanList(stub shim.ChaincodeStubInterface,  participatedLoan LoanApplication) ([]byte, error){

    var loanList []LoanApplication
    bytes , err := stub.GetState("loanlist")
	if err != nil {
		logger.Error("Could not fetch firstParticipant with id part1 from ledger", err)
		return nil, err
	}
	if ( bytes == nil) {
		loanList = append(loanList,participatedLoan);
	} else {
		err = json.Unmarshal(bytes,&loanList)
		if err != nil {
			logger.Error("unable to unmarshall loanlist")
		return nil, err
		}
	
		loanList = append(loanList,participatedLoan)
	//fmt.Println("firstParticipant Name" + firstParticipant.Name)
	}
	
	 loanbytes2, err := json.Marshal (&loanList)
	 if err != nil {
        fmt.Println("Could not marshal loanList object", err)
        return nil, err
	 }
	
	err = stub.PutState("loanlist", loanbytes2)
	if err != nil {
		return nil, err
	}	
//	fmt.Println("LoanList length is %d", len(loanList))
	
	return loanbytes2,nil

	
}

func ParticipateLoan(stub shim.ChaincodeStubInterface, participant string, loan_id string , participationAmount int) (error){
	
	partbytes, err := stub.GetState(participant)
	if err != nil || partbytes == nil {
		logger.Error("Could not fetch firstParticipant with id part1 from ledger", err)
		return  err
	}

	var firstParticipant Participant
	err = json.Unmarshal(partbytes,&firstParticipant)
	if err != nil {
		return err
	}
	fmt.Println("firstParticipant Name" + firstParticipant.Name)

	var newAsset Asset
	newAsset.AssetId= loan_id
	newAsset.SharePerCent = 80
	newAsset.ShareAmount = ( participationAmount * newAsset.SharePerCent / 100 )
	
	firstParticipant.AssetList = append(firstParticipant.AssetList, newAsset)
	
	fmt.Println("Total loans participated")
	fmt.Println(len(firstParticipant.AssetList))
	

	/*for _, elem := range firstParticipant.AssetList {
	fmt.Println("Participant Asset Details :" + elem.AssetId)
	firstParticipant.AssetList[0].ShareAmount = ( participationAmount * elem.SharePerCent / 100 )
	fmt.Println("elem.ShareAmount")
	fmt.Println(firstParticipant.AssetList[0].ShareAmount)
	}*/
 //firstParticipant.ShareAmount = loanApplicationInput.ApprovedAmount *  firstParticipant.SharePerCent

 //firstParticipant.ShareAmount = 1000 *  firstParticipant.SharePerCent

	 partbytes2, err := json.Marshal (&firstParticipant)
	 if err != nil {
        fmt.Println("Could not marshal firstParticipant info object", err)
        return err
	 }
	err = stub.PutState("part1", partbytes2)
	if err != nil {
		return err
	}
		
	return nil
	
}


func SettleLoanSyndication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering SettleLoanSyndication")

	if len(args) < 2 {
		logger.Error("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for loan settlement")
	}

	var loanApplicationId = args[0]
	var loanSettlementAmount = args[1]

	fmt.Printf("Settle Loan : %s, for :%s", loanApplicationId, loanSettlementAmount)

	v, err := strconv.Atoi(loanSettlementAmount)

	bytes, err := stub.GetState(loanApplicationId)
	if err != nil {
		logger.Error("Could not fetch loan application with id "+loanApplicationId+" from ledger", err)
		return nil, err
	}

	var participatedLoan LoanApplication
    err = json.Unmarshal(bytes,&participatedLoan)
    fmt.Println("participatedLoan ID and amount " + participatedLoan.ID, participatedLoan.ApprovedAmount)

	fmt.Println("updating outStandingSettlentAmount for ID for amount " + loanSettlementAmount)

	//participatedLoan.OutStandingSettlementAmount = participatedLoan.ApprovedAmount - v
	participatedLoan.OutStandingSettlementAmount = participatedLoan.OutStandingSettlementAmount - v

	laBytes, err := json.Marshal(&participatedLoan)
	if err != nil {
		fmt.Println("Could not marshal loan application", err)
		return nil, err
	}
	err = stub.PutState(loanApplicationId, laBytes)
	if err != nil {
		fmt.Println("Could not save loan application to ledger", err)
		return nil, err
	}

	//err = SettleParticipation(stub,"part1",loanApplicationId,v)
	
	var customEvent = "{eventType: 'loanApplicationCreation', description:" + loanApplicationId + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return nil, err
	}
	logger.Info("Successfully saved loan application")
	return laBytes, nil

return nil,nil
}

func SettleParticipation(stub shim.ChaincodeStubInterface, participant string, loan_id string , settlementAmount int) (error){
	partbytes, err := stub.GetState(participant)
	if err != nil || partbytes == nil {
		logger.Error("Could not fetch firstParticipant with id part1 from ledger", err)
		return err
	}

	 var firstParticipant Participant
	 err = json.Unmarshal(partbytes,&firstParticipant)
	 fmt.Println("firstParticipant Name" + firstParticipant.Name)

	for _, elem := range firstParticipant.AssetList {
	fmt.Println("Participant Asset Details :" + elem.AssetId)
	firstParticipant.AssetList[0].ShareAmount = firstParticipant.AssetList[0].ShareAmount - (firstParticipant.AssetList[0].SharePerCent*settlementAmount/100)
	fmt.Println("Update Participant ShareAmount")
	fmt.Println(firstParticipant.AssetList[0].ShareAmount)
	}
	partbytes2, err := json.Marshal (&firstParticipant)
	if err != nil {
       fmt.Println("Could not marshal firstParticipant info object", err)
       return  err
	 }
	 err = stub.PutState("part1", partbytes2)
	 if err != nil {
       fmt.Println("Could not put updated firstParticipant in world state", err)
       return  err
	 }
	 
	return nil
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

	bytes, err := CreateParticipants(stub,args)
	if err != nil {
			logger.Error("Could not create and save participants to ledger", err)
			return nil, err
		}
	return bytes, nil
}

func (t *SampleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "GetLoanApplication" {
		return GetLoanApplication(stub, args)
	} else if function == "GetLoanParticipant" {
		return GetLoanParticipant(stub, args)
	} else if (function == "GetParticipatedLoans"){
		return GetParticipatedLoans(stub, args)
	}else {
		return nil, errors.New("Invalid function name")
	}
	
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
	if (function == "CreateLoanParticipation") {
		//username, _ := GetCertAttribute(stub, "username")
		//role, _ := GetCertAttribute(stub, "role")
		return CreateLoanParticipation(stub, args)
	/*	if role == "Bank_Admin" {
		return CreateLoanApplication(stub, args)
		} else {
			return nil, errors.New(username + " with role " + role + " does not have access to create a loan application")
		}*/

	} else if (function == "SettleLoanSyndication") {
		return SettleLoanSyndication(stub, args)
	} else {
		return nil, errors.New("Invalid function name")
	}
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
