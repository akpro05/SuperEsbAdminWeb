/*Copyright (C) 2020-2030 Supernet Technologies India Private Limited. All rights reserved.
All material appearing on this repository, folder, sub folder and file(s) (“Content/Code etc”) is protected by copyright laws and is the property of Supernet Technologies India Pvt Ltd. You may not copy, reproduce, distribute, publish, display, perform, modify, create derivative works, transmit, or in any way exploit any such content, nor may you distribute any part of this content over any network, including a local area network, sell, offer it for sale, or use such content to construct any kind of database, you may not alter or remove any copyright or other notice from copies of contents on this file.
Copying or storing any contents except as provided by Supernet Technologies India Pvt Ltd is expressly prohibited without prior written permission of the Supernet Technologies.
For permission to use the content, please contact legal@supernet-india.com.
*/
package auditReport

import (
	"SuperEsbAdminWeb/session"
	"SuperEsbAdminWeb/utils/database/sql"

	"runtime/debug"

	"SuperEsbAdminWeb/utils"
	"errors"

	//	"strconv"

	"SuperEsbAdminWeb/model/db"
	//	"SuperEsbAdminWeb/services"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/astaxie/beego"

	log "github.com/sirupsen/logrus"
)

type Row struct {
	Timestamp   string
	Timestamp1  string
	Adminid     string
	URL         string
	Status      string
	IP          string
	Host        string
	HTTPSMethod string
}

type searchData struct {
	Email           string `json:"email"`
	Method          string `json:"method"`
	CustomStartDate string `json:"customStartDate"`
	CustomEndDate   string `json:"customEndDate"`
}

// type Display5 struct {
// 	Fields5 []Field5
// }

// type Field5 struct {
// 	Id     string
// 	Mobile string
// }
type AuditReport struct {
	beego.Controller
}

func (c *AuditReport) Get() {

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Audit Report Page Start")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Audit Report Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Audit Report Page Success")
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)
	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
		return
	}

	if err = session.ValidateSession(sess); err != nil {
		sess.SessionRelease(c.Ctx.ResponseWriter)
		session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		sessErr = true

		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()

	username := sess.Get("username").(string)
	username1 := strings.ToUpper(username)
	c.Data["username"] = username1

	uname := sess.Get("uname").(string)
	c.Data["uname"] = uname

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "username :- ", username)

	return
}

func (c *AuditReport) Post() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder Page Start")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Audit Report Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Audit Report Page Success")
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)
	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
		return
	}

	if err = session.ValidateSession(sess); err != nil {
		sess.SessionRelease(c.Ctx.ResponseWriter)
		session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		sessErr = true

		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)

	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()
	// utype := sess.Get("user_type").(string)
	// c.Data["utype"] = utype
	// //log.Println(beego.AppConfig.String("loglevel"), "Debug", "utype", utype)

	// username := sess.Get("username").(string)
	// username1 := strings.ToUpper(username)
	// c.Data["username"] = username1

	// user_type := sess.Get("user_type").(string)
	// user_type1 := strings.ToUpper(user_type)
	// c.Data["user_type"] = user_type1

	// mobile := sess.Get("mobile").(string)
	// mobile1 := strings.ToUpper(mobile)
	// c.Data["mobile"] = mobile1

	// language := sess.Get("language").(string)
	// c.Data["language"] = language
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	// role := sess.Get("role").(string)
	// c.Data["role"] = role
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "role :- ", role)
	// menu := sess.Get("menu").(string)
	// c.Data["menu"] = menu
	// submenu := sess.Get("submenu").(string)
	// c.Data["submenu"] = submenu
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "usertype :- ", user_type)

	// auth, err := utils.IsAuthorized(role, "reports-menu", "auditreport-active")
	// if !auth {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
	// 	Autherr = errors.New("UnAuthorized")
	// 	return
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

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

	input_email := Searchvalues.Email
	input_method := Searchvalues.Method

	from := Searchvalues.CustomStartDate
	to := Searchvalues.CustomEndDate

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "From Date - ", from)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "to Date - ", to)

	row, err := db.Db.Query(`select created_on,
	event::json->>'UserName' as uname,
	event::json->>'URL' as ul,
	event::json->>'Status' as stat,
	event::json->>'PIP' as ipp,
	event::json->>'Host' as hst,
	event::json->>'Method' as mthd
	from web_event where ((event::json->> 'Method' = '') OR (event::json->> 'Method' like $1)) AND ($2='' OR TO_DATE(created_on::text,'YYYY/MM/DD') >= TO_DATE( $2,'MM/DD/YYYY')) AND ($3='' OR TO_DATE(created_on::text,'YYYY/MM/DD') <= TO_DATE( $3,'MM/DD/YYYY')) AND ((event::json->> 'UserName' ISNULL) OR (event::json->> 'UserName' like $4)) ORDER BY created_on desc LIMIT 1000`, input_method+"%", from, to, input_email+"%")
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch System Audit Data")
		// if language == "english" {
		// 	err = errors.New(beego.AppConfig.String("EN_UNABLE_TO_FETCH_SYSTEM_AUDIT_DATA"))
		// 	log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
		// 	return
		// } else if language == "french" {
		// 	err = errors.New(beego.AppConfig.String("FN_UNABLE_TO_FETCH_SYSTEM_AUDIT_DATA"))
		// 	log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
		// 	return
		// }
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch System Audit Data")
		// if language == "english" {
		// 	err = errors.New(beego.AppConfig.String("EN_UNABLE_TO_FETCH_SYSTEM_AUDIT_DATA"))
		// 	log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
		// 	return
		// } else if language == "french" {
		// 	err = errors.New(beego.AppConfig.String("FN_UNABLE_TO_FETCH_SYSTEM_AUDIT_DATA"))
		// 	log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
		// 	return
		// }
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	var result []Row
	var ts, tdate, ttime string

	for i := range data {
		var r Row
		r.Timestamp1 = data[i][0]

		ts = data[i][0]
		tdate = ts[0:10]
		ttime = ts[11:19]
		r.Timestamp = tdate + " " + ttime

		r.Adminid = data[i][1]
		r.URL = data[i][2]
		r.Status = data[i][3]
		r.IP = data[i][4]
		r.Host = data[i][5]
		r.HTTPSMethod = data[i][6]
		result = append(result, r)
	}
	c.Data["CustomerData"] = result
	responseData := map[string]interface{}{
		"AuditData": result,
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
