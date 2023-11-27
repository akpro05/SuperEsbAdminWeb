package updateConsumers

import (
	"SuperEsbAdminWeb/model/db"
	//	"SuperEsbAdminWeb/services"
	"SuperEsbAdminWeb/session"
	"SuperEsbAdminWeb/utils"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime/debug"

	"SuperEsbAdminWeb/utils/database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	log "github.com/sirupsen/logrus"
)

type Row struct {
	Id           string
	FirstName    string
	MiddleName   string
	LastName     string
	Mobile       string
	Email        string
	Role         string
	Status       string
	Address1     string
	Address2     string
	Town         string
	City         string
	Pincode      string
	Language     string
	LocationType string
	LocationInfo string
}
type Display struct {
	Fields1 []Field1
	Fields2 []Field2
}
type Field struct {
	Id    string
	Name  string
	Email string
}

type Display1 struct {
	Fields1 []Field1
}

type Field1 struct {
	Id   string
	Name string
}
type Field2 struct {
	Id   string
	Name string
}
type UpdateConsumers struct {
	beego.Controller
}

type updateData struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	Consumercode     string `json:"consumercode"`
	Status           string `json:"input_status"`
	ConsumerServices string `json:"consumerservices"`
	ConsumerAddress  string `json:"consumeraddress"`
}

func (c *UpdateConsumers) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminId", AdminId)
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Update SystemUser Start")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Consumer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Consumer Page Success")
		}
		return
	}()
	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer is unable to process your request.Please contact customer care")
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
	row, err := db.Db.Query(`select id, consumer_name, status, email, consumer_services,consumer_code,consumer_address,created_at from public."consumers" where id=$1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get consumer data")
		sendFailureResponse(c, "Unable to get consumer data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get consumer data")
		sendFailureResponse(c, "Unable to get consumer data")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	var status string

	s1 := data[0][2]
	b1, _ := strconv.ParseBool(s1)

	if b1 == true {

		status = "ACTIVE"

	} else {

		status = "INACTIVE"

	}

	//Response for Forntend

	responseData := map[string]interface{}{
		"Id":               data[0][0],
		"Name":             data[0][1],
		"Status":           status,
		"Email":            data[0][3],
		"ConsumerServices": data[0][4],
		"ConsumerCode":     data[0][5],
		"ConsumerAddress":  data[0][6],
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
			c.Ctx.Output.SetStatus(500)
			c.Ctx.Output.Body([]byte("Error loading index.html"))
			return
		}
		c.Ctx.Output.Header("Content-Type", "text/html")
		c.Ctx.Output.Body(content)
	}
	return
}
func (c *UpdateConsumers) Post() {
	//var systemusermsg string

	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminId - ", AdminId)

	log.Println(beego.AppConfig.String("loglevel"), "Info", "add asset post page")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Update IP - ", pip)
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Consumer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "System User updated successfully"
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Consumer  Page Success")
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

	username := sess.Get("username").(string)
	username1 := strings.ToUpper(username)
	c.Data["username"] = username1

	uname := sess.Get("uname").(string)
	c.Data["uname"] = uname

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "username :- ", username)

	uid := sess.Get("uid").(string)
	c.Data["uid"] = uid

	language := sess.Get("language").(string)
	c.Data["language"] = language

	var createvalues updateData
	if err := json.NewDecoder(c.Ctx.Request.Body).Decode(&createvalues); err != nil {
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

	input_name := createvalues.Name
	input_email := createvalues.Email
	input_consumerservices := createvalues.ConsumerServices
	input_status := createvalues.Status
	input_consumercode := createvalues.Consumercode
	input_consumeraddress := createvalues.ConsumerAddress

	fmt.Println("ConsumerName", input_name)
	fmt.Println("ConsumerServices", input_consumerservices)
	fmt.Println("Email", input_email)
	fmt.Println("Status", input_status)
	fmt.Println("ConsumerCode", input_consumercode)
	fmt.Println("ConsumerAddress", input_consumeraddress)

	var channelstatus bool

	if input_status == "ACTIVE" {

		channelstatus = true

	} else {
		channelstatus = false
	}

	res, err := db.Db.Exec(`UPDATE public."consumers" SET email=$1,consumer_services=$2,status=$3,consumer_code=$4,consumer_address=$5,updated_by=$6,updated_at=now() WHERE id = $7`, input_email, input_consumerservices, channelstatus, input_consumercode, input_consumeraddress, uid, AdminId)
	if err != nil {
		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_CONSUMER_UPDATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("EN_CONSUMER_UPDATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_CONSUMER_UPDATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("FN_CONSUMER_UPDATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}
	i, err := res.RowsAffected()
	if err != nil || i == 0 {
		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_CONSUMER_UPDATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("EN_CONSUMER_UPDATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_CONSUMER_UPDATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("FN_CONSUMER_UPDATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}

	var pop_msg string

	if language == "french" {

		pop_msg = beego.AppConfig.String("FN_CONSUMER_UPDATED_SUCCESSFULLY")

	} else {

		pop_msg = beego.AppConfig.String("EN_CONSUMER_UPDATED_SUCCESSFULLY")

	}

	c.Ctx.Output.JSON(map[string]interface{}{
		"message": pop_msg,
	}, true, false)

	responseMap := map[string]interface{}{
		"success": true, // Indicate success in the response
		"message": "Login successful",
	}

	responseData, _ := json.Marshal(responseMap)

	c.Ctx.Output.Status = http.StatusOK
	c.Ctx.Output.Header("Content-Type", "application/json")
	c.Ctx.Output.Body(responseData)
	return
}

func ResetLoginCount(uname string) (err error) {
	count := 0
	result, err := db.Db.Exec("UPDATE public.'consumer' set pass_count=$1 where id=$2 ", count, uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Consumer Count Update Fail")
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer Count Update Fail")
		return
	}

	if n != 1 {
		err = errors.New("Consumer Count Update Fail")
		return
	}

	return
}

func sendFailureResponse(c *UpdateConsumers, message string) {
	c.Ctx.Output.JSON(map[string]interface{}{
		"title":   "FAILURE",
		"message": message,
		"status":  false,
	}, false, false)
}
