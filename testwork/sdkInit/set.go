package sdkInit

import "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"

func (t *Application) Set(args []string) (string, error) {
	request := channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}}
	response, err := t.SdkEnvInfo.Client.Execute(request)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}
func (t *Application) SetV(args []string) (string, error) {
	request := channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}}
	response, err := t.SdkEnvInfo.Client.Execute(request)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}
func (t *Application) SetVV(args []string) (string, error) {
	request := channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]),[]byte(args[4])}}
	response, err := t.SdkEnvInfo.Client.Execute(request)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}