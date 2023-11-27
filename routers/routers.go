package routers

import (
	"SuperEsbAdminWeb/controllers/general/changePassword"
	"SuperEsbAdminWeb/controllers/general/dashboard"
	"SuperEsbAdminWeb/controllers/general/forgotpassword"
	"SuperEsbAdminWeb/controllers/general/login"

	"SuperEsbAdminWeb/controllers/general/logout"

	"SuperEsbAdminWeb/controllers/error"
	"SuperEsbAdminWeb/model/db"
	"SuperEsbAdminWeb/session"

	"SuperEsbAdminWeb/controllers/sysusers/createsysUser"
	"SuperEsbAdminWeb/controllers/sysusers/searchsysUser"
	"SuperEsbAdminWeb/controllers/sysusers/updatesysUser"
	"SuperEsbAdminWeb/controllers/sysusers/viewsysUser"

	"SuperEsbAdminWeb/controllers/consumers/createConsumers"
	"SuperEsbAdminWeb/controllers/consumers/searchConsumers"
	"SuperEsbAdminWeb/controllers/consumers/updateConsumers"
	"SuperEsbAdminWeb/controllers/consumers/viewConsumers"

	"SuperEsbAdminWeb/controllers/producers/createProducers"
	"SuperEsbAdminWeb/controllers/producers/searchProducers"
	"SuperEsbAdminWeb/controllers/producers/updateProducers"
	"SuperEsbAdminWeb/controllers/producers/viewProducers"

	"SuperEsbAdminWeb/controllers/systemConfiguration/producerToConsumer/createproducerToConsumer"
	"SuperEsbAdminWeb/controllers/systemConfiguration/producerToConsumer/searchproducerToConsumer"
	"SuperEsbAdminWeb/controllers/systemConfiguration/producerToConsumer/updateproducerToConsumer"
	"SuperEsbAdminWeb/controllers/systemConfiguration/producerToConsumer/viewproducerToConsumer"

	"SuperEsbAdminWeb/controllers/systemConfiguration/role/createRole"
	"SuperEsbAdminWeb/controllers/systemConfiguration/role/searchRole"
	"SuperEsbAdminWeb/controllers/systemConfiguration/role/updateRole"
	"SuperEsbAdminWeb/controllers/systemConfiguration/role/viewRole"

	"SuperEsbAdminWeb/controllers/reports/auditReport"
	"SuperEsbAdminWeb/controllers/reports/consumerReport"
	"SuperEsbAdminWeb/controllers/reports/esbLogs"
	"SuperEsbAdminWeb/controllers/reports/producerReport"

	"github.com/astaxie/beego"
)

func init() {
	if session.Init() != nil {
		return
	}
	if db.Init() != nil {
		return
	}

	beego.SetStaticPath("/frontend/dist", "frontend/dist")

	beego.ErrorController(&error.Error{})

	beego.Router(beego.AppConfig.String("MAIN_PATH"), &login.Login{})
	beego.Router(beego.AppConfig.String("LOGIN_PATH"), &login.Login{})
	beego.Router(beego.AppConfig.String("LOGOUT_PATH"), &logout.Logout{})

	beego.Router(beego.AppConfig.String("DASHBOARD_PATH"), &dashboard.Dashboard{})
	beego.Router(beego.AppConfig.String("CHANGE_PASSWORD_PATH"), &changePassword.ChangePassword{})
	beego.Router(beego.AppConfig.String("FORGOT_PASSWORD_PATH"), &forgotpassword.Forgotpassword{})

	beego.Router(beego.AppConfig.String("CREATE_SYS_USER_PATH"), &createsysUser.CreatesysUser{})
	beego.Router(beego.AppConfig.String("SEACRH_SYS_USER_PATH"), &searchsysUser.SearchsysUser{})
	beego.Router(beego.AppConfig.String("UPDATE_SYS_USER_PATH"), &updatesysUser.UpdatesysUser{})
	beego.Router(beego.AppConfig.String("VIEW_SYS_USER_PATH"), &viewsysUser.ViewsysUser{})

	beego.Router(beego.AppConfig.String("CREATE_PRODUCERS_PATH"), &createProducers.CreateProducer{})
	beego.Router(beego.AppConfig.String("SEACRH_PRODUCERS_PATH"), &searchProducers.SearchProducer{})
	beego.Router(beego.AppConfig.String("UPDATE_PRODUCERS_PATH"), &updateProducers.UpdateProducer{})
	beego.Router(beego.AppConfig.String("VIEW_PRODUCERS_PATH"), &viewProducers.ViewProducer{})

	beego.Router(beego.AppConfig.String("CREATE_CONSUMERS_PATH"), &createConsumers.CreateConsumers{})
	beego.Router(beego.AppConfig.String("SEACRH_CONSUMERS_PATH"), &searchConsumers.SearchConsumers{})
	beego.Router(beego.AppConfig.String("UPDATE_CONSUMERS_PATH"), &updateConsumers.UpdateConsumers{})
	beego.Router(beego.AppConfig.String("VIEW_CONSUMERS_PATH"), &viewConsumers.ViewConsumers{})

	beego.Router(beego.AppConfig.String("CREATE_PROD2CONS_PATH"), &createProducerToConsumer.CreateProducerToConsumer{})
	beego.Router(beego.AppConfig.String("SEACRH_PROD2CONS_PATH"), &searchProducerToConsumer.SearchProducerToConsumer{})
	beego.Router(beego.AppConfig.String("UPDATE_PROD2CONS_PATH"), &updateProducerToConsumer.UpdateProducerToConsumer{})
	beego.Router(beego.AppConfig.String("VIEW_PROD2CONS_PATH"), &viewProducerToConsumer.ViewProducerToConsumer{})
	beego.Router(beego.AppConfig.String("GET_CONSUMERLISTS_PATH"), &createProducerToConsumer.GetConsumersServicelist{})

	beego.Router(beego.AppConfig.String("CREATE_ROLE_PATH"), &createRole.CreateRole{})
	beego.Router(beego.AppConfig.String("SEACRH_ROLE_PATH"), &searchRole.SearchRole{})
	beego.Router(beego.AppConfig.String("UPDATE_ROLE_PATH"), &updateRole.UpdateRole{})
	beego.Router(beego.AppConfig.String("VIEW_ROLE_PATH"), &viewRole.ViewRole{})

	beego.Router(beego.AppConfig.String("ESB_LOGS_REPORT"), &esbLogs.EsbLogs{})
	beego.Router(beego.AppConfig.String("PRODUCER_REPORT"), &producerReport.ProducerReport{})
	beego.Router(beego.AppConfig.String("CONSUMER_REPORT"), &consumerReport.ConsumerReport{})
	beego.Router(beego.AppConfig.String("AUDIT_REPORT"), &auditReport.AuditReport{})

}
