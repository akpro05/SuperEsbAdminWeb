package searchProducers

import (
	"SuperEsbAdminWeb/session"

	"runtime/debug"

	"SuperEsbAdminWeb/utils"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"

	// "SuperEsbAdminWeb/model/db"
	"SuperEsbAdminWeb/services"
	// "fmt"
	"strconv"
	"strings"

	//	"SuperEsbAdminWeb/services"

	"github.com/astaxie/beego"

	// "proyava.com/database/sql"

	log "github.com/sirupsen/logrus"
)

type Row struct {
	Id           string
	PartialId    string
	ProducerName string
	Email        string
	Timestamp    string
	Status       string
	AccessCode   string
}

type SearchProducer struct {
	beego.Controller
}
type searchData struct {
	Producername    string `json:"producerName"`
	Email           string `json:"email"`
	Status          string `json:"input_status"`
	CustomStartDate string `json:"customStartDate"`
	CustomEndDate   string `json:"customEndDate"`
}

func (c *SearchProducer) Get() {
	Utype := c.Ctx.Input.Param(":Utype")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Utype", Utype)
	log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchSubscriber Start")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchSubscriber Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchSubscriber  Page Success")
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

	username := sess.Get("username").(string)
	username1 := strings.ToUpper(username)
	c.Data["username"] = username1

	uname := sess.Get("uname").(string)
	c.Data["uname"] = uname

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "username :- ", username)
	data, err := services.SearchProducer()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producers fetch Failed")
		return
	}

	var result []Row
	var ts, tdate, ttime string

	for i := range data {
		var r Row
		r.Id = data[i][0]
		r.PartialId = r.Id[0:8]
		r.ProducerName = data[i][1]
		r.Email = data[i][2]
		r.AccessCode = data[i][3]

		ts = data[i][4]
		tdate = ts[0:10]
		ttime = ts[11:19]
		r.Timestamp = tdate + " " + ttime

		s1 := data[i][5]
		b1, _ := strconv.ParseBool(s1)

		if b1 == true {

			r.Status = "ACTIVE"

		} else {

			r.Status = "INACTIVE"

		}

		result = append(result, r)

	}
	c.Data["CustomerData"] = result

	//Response for Forntend

	responseData := map[string]interface{}{
		"CustomerData": c.Data["CustomerData"],
	}

	acceptHeader := c.Ctx.Input.Header("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
		json.NewEncoder(c.Ctx.ResponseWriter).Encode(responseData)
	} else {
		indexPath := filepath.Join("views", "index.html") // Adjust the file path as needed

		content, err := ioutil.ReadFile(indexPath)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			c.Ctx.Output.SetStatus(500)
			c.Ctx.Output.Body([]byte("Error loading index.html"))
			return
		}

		c.Ctx.Output.Header("Content-Type", "text/html")
		c.Ctx.Output.Body(content)
	}

	return
}

func (c *SearchProducer) Post() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchProducer page")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchProducer Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchProducer  Page Success")
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

	username := sess.Get("username").(string)
	username1 := strings.ToUpper(username)
	c.Data["username"] = username1

	uname := sess.Get("uname").(string)
	c.Data["uname"] = uname

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "username :- ", username)

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

	input_producername := Searchvalues.Producername
	input_email := Searchvalues.Email
	input_status := Searchvalues.Status

	from := Searchvalues.CustomStartDate
	to := Searchvalues.CustomEndDate

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "From Date - ", from)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "to Date - ", to)

	data, err := services.SearchProducerByFilter(input_producername+"%", input_email+"%", from, to, input_status)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producers fetch Failed")
		return
	}

	var result []Row
	var ts, tdate, ttime string

	for i := range data {
		var r Row
		r.Id = data[i][0]
		r.PartialId = r.Id[0:8]
		r.ProducerName = data[i][1]
		r.Email = data[i][2]
		r.AccessCode = data[i][3]

		ts = data[i][5]
		tdate = ts[0:10]
		ttime = ts[11:19]
		r.Timestamp = tdate + " " + ttime

		s1 := data[i][4]
		b1, _ := strconv.ParseBool(s1)

		if b1 == true {

			r.Status = "ACTIVE"

		} else {

			r.Status = "INACTIVE"

		}

		result = append(result, r)

	}
	c.Data["CustomerData"] = result

	//Response for Forntend

	responseData := map[string]interface{}{
		"CustomerDataPostMethod": c.Data["CustomerData"],
	}

	acceptHeader := c.Ctx.Input.Header("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
		json.NewEncoder(c.Ctx.ResponseWriter).Encode(responseData)
	} else {
		indexPath := filepath.Join("views", "index.html") // Adjust the file path as needed

		content, err := ioutil.ReadFile(indexPath)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			c.Ctx.Output.SetStatus(500)
			c.Ctx.Output.Body([]byte("Error loading index.html"))
			return
		}

		c.Ctx.Output.Header("Content-Type", "text/html")
		c.Ctx.Output.Body(content)
	}

	return
}
