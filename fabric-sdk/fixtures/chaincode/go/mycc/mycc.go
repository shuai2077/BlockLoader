package main
import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"strconv"
)
type AuChaincode struct {
}
func (s *AuChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("init succeed"))
}
func (s *AuChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fn, args := stub.GetFunctionAndParameters()
	var result string
	var err error
	if fn == "initAccount" {
		result, err = initAccount(stub, args)
	} else if fn == "transfer" {
		err = transfer(stub, args)
	} else if fn == "queryAccount" {
		result, err = queryAccount(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(result))
}
func initAccount(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return "",fmt.Errorf("Failed to initAccount asset: %s", err)
	}
	return "Success", nil
}
func transfer(stub shim.ChaincodeStubInterface, args []string) (error) {
	var acount_1 int
	var acount_2 int
	transfer, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("Failed to Atoi asset: %s", err)
	}
	asset, err := stub.GetState(args[0])
	if err != nil {
		return fmt.Errorf("Failed to set asset: %s", err)
	}
	acount_1, err = strconv.Atoi(string(asset))
	if err != nil {
		return fmt.Errorf("Failed to Atoi asset: %s", err)
	}
	asset, err = stub.GetState(args[0])
	if err != nil {
		return fmt.Errorf("Failed to set asset: %s", err)
	}
	acount_2, err = strconv.Atoi(string(asset))
	if err != nil {
		return fmt.Errorf("Failed to Atoi asset: %s", err)
	}
	if transfer > acount_1{
		return fmt.Errorf("account balances are insufficient")
	}
	acount_2 = acount_2 + transfer
	acount_1 = acount_1 - transfer
	err = stub.PutState(args[0], []byte(strconv.Itoa(acount_1)))
	if err != nil {
		return fmt.Errorf("Failed to set asset: %s", err)
	}
	err = stub.PutState(args[1], []byte(strconv.Itoa(acount_2)))
	if err != nil {
		return fmt.Errorf("Failed to set asset: %s", err)
	}
	return nil
}
func queryAccount(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	asset, err := stub.GetState(args[0])
	if err != nil {
		return "",fmt.Errorf("Failed to set asset: %s", err)
	}
	return string(asset), nil
}
func main() {
	err := shim.Start(new(AuChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}