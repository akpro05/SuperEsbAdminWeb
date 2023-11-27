package createProducers

import (
	"SuperEsbAdminWeb/session"

	"encoding/json"
	"net/http"

	"net/smtp"
	"runtime/debug"
	"strings"

	//	"time"

	"SuperEsbAdminWeb/model/db"
	"SuperEsbAdminWeb/utils"
	"errors"

	"github.com/astaxie/beego"

	"SuperEsbAdminWeb/services"

	// "proyava.com/database/sql"
	// "proyava.com/util/log"

	"SuperEsbAdminWeb/utils/database/sql"

	"fmt"

	//	"SuperEsbAdminWeb/utils/util/password"
	//	"SuperEsbAdminWeb/utils/util/txnno"
	"SuperEsbAdminWeb/utils/util/password"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type CreateProducer struct {
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
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"input_status"`
}

func (c *CreateProducer) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "CreateProducer Start")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "CreateProducer  Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "CreateProducer  Page Success")
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
func (c *CreateProducer) Post() {
	//	var systemusermsg string

	log.Println(beego.AppConfig.String("loglevel"), "Info", "CreateProducer page")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "CreateProducer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "System User created Successfully"
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "CreateProducer  Page Success")
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
	input_status := createvalues.Status

	var channelstatus bool

	if input_status == "ACTIVE" {

		channelstatus = true
	} else {
		channelstatus = false
	}

	err = CheckUserAlreadyExists(input_name)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)

		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_PRODUCER_ALREADY_EXISTS"))
			sendFailureResponse(c, beego.AppConfig.String("EN_PRODUCER_ALREADY_EXISTS"))
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_PRODUCER_ALREADY_EXISTS"))
			sendFailureResponse(c, beego.AppConfig.String("FN_PRODUCER_ALREADY_EXISTS"))
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}

		return
	}

	id := uuid.New()
	fmt.Println(id.String())

	//	txn_id := txnno.Generate13Digit()

	accesscode1, _ := password.AlphaNumericSpecial(6)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "access code - ", accesscode1)

	accesscode, err := services.EncryptPassword(accesscode1)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	result, err := db.Db.Exec(`INSERT INTO public."producers" (id,
	    producer_name,
		email,
		status,
		access_code,
		created_by,
		created_at)
		VALUES ($1, $2, $3, $4,$5,$6,now())`,
		id,
		input_name,
		input_email,
		channelstatus,
		accesscode,
		uid)
	if err != nil {
		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_PRODUCER_CREATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("EN_PRODUCER_CREATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_PRODUCER_CREATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("FN_PRODUCER_CREATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {

		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_PRODUCER_CREATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("EN_PRODUCER_CREATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_PRODUCER_CREATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("FN_PRODUCER_CREATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}

		return
	}

	var pop_msg string

	if language == "french" {

		pop_msg = beego.AppConfig.String("FN_PRODUCER_CREATED_SUCCESSFULLY")

	} else {

		pop_msg = beego.AppConfig.String("EN_PRODUCER_CREATED_SUCCESSFULLY")

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

	row, err := db.Db.Query(`SELECT count(*) FROM public."producers" where producer_name = $1`, email)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch Producer")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch Producer")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	countlen := data[0][0]

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

	if countlen != "0" {
		err = errors.New("Producer already exists")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}

func sendFailureResponse(c *CreateProducer, message string) {
	c.Ctx.Output.JSON(map[string]interface{}{
		"title":   "FAILURE",
		"message": message,
		"status":  false,
	}, false, false)
}
