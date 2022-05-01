package controller

import (
	"fmt"
	"github.com/KangChain.com/KangChain/service"
	"net/http"
)

type Application struct {
	Fabric *service.ServiceSetup
}

func (app *Application) PatientView(w http.ResponseWriter, r *http.Request) {
	firstdata := service.Patient{}
	firstdata.State = false
	showView(w, r, "patient.html", firstdata)
}

func (app *Application) IndexView(w http.ResponseWriter, r *http.Request) {
	firstdata := service.Data{}
	firstdata.State = false
	showView(w, r, "index.html", firstdata)
}

func (app *Application) InsureView(w http.ResponseWriter, r *http.Request) {
	firstdata := service.InsuranceForm{}
	firstdata.State = false
	showView(w, r, "insure.html", firstdata)
}

func (app *Application) DrugView(w http.ResponseWriter, r *http.Request) {
	firstdata := service.Drug{}
	firstdata.State = false
	showView(w, r, "drug.html", firstdata)
}

func (app *Application) RegistrationView(w http.ResponseWriter, r *http.Request) {
	firstdata := service.RegistrationForm{}
	firstdata.State = false
	showView(w, r, "registration.html", firstdata)
}


// 设置病人
func (app *Application) SetPatient(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	Id := r.FormValue("Id")
	Name := r.FormValue("Name")
	Money := r.FormValue("Money")

	MedicalID := r.FormValue("MedicalID")
	MedicalDate := r.FormValue("MedicalDate")
	MedicalDoctorID := r.FormValue("MedicalDoctorID")
	MedicalDoctorName := r.FormValue("MedicalDoctorName")
	MedicalDoctorDepartment := r.FormValue("MedicalDoctorDepartment")
	MedicalDoctorDepartmentID := r.FormValue("MedicalDoctorDepartmentID")
	MedicalDiagnosisResults := r.FormValue("MedicalDiagnosisResults")

	RegistrationID := r.FormValue("RegistrationID")
	RegistrationDate := r.FormValue("RegistrationDate")
	RegistrationDoctorID := r.FormValue("MedicalDoctorID")
	RegistrationDoctorName := r.FormValue("MedicalDoctorName")
	RegistrationDoctorDepartment := r.FormValue("MedicalDoctorDepartment")
	RegistrationDoctorDepartmentID := r.FormValue("MedicalDoctorDepartmentID")

	InsuranceID := r.FormValue("InsuranceID")
	InsuranceDate := r.FormValue("InsuranceDate")
	InsuranceCompany := r.FormValue("InsuranceCompany")
	InsuranceContent := r.FormValue("InsuranceContent")

	Medicaldoctor := service.Doctor{Id: MedicalDoctorID, Name: MedicalDoctorName, DepartmentName: MedicalDoctorDepartment, DepartmentId: MedicalDoctorDepartmentID}
	Regisdoctor := service.Doctor{Id: RegistrationDoctorID, Name: RegistrationDoctorName, DepartmentName: RegistrationDoctorDepartment, DepartmentId: RegistrationDoctorDepartmentID}

	insform := service.InsuranceForm{Id: InsuranceID, Date: InsuranceDate, Company: InsuranceCompany, InsuranceContent: InsuranceContent}
	regisForm := service.RegistrationForm{Id: RegistrationID, Date: RegistrationDate, Doctor: Regisdoctor}
	medicalForm := service.MedicalRecord{Id: MedicalID, Date: MedicalDate, Doctor: Medicaldoctor, DiagnosisResults: MedicalDiagnosisResults}

	// 调用业务层, 反序列化
	transactionID, err := app.Fabric.SetPatient(Id, Name, Money, insform, regisForm, medicalForm)
	// 封装响应数据
	setdata := service.Patient{}
	if err != nil {
		fmt.Println("报错")
		fmt.Println(err.Error())
		setdata.State = false
		setdata.TransactionID = err.Error()
	} else {
		setdata.State = true
		setdata.TransactionID = "操作成功，交易ID: " + transactionID
	}
	// 响应客户端
	showView(w, r, "patient.html", setdata)
}

// 查询病人
func (app *Application) GetPatient(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	id := r.FormValue("id")
	// 调用业务层, 反序列化
	getdata, err := app.Fabric.GetPatient(id)

	// 封装响应数据
	getdata.State = false
	getdata.TransactionID = ""
	if err != nil {
		getdata.Id = "没有查询到 " + id + " 对应的信息"
	}
	getdata.Id = getdata.Id
	// 响应客户端
	showView(w, r, "patient.html", getdata)
}

// 设置保险
func (app *Application) SetIsure(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	id := r.FormValue("id")
	company := r.FormValue("company")
	insuranceContent := r.FormValue("insuranceContent")
	date := r.FormValue("date")
	// 调用业务层, 反序列化
	transactionID, err := app.Fabric.SetInsure(id, company, insuranceContent, date)
	// 封装响应数据
	setdata := service.InsuranceForm{}
	if err != nil {
		fmt.Println("报错")
		fmt.Println(err.Error())
		setdata.State = false
		setdata.TransactionID = err.Error()
	} else {
		setdata.State = true
		setdata.TransactionID = "操作成功，交易ID: " + transactionID
	}
	// 响应客户端
	showView(w, r, "insure.html", setdata)
}

// 查询保险
func (app *Application) QueryInsurance(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	id := r.FormValue("id")
	var getdata service.InsuranceForm
	// 调用业务层, 反序列化
	getdata, err := app.Fabric.GetInsurance(id)

	// 封装响应数据
	getdata.State = false
	getdata.TransactionID = ""
	if err != nil {
		getdata.Id = "没有查询到 " + id + " 对应的信息"
	}
	// 响应客户端
	showView(w, r, "insure.html", getdata)
}

// 存储药品
func (app *Application) RecordDrug(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	name := r.FormValue("storeName")
	id := r.FormValue("storeId")
	lname := r.FormValue("storeLName")
	ltime := r.FormValue("storeLTime")
	lplace := r.FormValue("storeLPlace")
	lcontent := r.FormValue("storeLContent")
	// 调用业务层, 反序列化
	transactionID, err := app.Fabric.RecordDrug(name,id,lname,ltime,lplace,lcontent)
	// 封装响应数据
	setdata := service.Drug{}
	if err != nil {
		fmt.Println("报错")
		fmt.Println(err.Error())
		setdata.State = false
		setdata.TransactionID = err.Error()
	} else {
		setdata.State = true
		setdata.TransactionID = "操作成功，交易ID: " + transactionID
	}
	// 响应客户端
	showView(w, r, "drug.html", setdata)
}

//查询药品
func (app *Application) FindDrug(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	id := r.FormValue("queryId")
	//fmt.Println("查询参数: "+id)
	// 调用业务层, 反序列化
	getdata := service.Drug{}
	getdata, err := app.Fabric.FindDrug(id)
	// 封装响应数据
	getdata.State=false
	getdata.TransactionID = ""
	if err != nil {
		fmt.Println(err.Error())
		getdata.Id = "没有查询到 "+id+" 对应的信息"
	}
	// 响应客户端
	showView(w, r, "drug.html", getdata)
}

// 存储挂号单
func (app *Application) RecordRegis(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	id := r.FormValue("storeId")
	date := r.FormValue("storeDate")
	doctorid := r.FormValue("storeDoctorId")
	doctorname := r.FormValue("storeDoctorName")
	departmentname := r.FormValue("storeDepartmentName")
	departmentid := r.FormValue("storeDepartmentId")
	// 调用业务层, 反序列化
	transactionID, err := app.Fabric.RecordRegis(id,date,doctorid,doctorname,departmentname,departmentid)
	// 封装响应数据
	setdata := service.RegistrationForm{}
	if err != nil {
		fmt.Println("报错")
		fmt.Println(err.Error())
		setdata.State = false
		setdata.TransactionID = err.Error()
	} else {
		setdata.State = true
		setdata.TransactionID = "操作成功，交易ID: " + transactionID
	}
	// 响应客户端
	showView(w, r, "registration.html", setdata)
}

//查询挂号单
func (app *Application) FindRegis(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	id := r.FormValue("queryId")
	//fmt.Println("查询参数: "+id)
	// 调用业务层, 反序列化
	getdata := service.RegistrationForm{}
	getdata, err := app.Fabric.FindRegis(id)
	// 封装响应数据
	getdata.State=false
	getdata.TransactionID = ""
	if err != nil {
		fmt.Println(err.Error())
		getdata.Id = "没有查询到 "+id+" 对应的信息"
	}
	// 响应客户端
	showView(w, r, "registration.html", getdata)
}

// 根据指定的 key 设置/修改 value 信息
func (app *Application) SetInfo(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	key := r.FormValue("storeKey")
	value := r.FormValue("storeValue")

	// 调用业务层, 反序列化
	transactionID, err := app.Fabric.SetInfo(key, value)

	// 封装响应数据
	setdata := service.Data{QueryData: ""}
	if err != nil {
		setdata.State = false
		setdata.TransactionID = err.Error()
	} else {
		setdata.State = true
		setdata.TransactionID = "操作成功，交易ID: " + transactionID
	}

	// 响应客户端
	showView(w, r, "index.html", setdata)
}

// 根据指定的 Key 查询信息
func (app *Application) QueryInfo(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	name := r.FormValue("queryKey")

	// 调用业务层, 反序列化
	msg, err := app.Fabric.GetInfo(name)

	// 封装响应数据
	getdata := service.Data{}
	getdata.State = false
	getdata.TransactionID = ""
	if err != nil {
		getdata.QueryData = "没有查询到 " + name + " 对应的信息"
	} else {
		getdata.QueryData = msg
	}
	// 响应客户端
	showView(w, r, "index.html", getdata)
}
