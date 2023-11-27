package createConsumers

import (
	"SuperEsbAdminWeb/session"

	"io/ioutil"

	"encoding/json"
	"net/http"
	"net/mail"

	"net/smtp"
	"runtime/debug"
	"strings"

	//	"time"

	"SuperEsbAdminWeb/model/db"
	"SuperEsbAdminWeb/utils"
	"errors"

	"github.com/scorredoira/email"

	"github.com/astaxie/beego"

	//	"SuperEsbAdminWeb/services"

	// "proyava.com/database/sql"
	// "proyava.com/util/log"

	"SuperEsbAdminWeb/utils/database/sql"

	"fmt"

	//	"SuperEsbAdminWeb/utils/util/password"
	//	"SuperEsbAdminWeb/utils/util/txnno"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type CreateConsumers struct {
	beego.Controller
}

type unencryptedAuth struct {
	smtp.Auth
}

type Field struct {
	Id    string
	Name  string
	Email string
}

type Display struct {
	Fields1 []Field1
	Fields2 []Field2
}
type Field1 struct {
	Id   string
	Name string
}
type Field2 struct {
	Id   string
	Name string
}
type createData struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	Consumercode     string `json:"consumercode"`
	Status           string `json:"input_status"`
	ConsumerServices string `json:"consumerservices"`
	ConsumerAddress  string `json:"consumeraddress"`
}

func (c *CreateConsumers) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Add Assets Start")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "CreateConsumers  Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "CreateConsumers   Page Success")
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

	uname := sess.Get("uname").(string)
	c.Data["uname"] = uname

	return
}
func (c *CreateConsumers) Post() {
	//	var systemusermsg string

	log.Println(beego.AppConfig.String("loglevel"), "Info", "add asset post page")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Create Consumer  Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Consumer created Successfully"
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Create Consumer   Page Success")
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

	uid := sess.Get("uid").(string)
	c.Data["uid"] = uid

	language := sess.Get("language").(string)
	c.Data["language"] = language

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user email :- ", username)

	log.Printf("Request Body: %s", string(c.Ctx.Input.RequestBody))

	var createvalues createData
	if err := json.NewDecoder(c.Ctx.Request.Body).Decode(&createvalues); err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		log.Println(beego.AppConfig.String("loglevel"), "Error: JSON Decoding Failed for Body", string(c.Ctx.Input.RequestBody))
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

	var channelstatus bool

	if input_status == "ACTIVE" {

		channelstatus = true
	} else {
		channelstatus = false
	}

	// data1, err := TemplateFormate()
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	//err = errors.New("templates fetch Failed")
	// 	if language == "english" {
	// 		err = errors.New(beego.AppConfig.String("EN_TEMPLATE_FETCH_FAILED"))
	// 		log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
	// 		return
	// 	} else if language == "french" {
	// 		err = errors.New(beego.AppConfig.String("FN_TEMPLATE_FETCH_FAILED"))
	// 		log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
	// 		return
	// 	}
	// 	return
	// }
	// for i := range data1 {

	// 	c.Data["Id"] = data1[i][0]
	// 	c.Data["Title"] = data1[i][1]
	// 	c.Data["Desc"] = data1[i][2]
	// 	c.Data["Channel"] = data1[i][3]
	// 	c.Data["Url"] = data1[i][4]
	// 	c.Data["Template1"] = data1[i][5]
	// 	c.Data["Template2"] = data1[i][6]
	// 	c.Data["Template3"] = data1[i][7]
	// 	c.Data["DescribeUrl"] = data1[i][8]

	// }

	err = CheckUserAlreadyExists(input_email)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_CONSUMER_ALREADY_EXISTS"))
			sendFailureResponse(c, beego.AppConfig.String("EN_CONSUMER_ALREADY_EXISTS"))
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_CONSUMER_ALREADY_EXISTS"))
			sendFailureResponse(c, beego.AppConfig.String("FN_CONSUMER_ALREADY_EXISTS"))
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}

		return
	}

	id := uuid.New()
	fmt.Println(id.String())

	//	txn_id := txnno.Generate13Digit()

	result, err := db.Db.Exec(`INSERT INTO public."consumers" (id,
	    consumer_name,
		email,
		status,
		consumer_services,
		consumer_code,
		consumer_address,
		created_by,
		created_at)
		VALUES ($1, $2, $3, $4,$5,$6,$7,$8, now())`,
		id,
		input_name,
		input_email,
		channelstatus,
		input_consumerservices,
		input_consumercode,
		input_consumeraddress,

		uid)
	if err != nil {
		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_CONSUMER_CREATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("EN_CONSUMER_CREATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_CONSUMER_CREATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("FN_CONSUMER_CREATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_CONSUMER_CREATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("EN_CONSUMER_CREATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_CONSUMER_CREATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("FN_CONSUMER_CREATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}

	var pop_msg string

	if language == "french" {

		pop_msg = beego.AppConfig.String("FN_CONSUMER_CREATED_SUCCESSFULLY")

	} else {

		pop_msg = beego.AppConfig.String("EN_CONSUMER_CREATED_SUCCESSFULLY")

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

func CheckUserAlreadyExists(email string) (err error) {

	err = nil

	row, err := db.Db.Query(`SELECT count(*) FROM public."consumers" where email = $1`, email)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch Consumers")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch Consumers")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	countlen := data[0][0]

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

	if countlen != "0" {
		err = errors.New("Consumers already exists")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}

func SendEmail(emilid, name, password, emailtemplate string) {
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "called - ")

	uname := beego.AppConfig.String("EMAIL_NOTIFY_USERNAME")
	pass := beego.AppConfig.String("EMAIL_NOTIFY_PASSWORD")
	url := beego.AppConfig.String("EMAIL_NOTIFY_URL")
	to := beego.AppConfig.String("EMAIL_NOTIFY_TIMEOUT")
	loginurl := beego.AppConfig.String("EMAIL_APPLICATION_LOGIN_URL")
	recipients := strings.Split(emilid, "||")

	tmpFile := emailtemplate

	buff, err := ioutil.ReadFile(tmpFile)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "read file -", err)
		return
	}

	msg := string(buff)
	msg = strings.Replace(string(msg), "{{.Name}}", name, -1)
	msg = strings.Replace(string(msg), "{{.Email}}", emilid, -1)
	msg = strings.Replace(string(msg), "{{.Password}}", password, -1)
	msg = strings.Replace(string(msg), "{{.LoginURL}}", loginurl, -1)

	m := email.NewHTMLMessage("Email", msg)
	m.From = mail.Address{Name: "SuperEsbAdminWeb Admin Office", Address: uname}
	m.To = recipients

	// send it
	//auth := smtp.PlainAuth("", uname, pass, url)

	config := beego.AppConfig.String("EMAIL_AUTH_CONFIG_MODE")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "config", config)

	if config == "1" {
		auth := smtp.PlainAuth("", uname, pass, url)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "auth")
		if err = email.Send(url+":"+to, auth, m); err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}

	} else if config == "2" {

		auth := unencryptedAuth{
			smtp.PlainAuth(
				"",
				uname,
				pass,
				url,
			),
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "no tls auth")
		if err = email.Send(url+":"+to, auth, m); err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}
	} else {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "no auth")
		if err = email.Send(url+":"+to, nil, m); err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Email sent successfully")
	return
}

func sendFailureResponse(c *CreateConsumers, message string) {
	c.Ctx.Output.JSON(map[string]interface{}{
		"title":   "FAILURE",
		"message": message,
		"status":  false,
	}, false, false)
}
