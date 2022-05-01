/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package GoSDK

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"

	//mb "github.com/hyperledger/fabric-protos-go/msp"
	//pb "github.com/hyperledger/fabric-protos-go/peer"
	//"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	//"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"
	//"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl@v1.0.0-beta2"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

//初始化设置信息
type InitInfo struct {
	//组织名称 config.yaml ---> organizations ---> Org1
	Org1Name                 string              // 组织1名称
	Org1Admin                string              // 组织1管理员名称
	Org1AdminSigningIdentity msp.SigningIdentity //组织1管理员签名身份
	Org1User                 string              // 组织1普通用户名称

	Org2Name                 string              // 组织2名称
	Org2Admin                string              // 组织2管理员名称
	Org2AdminSigningIdentity msp.SigningIdentity //组织2管理员签名身份
	Org2User                 string              // 组织2普通用户名称

	OrdererOrgName  string //Orderer组织名称
	OrdererAdmin    string // Orderer 管理员名称
	OrdererEndpoint string //单独一个 Orderer名称

	ChannelID            string // 通道名称
	ChannelConfig        string // 通道配置文件所在路径
	Org1MSPanchorsConfig string // 组织1锚节点更新配置文件
	Org2MSPanchorsConfig string // 组织2锚节点更新配置文件
	SDKConfig            string // go sdk 配置文件

	ChaincodeID      string //链码名称/id
	ChaincodeGoPath  string //系统GOPATH路径
	ChaincodePath    string //链码所在路径
	ChaincodeVersion string //链码版本
}

// MSP 客户端
type MspClients struct {
	// Org MSP clients
	org1MspClient *mspclient.Client
	org2MspClient *mspclient.Client
}

// 资源管理上下文和客户端
type ResMgmtClients struct {
	// 资源管理客户端上下文
	ordererClientContext   contextAPI.ClientProvider //OrdererOrg客户端上下文
	org1AdminClientContext contextAPI.ClientProvider //Org1 客户端上下文
	org2AdminClientContext contextAPI.ClientProvider //Org2 客户端上下文

	// 通道资源管理客户端
	chResMgmtClient   *resmgmt.Client //ordererOrg通道资源管理客户端实例
	org1ResMgmtClient *resmgmt.Client //org1通道资源管理客户端实例
	org2ResMgmtClient *resmgmt.Client //org2通道资源管理客户端实例
}

// 通道上下文和客户端
type ChannelClients struct {
	// 通道客户端上下文
	//org1AdminChannelContext contextAPI.ChannelProvider//Org1 管理员 通道管理客户端上下文
	org1ChannelClientContext contextAPI.ChannelProvider //Org1 普通用户 通道管理客户端上下文
	org2ChannelClientContext contextAPI.ChannelProvider //Org2 普通用户 通道管理客户端上下文

	// 通道客户端
	Org1UserChClient *channel.Client
	Org2UserChClient *channel.Client
}

//实例信息
var (
//sdk   *fabsdk.FabricSDK   // SDK实例
//ccPkg *resource.CCPackage //链码包

// Peers
//org1Peers []fab.Peer
//org2Peers []fab.Peer
)

//创建SDK实例, 2个mspclient，3个资源管理客户端
func SetupSDK(info *InitInfo) (*fabsdk.FabricSDK, MspClients, ResMgmtClients, error) {
	var mspclients MspClients
	var res_mgmt_clients ResMgmtClients
	// 通过config.FromFile解析配置文件，然后通过fabsdk.New创建fabric go sdk的入口实例。
	sdk, err := fabsdk.New(config.FromFile(info.SDKConfig))
	if err != nil {
		return sdk, mspclients, res_mgmt_clients, fmt.Errorf("failed to create new SDK")
	}

	// 创建 Org1 MSP客户端 和 Org2 MSP客户端，用于获取签名身份
	mspclients.org1MspClient, err = mspclient.New(sdk.Context(), mspclient.WithOrg(info.Org1Name))
	if err != nil {
		return sdk, mspclients, res_mgmt_clients, fmt.Errorf("failed to create org1MspClient")
	}
	mspclients.org2MspClient, err = mspclient.New(sdk.Context(), mspclient.WithOrg(info.Org2Name))
	if err != nil {
		return sdk, mspclients, res_mgmt_clients, fmt.Errorf("failed to create org2MspClient")
	}

	// 创建资源管理客户端上下文
	res_mgmt_clients.ordererClientContext = sdk.Context(fabsdk.WithUser(info.OrdererAdmin), fabsdk.WithOrg(info.OrdererOrgName))
	res_mgmt_clients.org1AdminClientContext = sdk.Context(fabsdk.WithUser(info.Org1Admin), fabsdk.WithOrg(info.Org1Name))
	res_mgmt_clients.org2AdminClientContext = sdk.Context(fabsdk.WithUser(info.Org2Admin), fabsdk.WithOrg(info.Org2Name))

	//创建资源管理客户端，负责管理通道（创建/更新通道）
	res_mgmt_clients.chResMgmtClient, err = resmgmt.New(res_mgmt_clients.ordererClientContext)
	if err != nil {
		return sdk, mspclients, res_mgmt_clients, fmt.Errorf("failed to get a new channel management client")
	}
	// Org1资源管理上下文客户端，用于Org1锚节点更新通道
	res_mgmt_clients.org1ResMgmtClient, err = resmgmt.New(res_mgmt_clients.org1AdminClientContext)
	if err != nil {
		return sdk, mspclients, res_mgmt_clients, fmt.Errorf("failed to get a new channel management client for org1Admin")
	}
	// Org2资源管理上下文客户端，用于Org2锚节点更新通道
	res_mgmt_clients.org2ResMgmtClient, err = resmgmt.New(res_mgmt_clients.org2AdminClientContext)
	if err != nil {
		return sdk, mspclients, res_mgmt_clients, fmt.Errorf("failed to get a new channel management client for org2Admin")
	}

	fmt.Println("SetupSDK----创建实例、MSP客户端、资源管理客户端完成")

	return sdk, mspclients, res_mgmt_clients, nil
}

//创建通道
func CreatChannel(info *InitInfo, mspclients MspClients, res_mgmt_clients *ResMgmtClients) error {
	var err error
	// 获取签名的身份用来签名创建通道请求
	info.Org1AdminSigningIdentity, err = mspclients.org1MspClient.GetSigningIdentity(info.Org1Admin)
	if err != nil {
		return fmt.Errorf("failed to get org1AdminUser")
	}
	info.Org2AdminSigningIdentity, err = mspclients.org2MspClient.GetSigningIdentity(info.Org2Admin)
	if err != nil {
		return fmt.Errorf("failed to get org2AdminUser, err")
	}

	// 每次保存通道时查询一下上一次是否成功保存
	var lastConfigBlock uint64 //上一次的区块配置信息
	// 创建一个查询客户端
	configQueryClient, err := resmgmt.New(res_mgmt_clients.org1AdminClientContext)
	if err != nil {
		return fmt.Errorf("failed to get a new channel management client")
	}

	//使用Org1和Org2的签名身份 创建通道并为2个组织更新通道配置文件
	// 创建一个通道请求，根据通道ID，通道配置文件，两个组织的签名身份
	req := resmgmt.SaveChannelRequest{
		ChannelID:         info.ChannelID,
		ChannelConfigPath: info.ChannelConfig,
		SigningIdentities: []msp.SigningIdentity{info.Org1AdminSigningIdentity, info.Org2AdminSigningIdentity}}

	// 创建通道
	txID, err := res_mgmt_clients.chResMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererEndpoint))
	if err != nil {
		return err
	}

	// 等待至更新完成
	lastConfigBlock = WaitForOrdererConfigUpdate(configQueryClient, info.ChannelID, info.OrdererEndpoint, true, lastConfigBlock)

	//为Org1更新锚节点配置
	req = resmgmt.SaveChannelRequest{
		ChannelID:         info.ChannelID,
		ChannelConfigPath: info.Org1MSPanchorsConfig,
		SigningIdentities: []msp.SigningIdentity{info.Org1AdminSigningIdentity}}
	txID, err = res_mgmt_clients.org1ResMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererEndpoint))
	if err != nil {
		return fmt.Errorf("error should be nil for SaveChannel for Org1 anchor peer. txID = %s", txID)
	}

	// 等待至更新完成
	lastConfigBlock = WaitForOrdererConfigUpdate(configQueryClient, info.ChannelID, info.OrdererEndpoint, false, lastConfigBlock)

	//为Org2更新锚节点配置
	req = resmgmt.SaveChannelRequest{
		ChannelID:         info.ChannelID,
		ChannelConfigPath: info.Org2MSPanchorsConfig,
		SigningIdentities: []msp.SigningIdentity{info.Org2AdminSigningIdentity}}
	txID, err = res_mgmt_clients.org2ResMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererEndpoint))
	if err != nil {
		return fmt.Errorf("error should be nil for SaveChannel for Org2 anchor peer. txID = %s", txID)
	}

	// 等待至更新完成
	lastConfigBlock = WaitForOrdererConfigUpdate(configQueryClient, info.ChannelID, info.OrdererEndpoint, false, lastConfigBlock)
	fmt.Println("CreatChannel----创建通道、更新锚节点配置完成")
	return nil
}

// 等待状态 直到区块配置更新被提交
func WaitForOrdererConfigUpdate(client *resmgmt.Client, channelID string, OrdererEndpoint string, genesis bool, lastConfigBlock uint64) uint64 {

	blockNum, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
		func() (interface{}, error) {
			// 从Orderer查询通道配置信息
			chConfig, err := client.QueryConfigFromOrderer(channelID, resmgmt.WithOrdererEndpoint(OrdererEndpoint))
			if err != nil {
				return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), err.Error(), nil)
			}

			// 获取配置区块
			currentBlock := chConfig.BlockNumber()
			if !genesis && currentBlock <= lastConfigBlock {
				return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Block number was not incremented [%d, %d]", currentBlock, lastConfigBlock), nil)
			}

			block, err := client.QueryConfigBlockFromOrderer(channelID, resmgmt.WithOrdererEndpoint(OrdererEndpoint))
			if err != nil {
				return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), err.Error(), nil)
			}
			if block.Header.Number != currentBlock {
				return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Invalid block number [%d, %d]", block.Header.Number, currentBlock), nil)
			}

			return &currentBlock, nil
		},
	)
	if err != nil {
		fmt.Println("WaitForOrdererConfigUpdate 报错")
	}
	return *blockNum.(*uint64)
}

//加入通道
func JoinChannel(info InitInfo, res_mgmt_clients *ResMgmtClients) error {
	var err error
	// Org1 的peers 加入通道
	err = res_mgmt_clients.org1ResMgmtClient.JoinChannel(info.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererEndpoint))
	if err != nil {
		return err
	}

	// Org2 的peers 加入通道
	err = res_mgmt_clients.org2ResMgmtClient.JoinChannel(info.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererEndpoint))
	if err != nil {
		return err
	}
	fmt.Println("JoinChannel----Org1和Org2 peers加入通道成功")
	return nil
}

//打包链码
func PackageChaincode(info InitInfo, ) (*resource.CCPackage, error) {
	//创建链码包
	ccPkg, err := packager.NewCCPackage(info.ChaincodePath, info.ChaincodeGoPath)
	if err != nil {
		return ccPkg, fmt.Errorf("create ccPkg error")
	}
	fmt.Println("PackageChaincode----打包链码完成")
	return ccPkg, nil
}

//安装链码
func InstallChaincode(info InitInfo, res_mgmt_clients *ResMgmtClients, ccPkg *resource.CCPackage) ([]fab.Peer, []fab.Peer, error) {
	// 获取2个组织的节点
	org1Peers, err := DiscoverLocalPeers(res_mgmt_clients.org1AdminClientContext, 2)
	if err != nil {
		return org1Peers, nil, fmt.Errorf("获取Org1 peer节点失败")
	}
	org2Peers, err := DiscoverLocalPeers(res_mgmt_clients.org2AdminClientContext, 2)
	if err != nil {
		return org1Peers, org2Peers, fmt.Errorf("获取Org2 peer节点失败")
	}

	// 安装链码参数
	installCCReq := resmgmt.InstallCCRequest{Name: info.ChaincodeID, Path: info.ChaincodePath, Version: info.ChaincodeVersion, Package: ccPkg}

	// 在Org1 上安装链码
	_, err = res_mgmt_clients.org1ResMgmtClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return org1Peers, org2Peers, fmt.Errorf("InstallCC for Org1 failed")
	}

	// 在Org2 上安装链码
	_, err = res_mgmt_clients.org2ResMgmtClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return org1Peers, org2Peers, fmt.Errorf("InstallCC for Org2 failed")
	}

	// 确保Org1 上所有peer 安装了链码
	found := queryInstalledCC(res_mgmt_clients.org1ResMgmtClient, info.ChaincodeID, info.ChaincodeVersion, org1Peers)
	if !found {
		return org1Peers, org2Peers, fmt.Errorf("Org1 上不是所有peer都安装了链码")
	}
	// 确保Org2 上所有peer 安装了链码
	found = queryInstalledCC(res_mgmt_clients.org2ResMgmtClient, info.ChaincodeID, info.ChaincodeVersion, org2Peers)
	if !found {
		return org1Peers, org2Peers, fmt.Errorf("Org2 上不是所有peer都安装了链码")
	}
	fmt.Println("InstallChaincode----Org1 Org2安装链码完成")
	return org1Peers, org2Peers, nil
}

// 根据给出的MSP上下文，搜索本地Peers，返回所有的Peers。如果number不能匹配期望的个数，返回错误。
func DiscoverLocalPeers(ctxProvider contextAPI.ClientProvider, expectedPeers int) ([]fab.Peer, error) {
	ctx, err := context.NewLocal(ctxProvider)
	if err != nil {
		return nil, errors.Wrap(err, "error creating local context")
	}

	discoveredPeers, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
		func() (interface{}, error) {
			peers, serviceErr := ctx.LocalDiscoveryService().GetPeers()
			if serviceErr != nil {
				return nil, errors.Wrapf(serviceErr, "error getting peers for MSP [%s]", ctx.Identifier().MSPID)
			}
			if len(peers) < expectedPeers {
				return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Expecting %d peers but got %d", expectedPeers, len(peers)), nil)
			}
			return peers, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return discoveredPeers.([]fab.Peer), nil
}

//确保是否安装了链码
func queryInstalledCC(resMgmt *resmgmt.Client, ccName, ccVersion string, peers []fab.Peer) bool {
	installed, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
		func() (interface{}, error) {
			ok := isCCInstalled(resMgmt, ccName, ccVersion, peers)
			if !ok {
				return &ok, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Chaincode [%s:%s] is not installed on all peers in Org1", ccName, ccVersion), nil)
			}
			return &ok, nil
		},
	)
	if err != nil {
		fmt.Println("查询链码是否安装出现问题")
	}
	return *(installed).(*bool)
}

// 是否安装链码
func isCCInstalled(resMgmt *resmgmt.Client, ccName, ccVersion string, peers []fab.Peer) bool {
	installedOnAllPeers := true
	for _, peer := range peers {
		resp, err := resMgmt.QueryInstalledChaincodes(resmgmt.WithTargets(peer))
		if err != nil {
			fmt.Println("QueryInstalledChaincodes for peer failed")
			return false
		}
		found := false
		for _, ccInfo := range resp.Chaincodes {
			if ccInfo.Name == ccName && ccInfo.Version == ccVersion {
				found = true
				break
			}
		}
		if !found {
			installedOnAllPeers = false
		}
	}
	return installedOnAllPeers
}

//实例化链码,并创建通道客户端
func Instantiate(info InitInfo, sdk *fabsdk.FabricSDK, res_mgmt_clients *ResMgmtClients, org1Peers []fab.Peer, org2Peers []fab.Peer) (ChannelClients, error) {
	var channel_clients ChannelClients
	// 实例化链码，实例化只需要任意peer上执行1次
	// ****设置策略*****
	ccPolicy, err := cauthdsl.FromString("AND ('Org1MSP.member','Org2MSP.member')")
	if err != nil {
		return channel_clients, fmt.Errorf("error creating CC policy ")
	}
	instantiateCCReq := resmgmt.InstantiateCCRequest{
		Name:    info.ChaincodeID,
		Path:    info.ChaincodePath,
		Version: info.ChaincodeVersion,
		Args:    [][]byte{[]byte("init")},
		Policy:  ccPolicy}
	//参数设置方式：
	//[][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}
	_, err = res_mgmt_clients.org1ResMgmtClient.InstantiateCC(info.ChannelID, instantiateCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return channel_clients, err
	}

	// 确保Org1 上所有peer实例化链码
	found := queryInstantiatedCC(info.Org1Name, res_mgmt_clients.org1ResMgmtClient, info.ChannelID, info.ChaincodeID, info.ChaincodeVersion, org1Peers)
	if !found {
		return channel_clients, fmt.Errorf("Org1 上不是所有peer都实例化了链码")
	}
	fmt.Println("Instantiate----Org1 上peer实例化完成")
	// 确保Org2 上所有peer实例化链码
	found = queryInstantiatedCC(info.Org2Name, res_mgmt_clients.org2ResMgmtClient, info.ChannelID, info.ChaincodeID, info.ChaincodeVersion, org2Peers)
	if !found {
		return channel_clients, fmt.Errorf("Org2 上不是所有peer都实例化了链码")
	}
	fmt.Println("Instantiate----Org2 上peer实例化完成")
	// 创建通道上下文，用于创建通道客户端
	//org1AdminChannelContext = sdk.ChannelContext(mc.channelID, fabsdk.WithUser(info.Org1Admin), fabsdk.WithOrg(info.Org1Name))
	channel_clients.org1ChannelClientContext = sdk.ChannelContext(info.ChannelID, fabsdk.WithUser(info.Org1User), fabsdk.WithOrg(info.Org1Name))
	channel_clients.org2ChannelClientContext = sdk.ChannelContext(info.ChannelID, fabsdk.WithUser(info.Org2User), fabsdk.WithOrg(info.Org2Name))

	// 创建通道客户端。通道客户端可以查询链码，执行链码。
	// 创建Org1 的通道客户端。
	channel_clients.Org1UserChClient, err = channel.New(channel_clients.org1ChannelClientContext)
	if err != nil {
		return channel_clients, fmt.Errorf("Failed to create new channel client for Org1 user: %s", err)
	}
	// 创建Org2 的通道客户端。
	channel_clients.Org2UserChClient, err = channel.New(channel_clients.org2ChannelClientContext)
	if err != nil {
		return channel_clients, fmt.Errorf("Failed to create new channel client for Org2 user: %s", err)
	}
	fmt.Println("Instantiate----实例化链码完成，Org1 Org2 通道客户端创建完成")
	return channel_clients, nil
}

//确保组织上的所有节点都已经实例化链码
func queryInstantiatedCC(orgID string, resMgmt *resmgmt.Client, channelID, ccName, ccVersion string, peers []fab.Peer) bool {
	instantiated, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
		func() (interface{}, error) {
			ok := isCCInstantiated(resMgmt, channelID, ccName, ccVersion, peers)
			if !ok {
				return &ok, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Did NOT find instantiated chaincode [%s:%s] on one or more peers in [%s].", ccName, ccVersion, orgID), nil)
			}
			return &ok, nil
		},
	)
	if err != nil {
		fmt.Println("Got error checking if chaincode was instantiated")
	}
	return *(instantiated).(*bool)
}

//是否实例化
func isCCInstantiated(resMgmt *resmgmt.Client, channelID, ccName, ccVersion string, peers []fab.Peer) bool {
	InstantiatedOnAllPeers := true
	for _, peer := range peers {
		chaincodeQueryResponse, err := resMgmt.QueryInstantiatedChaincodes(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargets(peer))
		if err != nil {
			fmt.Println("QueryInstantiatedChaincodes return error")
			return false
		}
		found := false
		for _, chaincode := range chaincodeQueryResponse.Chaincodes {
			if chaincode.Name == ccName && chaincode.Version == ccVersion {
				found = true
				break
			}
		}
		if !found {
			InstantiatedOnAllPeers = false
		}
	}
	return InstantiatedOnAllPeers
}
