package main

import (
	"fmt"
	"github.com/KangChain.com/KangChain/GoSDK"
	"github.com/KangChain.com/KangChain/service"
	"github.com/KangChain.com/KangChain/web"
	"github.com/KangChain.com/KangChain/web/controller"
	"os"
)

func main() {
	/*
		1. 配置
	*/
	// 初始化Go SDK 配置
	initSDKInfo := &GoSDK.InitInfo{
		Org1Name:             "Org1",
		Org1Admin:            "Admin",
		Org1User:             "User1",
		Org2Name:             "Org2",
		Org2Admin:            "Admin",
		Org2User:             "User1",
		OrdererOrgName:       "OrdererOrg",
		OrdererAdmin:         "Admin",
		OrdererEndpoint:      "orderer.example.com",
		ChannelID:            "mychannel",
		ChannelConfig:        os.Getenv("GOPATH") + "/src/github.com/KangChain.com/KangChain/demo0/channel-artifacts/channel.tx",
		Org1MSPanchorsConfig: os.Getenv("GOPATH") + "/src/github.com/KangChain.com/KangChain/demo0/channel-artifacts/Org1MSPanchors.tx",
		Org2MSPanchorsConfig: os.Getenv("GOPATH") + "/src/github.com/KangChain.com/KangChain/demo0/channel-artifacts/Org2MSPanchors.tx",
		SDKConfig:            "./demo0/config.yaml",
		ChaincodeID:          "simplecc",
		ChaincodeGoPath:      os.Getenv("GOPATH"),
		ChaincodePath:        "github.com/KangChain.com/KangChain/chaincode/",
		ChaincodeVersion:     "3"}
	//fmt.Printf("%+v",initSDKInfo)
	// 创建SDK及msp客户端，资源管理客户端
	sdk, mspClients, resMgmtClients, err := GoSDK.SetupSDK(initSDKInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	// return前自动释放资源
	defer sdk.Close()
	/*
		2. 生命周期管理
	*/
	//创建通道
	err = GoSDK.CreatChannel(initSDKInfo, mspClients, &resMgmtClients)
	if err != nil {
		fmt.Println(err)
		return
	}
	//加入通道
	err = GoSDK.JoinChannel(*initSDKInfo, &resMgmtClients)
	if err != nil {
		fmt.Println(err)
		return
	}
	//打包链码
	ccPkg, err := GoSDK.PackageChaincode(*initSDKInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	//安装链码
	org1Peers, org2Peers, err := GoSDK.InstallChaincode(*initSDKInfo, &resMgmtClients, ccPkg)
	if err != nil {
		fmt.Println(err)
		return
	}
	//实例化链码
	channelClients, err := GoSDK.Instantiate(*initSDKInfo, sdk, &resMgmtClients, org1Peers, org2Peers)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(channelClients)
	/*
		3. web设置
	*/
	//serviceSetup0 := service.ServiceSetup{
	//	ChaincodeID: initSDKInfo.ChaincodeID,
	//	Client:      channelClients.Org2UserChClient,
	//}
	//
	//msg, err := serviceSetup0.SetInfo("test1", "result1")
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(msg)
	//}

	serviceSetup := service.ServiceSetup{
		ChaincodeID: initSDKInfo.ChaincodeID,
		Client:      channelClients.Org1UserChClient,
	}
	app := controller.Application{
		Fabric: &serviceSetup,
	}
	web.WebStart(&app)
}
