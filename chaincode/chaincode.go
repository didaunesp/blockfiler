package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

type Register struct {
	Key     string `json:"key"`
	Content string `json:"content"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately

	if function == "create" {
		return s.create(APIstub, args)
	} else if function == "query" {
		return s.query(APIstub, args)
	} else if function == "history" {
		return s.history(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) create(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var err error

	fmt.Println("content ", args[0])
	var key = args[0]

	keyAsBytes, _ := json.Marshal(args[1])
	fmt.Println("keyAsBytes ", keyAsBytes)

	err = APIstub.PutState(key, keyAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SmartContract) query(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	QueryAsBytes, _ := APIstub.GetHistoryForKey(args[0])
	return shim.Success(QueryAsBytes)
}

func (s *SmartContract) history(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	QueryAsBytes, _ := APIstub.GetHistoryForKey(args[0])
	return shim.Success(QueryAsBytes)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
