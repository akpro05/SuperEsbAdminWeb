package esbLogs

import (
	"SuperEsbAdminWeb/session"

	"runtime/debug"

	"SuperEsbAdminWeb/utils"
	"errors"

	// "SuperEsbAdminWeb/model/db"

	"strings"

	"SuperEsbAdminWeb/services"

	"github.com/astaxie/beego"

	// "proyava.com/database/sql"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type Row struct {
	Id                 string
	Timestamp          string
	Timestamp1         string
	RequestId          string
	Url                string
	Service            string
	In_Request         string
	Out_Response       string
	ProducerAccessCode string
}

type EsbLogs struct {
	beego.Controller
}
type searchData struct {
	Service            string `json:"service"`
	RequestId          string `json:"requestid"`
	Url                string `json:"Url"`
	Produceraccesscode string `json:"produceraccesscode"`
	CustomStartDate    string `json:"customStartDate"`
	CustomEndDate      string `json:"customEndDate"`
}

func (c *EsbLogs) Get() {
	Utype := c.Ctx.Input.Param(":Utype")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Utype", Utype)
	log.Println(beego.AppConfig.String("loglevel"), "Info", "EsbLogs Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	var Autherr error
	sessErr := false
	defer func() {

		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.Data["DisplayMessage"] = "Something went wrong.Please Contact CustomerCare."
			c.TplName = "error/error.html"
		}
		if Autherr != nil {
			c.Data["DisplayMessage"] = Autherr.Error()
			c.TplName = "error/error.html"
			return
		}
		if err != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)

			} else {
				c.Data["DisplayMessage"] = err.Error()
			}
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search EsbLogs Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search EsbLogs  Page Success")
		}
		return
	}()

	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
		sessErr = true
		return
	}

	if err = session.ValidateSession(sess); err != nil {
		sess.SessionRelease(c.Ctx.ResponseWriter)
		session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		sessErr = true
		return
	}
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()

	// auth, err := utils.IsAuthorized(utype, "UserManagment")
	// if !auth {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
	// 	Autherr = errors.New("UnAuthorized")
	// 	return
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	username := sess.Get("username").(string)
	username1 := strings.ToUpper(username)
	c.Data["username"] = username1

	uname := sess.Get("uname").(string)
	c.Data["uname"] = uname

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "username :- ", username)

	// auth, err := utils.IsAuthorized(role, "sysusermanagement-menu", "searchsysuser-active")
	// if !auth {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
	// 	Autherr = errors.New("UnAuthorized")
	// 	return
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "language - ", language)

	data, err := services.GetESBLogs()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System Users fetch Failed")
		return
	}

	var result []Row
	var ts, tdate, ttime string

	for i := range data {
		var r Row
		r.Id = data[i][0]
		r.RequestId = data[i][2]
		r.Url = data[i][3]
		r.Service = data[i][4]
		r.In_Request = data[i][5]
		r.Out_Response = data[i][6]
		r.ProducerAccessCode = data[i][7]
		ts = data[i][1]
		tdate = ts[0:10]
		ttime = ts[11:19]
		r.Timestamp = tdate + " " + ttime
		result = append(result, r)

	}
	c.Data["CustomerData"] = result

	responseData := map[string]interface{}{
		"esb_logs_data": result,
	}

	acceptHeader := c.Ctx.Input.Header("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
		err := json.NewEncoder(c.Ctx.ResponseWriter).Encode(responseData)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.Body([]byte("Error encoding JSON response"))
			return
		}
	} else {
		indexPath := filepath.Join("views", "index.html") // Adjust the file path as needed
		content, err := ioutil.ReadFile(indexPath)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.Body([]byte("Error loading index.html"))
			return
		}
		c.Ctx.Output.Header("Content-Type", "text/html")
		c.Ctx.Output.Body(content)

	}
}
func (c *EsbLogs) Post() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "EsbLogs page")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	var Autherr error
	sessErr := false
	defer func() {
		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.TplName = "error/error.html"
		}
		if Autherr != nil {
			c.Data["DisplayMessage"] = Autherr.Error()
			c.TplName = "error/error.html"
			return
		}
		if err != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)
			} else {
				c.Data["DisplayMessage"] = err.Error()
			}
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", " EsbLogs Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "EsbLogs  Page Success")
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)
	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
		sessErr = true
		return
	}

	if err = session.ValidateSession(sess); err != nil {
		sess.SessionRelease(c.Ctx.ResponseWriter)
		session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		sessErr = true
		return
	}
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()

	var Searchvalues searchData
	if err := json.NewDecoder(c.Ctx.Request.Body).Decode(&Searchvalues); err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		c.Data["DisplayMessage"] = "Invalid Request Received"
		c.Ctx.Output.Status = http.StatusBadRequest // Set the HTTP status to indicate a bad request
		c.Ctx.Output.JSON(map[string]string{
			"Tittle":  "FAILURE",
			"Message": "Invalid Request Received",
			"Type":    "failure",
		}, false, false)
		return
	}

	input_requestId := Searchvalues.RequestId
	input_service := Searchvalues.Service
	input_url := Searchvalues.Url
	input_produceraccesscode := Searchvalues.Produceraccesscode

	from := Searchvalues.CustomStartDate
	to := Searchvalues.CustomEndDate

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "RequestId - ", input_requestId)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Service - ", input_service)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "From Date - ", from)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "to Date - ", to)

	data, err := services.SearchEsbLogsByFilter(input_requestId+"%", input_service+"%", from, to, input_url+"%", input_produceraccesscode+"%")
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("SystemUsers fetch Failed")
		return
	}

	var result []Row
	var ts, tdate, ttime string

	for i := range data {
		var r Row
		r.Id = data[i][0]
		r.RequestId = data[i][2]
		r.Url = data[i][3]
		r.Service = data[i][4]
		r.In_Request = data[i][5]
		r.Out_Response = data[i][6]
		r.ProducerAccessCode = data[i][7]
		ts = data[i][1]
		tdate = ts[0:10]
		ttime = ts[11:19]
		r.Timestamp = tdate + " " + ttime
		result = append(result, r)

	}
	c.Data["CustomerData"] = result
	responseData := map[string]interface{}{
		"esblogdata": result,
	}

	acceptHeader := c.Ctx.Input.Header("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
		err := json.NewEncoder(c.Ctx.ResponseWriter).Encode(responseData)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.Body([]byte("Error encoding JSON response"))
			return
		}
	} else {
		indexPath := filepath.Join("views", "index.html") // Adjust the file path as needed
		content, err := ioutil.ReadFile(indexPath)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Ctx.Output.Body([]byte("Error loading index.html"))
			return
		}
		c.Ctx.Output.Header("Content-Type", "text/html")
		c.Ctx.Output.Body(content)
		return
	}
}
