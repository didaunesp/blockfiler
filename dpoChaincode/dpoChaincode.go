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

type RegisterPublic struct {
	Key  string `json:"key"`
	User string `json:"user"`
	Time string `json:"time"`
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
	} else if function == "getDpoKeys" {
		return s.getDpoKeys(APIstub, args)
	} else if function == "updateRegister" {
		return s.updateRegister(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) getDpoKeys(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	QueryAsBytes, err := APIstub.GetState("DPO")
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(QueryAsBytes)
}

func (s *SmartContract) query(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := args[0]

	collection, QueryAsBytes, err := s.doQueryRegister(APIstub, key)
	fmt.Println(collection)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(QueryAsBytes)
}

func (s *SmartContract) doQueryRegister(APIstub shim.ChaincodeStubInterface, key string) (string, []byte, error) {
	collection := "collectionPublico"
	QueryAsBytes, err := APIstub.GetPrivateData(collection, key)
	if len(QueryAsBytes) == 0 {
		collection = "collectionReativo"
		QueryAsBytes, err = APIstub.GetPrivateData(collection, key)
		if len(QueryAsBytes) == 0 {
			collection = "collectionAtivo"
			QueryAsBytes, err = APIstub.GetPrivateData(collection, key)
			if len(QueryAsBytes) == 0 {
				collection = "collectionEmpresa"
				QueryAsBytes, err = APIstub.GetPrivateData(collection, key)
			}
		}
	}
	return collection, QueryAsBytes, err
}

func (s *SmartContract) updateRegister(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := args[0]
	user := args[1]

	QueryAsBytes, err := APIstub.GetState(key)

	if err != nil {
		return shim.Error(err.Error())
	}

	if len(QueryAsBytes) > 0 {
		var register RegisterPublic
		json.Unmarshal(QueryAsBytes, &register)
		register.User = user
		register.Time = time.Now().Format("2006-01-02 15:04:05")

		registerAsBytes, err2 := json.Marshal(register)
		if err2 != nil {
			return shim.Error(err2.Error())
		}
		err3 := APIstub.PutState(key, registerAsBytes)
		if err3 != nil {
			return shim.Error(err3.Error())
		}

		fmt.Println("key: " + key)
		fmt.Println("register: " + string(registerAsBytes))
	}
	return shim.Success(nil)
}

func (s *SmartContract) history(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	resultsIterator, err := APIstub.GetHistoryForKey(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("args, results")
	fmt.Println(args[0])
	fmt.Println(resultsIterator)

	resultsIterator.Close()
	fmt.Println(resultsIterator)
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err2 := resultsIterator.Next()
		fmt.Println("response")
		fmt.Println(response)
		if err2 != nil {
			return shim.Error(err2.Error())
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
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	fmt.Println("key ", args[0])
	var key = args[0]
	var content = args[1]
	var collection = args[2]
	var user = "DPOcli1"

	register := &Register{key, content, user, time.Now().Format("2006-01-02 15:04:05")}
	registerPublic := &RegisterPublic{key, user, time.Now().Format("2006-01-02 15:04:05")}
	registerAsBytes, err := json.Marshal(register)
	registerPublicAsBytes, err := json.Marshal(registerPublic)
	if err != nil {
		return shim.Error(err.Error())
	}
	registerPublicAsBytes, errPublic := json.Marshal(register)
	if errPublic != nil {
		return shim.Error(errPublic.Error())
	}

	fmt.Println("registerAsBytes ", string(registerAsBytes))
	fmt.Println("collection", collection)

	err2 := APIstub.PutPrivateData(collection, key, registerAsBytes)
	if err2 != nil {
		return shim.Error(err2.Error())
	}

	err2Public := APIstub.PutState(key, registerPublicAsBytes)
	if err2Public != nil {
		return shim.Error(err2Public.Error())
	}

	return s.updateDpoList(APIstub, key)
}

func (s *SmartContract) updateDpoList(APIstub shim.ChaincodeStubInterface, key string) sc.Response {
	keysBytes, err := APIstub.GetState("DPO")
	if err != nil {
		return shim.Error(err.Error())
	}
	keysString := string(keysBytes)
	keys := keysString + key + "|"
	err2 := APIstub.PutState("DPO", []byte(keys))
	if err2 != nil {
		return shim.Error(err.Error())
	}
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
