package main

import (
	"fabric-sdk-go/sdkInit"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"

	"github.com/hyperledger/fabric/protoutil"
	"os"
)

const (
	cc_name    = "smallbank"
	cc_version = "1.0.0"
)

var App sdkInit.Application

func main() {
	// init orgs information

	//org信息
	orgs := []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: "/root/fabric-sdk/fixtures/channel-artifacts/Org1MSPanchors.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org2",
			OrgMspId:      "Org2MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: "/root/fabric-sdk/fixtures/channel-artifacts/Org2MSPanchors.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org3",
			OrgMspId:      "Org3MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: "/root/fabric-sdk/fixtures/channel-artifacts/Org3MSPanchors.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org4",
			OrgMspId:      "Org4MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: "/root/fabric-sdk/fixtures/channel-artifacts/Org4MSPanchors.tx",
		},
	}
	// 初始化info
	info := sdkInit.SdkEnvInfo{
		ChannelID:        "mychannel",
		ChannelConfig:    "/root/fabric-sdk/fixtures/channel-artifacts/channel.tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer1.example.com",
		ChaincodeID:      cc_name,
		ChaincodePath:    "/root/fabric-sdk/chaincode/go/smallbank",
		ChaincodeVersion: cc_version,
	}

	// sdk setup
	sdk, err := sdkInit.Setup("config.yaml", &info)
	if err != nil {
		fmt.Println(">> SDK setup error:", err)
		os.Exit(-1)
	}

	// create channel and join
	if err := sdkInit.CreateAndJoinChannel(&info); err != nil {
		fmt.Println(">> Create channel and join error:", err)
		os.Exit(-1)
	}

	// create chaincode lifecycle
	if err := sdkInit.CreateCCLifecycle(&info, 1, false, sdk); err != nil {
		fmt.Println(">> create chaincode lifecycle error: %v", err)
		os.Exit(-1)
	}

	// invoke chaincode set status
	fmt.Println(">> 通过链码外部服务设置链码状态......")

	if err := info.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk); err != nil {

		fmt.Println("InitService successful")
		os.Exit(-1)
	}

	App = sdkInit.Application{
		SdkEnvInfo: &info,
	}
	fmt.Println(">> 设置链码状态完成")

	bereg, notifier := sdkInit.BlockListener(info.EvClient)

	for {
		e := <-notifier
		for _, data := range e.Block.Data.Data {
			// 解析区块中的每一笔交易
			env, err := protoutil.UnmarshalEnvelope(data)
			if err != nil {
				fmt.Printf("failed to unmarshal envelop:%s\n", err)
				continue
			}

			payload, err := protoutil.UnmarshalPayload(env.Payload)
			if err != nil {
				fmt.Printf("failed to unmarshal payload:%s\n", err)
				continue
			}

			channelHeader, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
			if err != nil {
				fmt.Printf("failed to unmarshal channel header:%s\n", err)
				continue
			}

			txID := channelHeader.TxId
			timestamp := channelHeader.Timestamp

			writeToFile(txID, timestamp)
		}
	}

	defer info.EvClient.Unregister(bereg)
	//defer info.EvClient.Unregister(sdkInit.ChainCodeEventListener(info.EvClient, info.ChaincodeID))
}

func writeToFile(txID string, timestamp *timestamppb.Timestamp) {
	// 创建或打开文件
	file, err := os.OpenFile("transaction_info.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// 格式化时间戳
	timeFormatted := time.Unix(timestamp.Seconds, int64(timestamp.Nanos)).Format(time.RFC3339)

	// 格式化字符串输出
	content := fmt.Sprintf("Transaction ID: %s\nTimestamp: %s\n", txID, timeFormatted)

	// 将内容写入文件
	if _, err := file.WriteString(content); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
}
