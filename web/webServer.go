package web

import (
	"fmt"
	"github.com/KangChain.com/KangChain/web/controller"
	"net/http"
)

func WebStart(app *controller.Application) {

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/patient.html", app.PatientView)
	http.HandleFunc("/index.html", app.IndexView)
	http.HandleFunc("/insure.html", app.InsureView)
	http.HandleFunc("/drug.html", app.DrugView)
	http.HandleFunc("/registration.html", app.RegistrationView)

	http.HandleFunc("/recordDrugReq", app.RecordDrug)
	http.HandleFunc("/findDrugReq", app.FindDrug)

	http.HandleFunc("/recordRegisReq", app.RecordRegis)
	http.HandleFunc("/findRegisReq", app.FindRegis)

	http.HandleFunc("/SetPatient", app.SetPatient)
	http.HandleFunc("/GetPatient", app.GetPatient)

	http.HandleFunc("/setInsureReq", app.SetIsure)
	http.HandleFunc("/queryInsureReq", app.QueryInsurance)

	http.HandleFunc("/setReq", app.SetInfo)
	http.HandleFunc("/queryReq", app.QueryInfo)

	fmt.Println("启动Web服务, 监听端口号: 9999")

	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		fmt.Println("启动Web服务错误")
	}

}
