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
	address       string `json:"address"`
	Email     string `json:"email"`
	contact    string `json:"contact"`
}

type FinancialInfo struct {
	spRating      string `json:"spRating"`
	moodyRating        string `json:"moodyRating"`
	dcr   			   int `json:"dcr"`
	turnover	   int `json:"turnover"`
}

type LoanApplication struct {
	ID                     string        `json:"id"`
	DealType			   string 		 `json:"dealType"`
	BaseRateType		   string        `json:"baseRateType"`
	AllInRate			   int			 `json:"allInRate"`
	Spread				   int			 `json:"spread"`
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
	DealAmount         int          	 `json:"dealAmount"`
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
	SharePerCent           int           `json:"share"`
	AssetList []Asset
}
type Asset struct{
	AssetId								 string        `json:"loanId"`
	
	ShareAmount            int 					 `json:"shareAmount"`
	SyndicatedAmount 			 int					 `json:"syndicatedAmount"`
	SettlementFees		   int                     `json:"settlementFees"`
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
	
	var secondParticipant Participant

	/*secondParticipant = Participant{ID:"part1",Name:"DeucheBank",
								AssetList: []Asset{
									{AssetId:"la1",
									SharePerCent:80,
									ShareAmount:0,
									SyndicatedAmount:1000},
								},
						}*/
	
	firstParticipant = Participant{ID:participantID ,Name:"DeucheBank",SharePerCent:80}
	
	secondParticipant = Participant{ID:"part2",Name:"CitiBank",SharePerCent:20}

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

		
	bytes2, err2 := json.Marshal (&secondParticipant)
	 if err2 != nil {
		         fmt.Println("Could not marshal secondParticipant object", err2)
			         return nil, err2
				  }

	err3 := stub.PutState("part2", bytes2 )
	if err3 != nil {
			logger.Error("Could not save secondParticipant to ledger", err3)
			return nil, err3
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
    fmt.Println("CreateLoanParticipation : ParticipatedLoan ID and amount " + participatedLoan.ID, participatedLoan.DealAmount)
	
	fmt.Println("CreateLoanParticipation : PropertyId " + participatedLoan.PropertyId)
   
    fmt.Println("CreateLoanParticipation:  All In Rate ", participatedLoan.AllInRate)
	
	fmt.Println("CreateLoanParticipation : baseRateType " + participatedLoan.BaseRateType)


	loanbytes2, err := AppendToLoanList(stub,participatedLoan)
	    
    err = ParticipateLoan(stub, "part1",loanApplicationId, participatedLoan.DealAmount)
	err = ParticipateLoan(stub, "part2",loanApplicationId, participatedLoan.DealAmount)

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
		logger.Error("Could not fetch firstParticipant from ledger", err)
		return  err
	}

	var firstParticipant Participant
	err = json.Unmarshal(partbytes,&firstParticipant)
	if err != nil {
		return err
	}
	fmt.Println("ParticipateLoan: firstParticipant Name" + firstParticipant.Name)
	fmt.Println("ParticipateLoan: participationAmount" ,participationAmount)
		fmt.Println("ParticipateLoan: SharePerCent" ,firstParticipant.SharePerCent)
	

	var newAsset Asset
	newAsset.AssetId = loan_id
	//newAsset.SharePerCent = 80
	newAsset.ShareAmount = ( participationAmount * firstParticipant.SharePerCent / 100 )
	newAsset.SettlementFees = 0
	
	firstParticipant.AssetList = append(firstParticipant.AssetList, newAsset)
	
	fmt.Println("Total loans participated")
	fmt.Println(len(firstParticipant.AssetList))
	
	 partbytes2, err := json.Marshal (&firstParticipant)
	 if err != nil {
        fmt.Println("Could not marshal firstParticipant info object", err)
        return err
	 }
	err = stub.PutState(participant, partbytes2)
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
    fmt.Println("SettleLoanSyndication : participatedLoan ID and amount " + participatedLoan.ID, participatedLoan.DealAmount)
	
	fmt.Println("SettleLoanSyndication: All In Rate ", participatedLoan.AllInRate)
	
	fmt.Println("SettleLoanSyndication :  baseRateType " + participatedLoan.BaseRateType)

	fmt.Println("SettleLoanSyndication : updating outStandingSettlentAmount for ID for amount " + loanSettlementAmount)

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

	err = SettleParticipation(stub,"part1",loanApplicationId,participatedLoan.AllInRate,v)
	
	var customEvent = "{eventType: 'loanApplicationCreation', description:" + loanApplicationId + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return nil, err
	}
	logger.Info("Successfully saved loan application")
	return laBytes, nil

return nil,nil
}

func SettleParticipation(stub shim.ChaincodeStubInterface, participant string, loan_id string , allinRate int,  settlementAmount int) (error){
	fmt.Println("Entering SettleParticipation")
	partbytes, err := stub.GetState(participant)
	if err != nil || partbytes == nil {
		logger.Error("Could not fetch firstParticipant with id part1 from ledger", err)
		return err
	}

	 var firstParticipant Participant
	 err = json.Unmarshal(partbytes,&firstParticipant)
	 fmt.Println("SettleParticipation: firstParticipantName " + firstParticipant.Name)

	for i := range firstParticipant.AssetList {
		fmt.Println("elem.AssetId" , firstParticipant.AssetList[i].AssetId)
		fmt.Println("loan_id",loan_id)
		if(firstParticipant.AssetList[i].AssetId == loan_id){
			fmt.Println("SettleParticipation:Participant Asset Details :" + firstParticipant.AssetList[i].AssetId)
			fmt.Println("SettleParticipation:Index value is :" , i )
			var settlementPortion int
			//settlementPortion = firstParticipant.AssetList[i].SharePerCent*settlementAmount/100
			settlementPortion = firstParticipant.SharePerCent*settlementAmount/100
			fmt.Println("SettleParticipation:settlementPortion Portion :", settlementPortion)
			var orginalShareAmt int
			orginalShareAmt = firstParticipant.AssetList[i].ShareAmount
			fmt.Println("SettleParticipation:orginalShareAmt" , orginalShareAmt)
			firstParticipant.AssetList[i].SettlementFees = firstParticipant.AssetList[i].SettlementFees + ((orginalShareAmt*30*allinRate)/(100*365))
			firstParticipant.AssetList[i].ShareAmount = orginalShareAmt - settlementPortion
			
			fmt.Println("SettleParticipation:Update Participant ShareAmount")
			fmt.Println(firstParticipant.AssetList[i].ShareAmount)
		}
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
	 	fmt.Println("Exiting SettleParticipation")
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
