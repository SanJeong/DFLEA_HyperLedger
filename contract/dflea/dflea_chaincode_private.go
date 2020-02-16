package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type User struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Mail       string `json:"mail"`
	Name 	   string  `json:"name"`
	Rank	   int    `json:"rank"`
	Historys   []History `json:"history"`
}

type History struct{
	PDataNum	  string `json:"pdataNum"`
	PDate  string `json:"pdate"`
}

type Dataset struct {
	ObjectType string `json:"docType"`
	DataNum	  string `json:"dataNum"`
	Mail       string `json:"mail"`
	DataName  string `json:"dataName"`
	UpDate	  string `json:"update"`
}

type Purchase struct {
	ObjectType string `json:"docType"`
	PDataNum	  string `json:"pdataNum"`
	Mail       string `json:"mail"`
	PDataName  string `json:"pdataName"`
	PDate  string `json:"pdate"`
	PStatus string `json:"pstatus"`
}

type datasetPrivateDetails struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	DataNum    string `json:"dataNum"`    //the fieldtags are needed to keep case from bouncing around
	Datapath   string `json:"datapath"`
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	switch function {
	case "addUser":
		//create a new marble
		return t.addUser(stub, args)
	case "addDataset":
		//read a marble
		return t.addDataset(stub, args)
	case "addPurchase":
		return t.addPurchase(stub,args)
	case "readUser":
		//read a marble private details
		return t.readUser(stub, args)
	case "readDataset":
		//read a marble private details
		return t.readDataset(stub, args)
	case "readDatasetPrivateDetails":
		//change owner of a specific marble
		return t.readDatasetPrivateDetails(stub, args)
	case "readPurchase":
		return t.readPurchase(stub, args)
	case "transferPurchase":
		return t.transferPurchase(stub, args)
	default:
		//error
		fmt.Println("invoke did not find func: " + function)
		return shim.Error("Received unknown function invocation")
	}
}



// ============================================================
// add User
// ============================================================


func (t *SimpleChaincode) addUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("fail!. Incorrect number of argu")
	}
	var user = User{ObjectType: "user", Mail: args[0], Name: args[1], Rank: 0}
	userAsBytes, _ := json.Marshal(user)
	stub.PutState(args[0], userAsBytes)

	return shim.Success(nil)
}

func (t *SimpleChaincode) addDataset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		return shim.Error("fail!. Incorrect number of argu")
	}

	var dataset = Dataset{ObjectType: "dataset", DataNum: args[0], Mail: args[1], DataName: args[2], UpDate: args[3]}
	var datasetPrivateDetails = datasetPrivateDetails{ObjectType: "datasetPrivateDetails", DataNum: args[0], Datapath: args[4]}

	// ==== Check if already exists ====
	datasetAsBytes, err := stub.GetPrivateData("collectionDatasets", args[0])
	if err != nil {
		return shim.Error("Failed to get dataset: " + err.Error())
	} else if datasetAsBytes != nil {
		fmt.Println("This dataset already exists: " + args[0])
		return shim.Error("This dataset already exists: " + args[0])
	}

	datasetJSONasBytes, err := json.Marshal(dataset)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("new dataset number: " + args[0])
	err = stub.PutPrivateData("collectionDatasets", args[0], datasetJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	datasetPrivateDetailsBytes, err := json.Marshal(datasetPrivateDetails)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutPrivateData("collectionDatasetPrivateDetails", args[0], datasetPrivateDetailsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//add rank counting
	userAsBytes, err := stub.GetState(args[1])
	if err != nil{
		jsonResp := "\"Error\":\"Failed to get state for "+ args[0]+"\"}"
		return shim.Error(jsonResp)
	} else if userAsBytes == nil{ // no State! error
		jsonResp := "\"Error\":\"User does not exist: "+ args[0]+"\"}"
		return shim.Error(jsonResp)
	}

	user := User{}
	err = json.Unmarshal(userAsBytes, &user)
	if err != nil {
		return shim.Error(err.Error())
	}

	user.Rank =  (user.Rank + 1)
	userAsBytes, err = json.Marshal(user);
	stub.PutState(args[1], userAsBytes)

	fmt.Println("- end init dataset")
	return shim.Success(nil)
}

func (t *SimpleChaincode) addPurchase(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("fail!. Incorrect number of argu")
	}

	var purchase = Purchase{ObjectType: "purchase", PDataNum: args[0], Mail: args[1], PDate: args[3], PDataName: args[2], PStatus: "Pending"}

	// ==== Check if already exists ====
	purchaseAsBytes, err := stub.GetPrivateData("collectionPurchases", args[0])
	if err != nil {
		return shim.Error("Failed to get purchase: " + err.Error())
	} else if purchaseAsBytes != nil {
		fmt.Println("This purchase already exists: " + args[0])
		return shim.Error("This purchase already exists: " + args[0])
	}

	purchaseJSONasBytes, err := json.Marshal(purchase)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("new purchase number: " + args[0])
	err = stub.PutPrivateData("collectionPurchases", args[0], purchaseJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//history add
	userAsBytes, err := stub.GetState(args[1])
	if err != nil{
		jsonResp := "\"Error\":\"Failed to get state for "+ args[0]+"\"}"
		return shim.Error(jsonResp)
	} else if userAsBytes == nil{ // no State! error
		jsonResp := "\"Error\":\"User does not exist: "+ args[0]+"\"}"
		return shim.Error(jsonResp)
	}

	user := User{}
	err = json.Unmarshal(userAsBytes, &user)
	if err != nil {
		return shim.Error(err.Error())
	}

	var History = History{PDataNum: args[0], PDate: args[3]}
	user.Historys = append(user.Historys, History)
	userAsBytes, err = json.Marshal(user);

	stub.PutState(args[1], userAsBytes)

	fmt.Println("- end init purchase")
	return shim.Success(nil)
}

func (t *SimpleChaincode) readUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	UserAsBytes, _ := stub.GetState(args[0])
	return shim.Success(UserAsBytes)
}

func (t *SimpleChaincode) readDataset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the dataset to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateData("collectionDatasets", name) //get the marble from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\" does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ===============================================
// readMarblereadMarblePrivateDetails - read a marble private details from chaincode state
// ===============================================
func (t *SimpleChaincode) readDatasetPrivateDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the dataset to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateData("collectionDatasetPrivateDetails", name) //get the marble private details from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get private details for " + name + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"private details does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

func (t *SimpleChaincode) readPurchase(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the dataset to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateData("collectionPurchases", name) //get the marble from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

func (t *SimpleChaincode) transferPurchase(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start transfer purchase")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}

	Name := args[0]
	purchaseAsBytes, err := stub.GetPrivateData("collectionPurchases", Name)
	if err != nil {
		return shim.Error("Failed to get purchase:" + err.Error())
	} else if purchaseAsBytes == nil {
		return shim.Error("does not exist: " + Name)
	}
	PurchaseToTransfer := Purchase{}
	err = json.Unmarshal(purchaseAsBytes, &PurchaseToTransfer)
	if err != nil {
		return shim.Error(err.Error())
	}

	PStatus := args[1]
	PurchaseToTransfer.PStatus = PStatus

	purchaseJSONasBytes, _ := json.Marshal(PurchaseToTransfer)
	err = stub.PutPrivateData("collectionPurchases", PurchaseToTransfer.PDataNum, purchaseJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end transferpurchase (success)")
	return shim.Success(nil)
}


