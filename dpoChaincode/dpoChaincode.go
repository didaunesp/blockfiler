package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

type Register struct {
	Key     string `json:"key"`
	Content string `json:"content"`
	User    string `json:"user"`
	Time    string `json:"time"`
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

func (s *SmartContract) query(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	QueryAsBytes, _ := APIstub.GetState(args[0])

	// Unmarshal string into structs.
	var registers []Register
	json.Unmarshal(QueryAsBytes, &registers)

	// Loop over structs and update user and time
	for r := range registers {
		registers[r].User = "user"
		registers[r].Time = time.Now().String()
		registerAsBytes, err := json.Marshal(registers[r])
		if err != nil {
			return shim.Error(err.Error())
		}
		err2 := APIstub.PutState(registers[r].Key, registerAsBytes)
		if err2 != nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success(QueryAsBytes)
}

func (s *SmartContract) history(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	resultsIterator, _ := APIstub.GetHistoryForKey(args[0])

	resultsIterator.Close()
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForMarble returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) create(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var err error

	fmt.Println("key ", args[0])
	var key = args[0]
	var content = args[1]

	// keyAsBytes, _ := json.Marshal(args[1])
	// fmt.Println("contentAsBytes ", keyAsBytes)
	register := &Register{key, content, "user", time.Now().String()}
	registerAsBytes, err := json.Marshal(register)
	if err != nil {
		return shim.Error(err.Error())
	}

	err2 := APIstub.PutState(key, registerAsBytes)
	if err2 != nil {
		return shim.Error(err.Error())
	}
	// err2 := stub.PutPrivateData("collectionAtivo", key, registerAsBytes)
	// if err2 != nil {
	// 	return shim.Error(err.Error())
	// }

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
