package viewConsumers

import (
	"SuperEsbAdminWeb/session"

	"SuperEsbAdminWeb/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime/debug"

	"errors"
	"strconv"
	"strings"

	"SuperEsbAdminWeb/model/db"

	"github.com/astaxie/beego"

	"SuperEsbAdminWeb/utils/database/sql"

	log "github.com/sirupsen/logrus"
)

type Row struct {
	Id       string
	FullName string
	Mobile   string
	Email    string
	Address  string
}
type ViewConsumers struct {
	beego.Controller
}

func (c *ViewConsumers) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminId", AdminId)
	log.Println(beego.AppConfig.String("loglevel"), "Info", "View Clients Start")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "View Consumer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "View Consumer  Page Success")
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
	row, err := db.Db.Query(`select id, consumer_name, status, email, consumer_services,consumer_code,consumer_address,created_at from public."consumers" where id=$1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get SystemUser data")
		sendFailureResponse(c, "Unable to get Client data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get SystemUser data")
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
		"ConsumerName":     data[0][1],
		"Status":           status,
		"Email":            data[0][3],
		"consumerservices": data[0][4],
		"ConsumerCode":     data[0][5],
		"consumeraddress":  data[0][6],
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

func sendFailureResponse(c *ViewConsumers, message string) {
	c.Ctx.Output.JSON(map[string]interface{}{
		"title":   "FAILURE",
		"message": message,
		"status":  false,
	}, false, false)
}
