package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	err := stub.PutState("lzp", []byte("ComeOn"))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	f, args := stub.GetFunctionAndParameters()
	if f != "invoke" {
		return shim.Error("Unknown function call")
	}
	fun := args[0]
	args = args[1:]
	switch fun {
	case "set":
		{
			return t.set(stub, args)
		}
	case "get":
		{
			return t.get(stub, args)
		}
	case "setPatient":
		{
			return t.setPatient(stub, args)
		}
	case "getPatient":
		{
			return t.getPatient(stub, args)
		}
	case "insure":
		{
			return t.insure(stub, args)
		}
	case "getInsurance":
		{
			return t.getInsurance(stub, args[0])
		}
	case "recordDrug":
		{
			return t.recordDrug(stub, args)
		}
	case "findDrug":
		{
			return t.findDrug(stub, args)
		}
	case "setRegis":
		{
			return t.setRegis(stub,args)
		}
	case "getRegis":
		{
			return t.getRegis(stub,args)
		}
	default:
		{
			return shim.Error("暂无此功能，敬请期待")
		}
	}
}


//记录病人
func (t *SimpleChaincode) setPatient(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("给定的参数个数不是3个" + args[0] + "  " + args[len(args)-1])
	}
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.SetEvent(args[2], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("设置病人成功"))
}

//查询病人
func (t *SimpleChaincode) getPatient(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	bytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if bytes == nil {
		return shim.Error("bytes==nil")
	}
	return shim.Success(bytes)
}

//记录药品
func (t *SimpleChaincode) recordDrug(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.SetEvent(args[2], []byte{}) //记录事件
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("存储药品成功"))
}

//查询药品
func (t *SimpleChaincode) findDrug(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	bytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if bytes == nil {
		return shim.Error("没有获取到相应的数据")
	}
	return shim.Success(bytes)
}

//记录挂号单
func (t *SimpleChaincode) setRegis(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.SetEvent(args[2], []byte{}) //记录事件
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("存储挂号单成功"))
}

//查询挂号单
func (t *SimpleChaincode) getRegis(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	bytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if bytes == nil {
		return shim.Error("没有获取到相应的数据")
	}
	return shim.Success(bytes)
}

//存储保险单
func (t *SimpleChaincode) insure(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("给定的参数个数不是3个" + args[0] + "  " + args[len(args)-1])
	}
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.SetEvent(args[2], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("投保成功"))
}

//查询保险单
func (t *SimpleChaincode) getInsurance(stub shim.ChaincodeStubInterface, id string) peer.Response {
	bytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error(err.Error())
	}
	if bytes == nil {
		return shim.Error("bytes==nil")
	}
	return shim.Success(bytes)
}

func (t *SimpleChaincode) set(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 3 {
		return shim.Error("这是set函数，给定的参数个数不符合要求")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent(args[2], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("存储成功"))

}

func (t *SimpleChaincode) get(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("这是get函数给定的参数个数不符合要求")
	}
	result, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("获取数据发生错误")
	}
	if result == nil {
		return shim.Error("没有获取到相应的数据")
	}
	return shim.Success(result)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("启动SimpleChaincode时发生错误: %s", err)
	}
}
