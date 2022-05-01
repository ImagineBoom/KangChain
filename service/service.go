package service

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
)

type Common struct {
	State         bool   `json:"State"`
	TransactionID string `json:"TransactionID"`
}
type Patient struct {
	Common
	Id                string             `json:"id"` //身份证号
	Name              string             `json:"name"`
	Money             string             `json:"money"`
	MedicalRecords    []MedicalRecord    `json:"medicalRecords"`    //病历数组
	RegistrationForms []RegistrationForm `json:"registrationForms"` //挂号单数组
	InsuranceForms    []InsuranceForm    `json:"insuranceForms"`    //保险单数组
}
type MedicalRecord struct {
	Common
	Id               string `json:"id"`               //病历id
	Date             string `json:"date"`             //病历写入日期
	Doctor           Doctor `json:"doctor"`           //医生
	DiagnosisResults string `json:"diagnosisResults"` //诊断结果
}
type RegistrationForm struct {
	Common
	Id     string `json:"id"`     //挂号单id
	Date   string `json:"date"`   //挂号单写入日期
	Doctor Doctor `json:"doctor"` //医生
}
type InsuranceForm struct {
	Common
	Id               string `json:"id"`               //保险单id
	Date             string `json:"date"`             //保险单写入日期
	Company          string `json:"company"`          //保险公司名
	InsuranceContent string `json:"insuranceContent"` //保险单详情
}
type Doctor struct {
	Common
	Id             string `json:"id"`             //医生身份证号
	Name           string `json:"name"`           //医生姓名
	DepartmentName string `json:"departmentName"` //科室名称
	DepartmentId   string `json:"departmentId"`   //科室id
}
type Drug struct {
	Common
	Name  string `json:"name"`  //药品名
	Id    string `json:"id"`    //药品id
	Links []Link `json:"links"` //流通环节
}
type Link struct {
	LinkName    string `json:"linkName"`    //经过的环节名称
	LinkTime    string `json:"linkTime"`    //经过环节的时间
	LinkPlace   string `json:"linkPlace"`   //经过环节的地点
	LinkContent string `json:"linkContent"` //在此环节进行的操作
}
type Data struct {
	Common
	QueryData string
}

//存储病人
func (t *ServiceSetup) SetPatient(id, name, money string, insform InsuranceForm, regisform RegistrationForm, medicalRecord MedicalRecord) (string, error) {
	eventID := "eventsetpatient"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	var patient Patient
	patient = Patient{
		Id:                id,
		Name:              name,
		Money:             money,
		InsuranceForms:    append(patient.InsuranceForms, insform),
		RegistrationForms: append(patient.RegistrationForms, regisform),
		MedicalRecords:    append(patient.MedicalRecords, medicalRecord),
	}
	funBytes:=[]byte("setPatient")
	idBytes:=[]byte(id)
	eventIDBytes:=[]byte(eventID)
	patientBytes, err := json.Marshal(patient)
	if err != nil {
		return "", err
	}
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "invoke", Args: [][]byte{funBytes, idBytes, patientBytes, eventIDBytes}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

//查询病人
func (t *ServiceSetup) GetPatient(id string) (Patient, error) {
	funBytes:=[]byte("getPatient")
	idBytes:=[]byte(id)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "invoke", Args: [][]byte{funBytes, idBytes}}
	response, err := t.Client.Query(req)
	if err != nil {
		return Patient{}, fmt.Errorf("查询错误")
	}
	var patient Patient
	err = json.Unmarshal(response.Payload, &patient) //反序列化
	if err != nil {
		return Patient{}, err
	}
	//spew.Dump(patient)
	return patient, nil
}

//存储药品
func (t *ServiceSetup) RecordDrug(name ,id,lname ,ltime ,lplace ,lcontent string) (string, error) {
	var drug Drug
	eventID := "eventrecorddrug"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	link:=Link{
		LinkName: lname,LinkTime: ltime,LinkPlace: lplace,LinkContent: lcontent}
	drug=Drug{
		Name: name,Id: id,Links:append(drug.Links,link),
	}
	//序列化
	drugBytes,err:=json.Marshal(drug)
	if err!=nil{
		return "", err
	}
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "invoke", Args: [][]byte{[]byte("recordDrug"),[]byte(drug.Id),drugBytes,[]byte(eventID)}}
	//respone, err := t.Client.Execute(req,channel.WithRetry(retry.DefaultChannelOpts))
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", fmt.Errorf("存储数据失败")
	}
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

//查询药品
func (t *ServiceSetup) FindDrug(id string) (Drug, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "invoke", Args: [][]byte{[]byte("findDrug"),[]byte(id)}}
	respone, err := t.Client.Query(req,channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return Drug{}, fmt.Errorf("查询错误"+err.Error())
	}
	var drug Drug
	//反序列化
	err=json.Unmarshal(respone.Payload,&drug)
	if err!=nil{
		return Drug{}, err
	}
	return drug, nil
}

//存储挂号单
func (t *ServiceSetup) RecordRegis(id,date,doctorid,doctorname,departmentname,departmentid string) (string, error) {
	eventID := "eventsetregistration"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	doctor:=Doctor{
		Id: doctorid,Name: doctorname,DepartmentName: departmentname,DepartmentId: departmentid,
	}
	registration:=RegistrationForm{
		Id: id,Date: date,Doctor: doctor,
	}
	funBytes:=[]byte("setRegis")
	idBytes:=[]byte(id)
	eventIDBytes:=[]byte(eventID)
	//序列化
	regisBytes,err:=json.Marshal(registration)
	if err!=nil{
		return "", err
	}
	fmt.Printf("%+v",registration)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "invoke", Args: [][]byte{funBytes,idBytes,regisBytes,eventIDBytes}}
	//response, err := t.Client.Execute(req,channel.WithRetry(retry.DefaultChannelOpts))
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", fmt.Errorf("存储数据失败")
	}
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

//查询挂号单
func (t *ServiceSetup) FindRegis(id string) (RegistrationForm, error) {
	funArgs:=[]byte("getRegis")
	idArgs:=[]byte(id)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "invoke", Args: [][]byte{funArgs,idArgs}}
	respone, err := t.Client.Query(req,channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return RegistrationForm{}, fmt.Errorf("查询错误"+err.Error())
	}
	var registration RegistrationForm
	//反序列化
	err=json.Unmarshal(respone.Payload,&registration)
	if err!=nil{
		return RegistrationForm{}, err
	}
	return registration, nil
}

//存储保险单
func (t *ServiceSetup) SetInsure(id, company, insuranceContent, date string) (string, error) {
	eventID := "eventSetInsure"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	inform := InsuranceForm{
		Id: id, Company: company, InsuranceContent: insuranceContent, Date: date}
	insformBytes, err := json.Marshal(inform)
	if err != nil {
		return "", err
	}
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "invoke", Args: [][]byte{[]byte("insure"), []byte(inform.Id), insformBytes, []byte(eventID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

//查询保险单
func (t *ServiceSetup) GetInsurance(id string) (InsuranceForm, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "invoke", Args: [][]byte{[]byte("getInsurance"), []byte(id)}}
	response, err := t.Client.Query(req)
	if err != nil {
		return InsuranceForm{}, err
	}
	var insform InsuranceForm
	err = json.Unmarshal(response.Payload, &insform) //反序列化
	if err != nil {
		return InsuranceForm{}, err
	}
	return insform, nil
}

func (t *ServiceSetup) SetInfo(key, value string) (string, error) {

	eventID := "eventSetInfo"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "invoke", Args: [][]byte{[]byte("set"), []byte(key), []byte(value), []byte(eventID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(response.TransactionID), nil
}
func (t *ServiceSetup) GetInfo(name string) (string, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "invoke", Args: [][]byte{[]byte("get"), []byte(name)}}
	response, err := t.Client.Query(req)
	if err != nil {
		return "", err
	}

	return string(response.Payload), nil
}
