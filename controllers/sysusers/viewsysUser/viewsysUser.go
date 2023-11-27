package viewsysUser

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
type ViewsysUser struct {
	beego.Controller
}

func (c *ViewsysUser) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminId", AdminId)
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Search SystemUser Start")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search SystemUser Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search SystemUser  Page Success")
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
	row, err := db.Db.Query(`select 
	sysuser.id,
	sysuser.fullname,
	sysuser.mobile,
	sysuser.email,
	sysuser.address,
	sysuser.status,
	sysuser.created_at,
	roles.role_name,
	sysuser.language from sysuser
	LEFT JOIN roles ON roles.id = sysuser.role_id where sysuser.id=$1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get SystemUser data")
		sendFailureResponse(c, "Unable to get SystemUser data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get SystemUser data")
		sendFailureResponse(c, "Unable to get SystemUser data")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	for i := range data {

		c.Data["Id"] = data[i][0]
		c.Data["FullName"] = data[i][1]
		c.Data["Mobile"] = data[i][2]
		c.Data["Email"] = data[i][3]
		c.Data["Address"] = data[i][4]
		s1 := data[i][5]
		b1, _ := strconv.ParseBool(s1)

		if b1 == true {

			c.Data["Status"] = "ACTIVE"

		} else {

			c.Data["Status"] = "INACTIVE"

		}
		c.Data["Rolename"] = data[i][7]
		c.Data["Language"] = data[i][8]

	}

	//Response for Forntend

	responseData := map[string]interface{}{
		"Id":       c.Data["Id"],
		"FullName": c.Data["FullName"],
		"Mobile":   c.Data["Mobile"],
		"Email":    c.Data["Email"],
		"Address":  c.Data["Address"],
		"Status":   c.Data["Status"],
		"Rolename": c.Data["Rolename"],
		"Language": c.Data["Language"],
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

func sendFailureResponse(c *ViewsysUser, message string) {
	c.Ctx.Output.JSON(map[string]interface{}{
		"title":   "FAILURE",
		"message": message,
		"status":  false,
	}, false, false)
}
