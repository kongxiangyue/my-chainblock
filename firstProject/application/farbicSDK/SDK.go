package farbicSDK

import (
	"os"

	// 1. 导入关键包
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var (
	sdk           *fabsdk.FabricSDK
	channelName   = "mychannel"
	userName      = "Admin"
	orgName       = "Org1"
	chaincodeName = "mycc"
	configPath    = os.Getenv("HOME") + "/firstProject/application/config_e2e.yaml"
)

func Init() {
	var err error
	// 1. 初始化SDK
	sdk, err = fabsdk.New(config.FromFile(configPath))
	if err != nil {
		panic(err.Error())
	}
}

func createClient() (*channel.Client, error) {
	// 2. 生成通道配置文件
	context := sdk.ChannelContext(channelName,
		fabsdk.WithUser(userName),
		fabsdk.WithOrg(orgName),
	)
	// 3. 根据配置文件生成通道
	cli, err := channel.New(context)
	if err != nil {
		return new(channel.Client), err
	}
	return cli, nil
}

func ChannelQuery(fcn string, args [][]byte) ([]byte, error) {
	// 4. 与链码层交互
	cli, err := createClient()
	if err != nil {
		return []byte{}, err
	}
	resp, err := cli.Query(channel.Request{
		ChaincodeID: chaincodeName,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithTargetEndpoints("peer0.Organization1.test.com", "node1.Organization2.test.com"))
	if err != nil {
		return []byte{}, err
	}
	return resp.Payload, nil
}

func ChannelExecute(fcn string, args [][]byte) ([]byte, error) {
	// 4. 与链码层交互
	cli, err := createClient()
	if err != nil {
		return []byte{}, err
	}
	resp, err := cli.Execute(channel.Request{
		ChaincodeID: chaincodeName,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithTargetEndpoints("peer0.Organization1.test.com", "node1.Organization2.test.com"))
	if err != nil {
		return []byte{}, err
	}
	return resp.Payload, nil
}
