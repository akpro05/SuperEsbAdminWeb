package changePassword

import (
	"SuperEsbAdminWeb/session"
	"SuperEsbAdminWeb/utils"
	"errors"
	"runtime/debug"

	"SuperEsbAdminWeb/services"

	"SuperEsbAdminWeb/model/db"
	//	"strings"
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego"

	"SuperEsbAdminWeb/utils/database/sql"

	log "github.com/sirupsen/logrus"

	"SuperEsbAdminWeb/utils/encoding/base64"
	"SuperEsbAdminWeb/utils/util/pbkdf2"
)

type ChangePassword struct {
	beego.Controller
}

type createData struct {
	Oldpassword     string `json:"oldpassword"`
	Newpassword     string `json:"newpassword"`
	Confirmpassword string `json:"confirmpassword"`
}

func (c *ChangePassword) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Change Password Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	defer func() {
		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)

			c.TplName = "error/error.html"
		}

		if err != nil && err.Error() == "Session Time Out.Please Logout and Login Again." {
			c.Abort("500")
		}

		if err != nil {
			c.Data["DisplayMessage"] = err.Error()
			c.TplName = "error/error.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Change Password Page Fail")
		} else {
			c.TplName = "index.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Change Password Page Success")
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
		return
	}

	uname := sess.Get("uname").(string)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "UserEmail - ", uname)

	// defer func() {
	// 	utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
	// 	sess.SessionRelease(c.Ctx.ResponseWriter)
	// }()
	// return
}

func (c *ChangePassword) Post() {
	//	var changepasswordmsg string
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Change Password Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	defer func() {

		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)

			c.TplName = "error/error.html"
		}
		if err != nil && err.Error() == "Session Time Out.Please Logout and Login Again." {
			c.Abort("500")
		}

		if err != nil {
			c.Data["DisplayMessage"] = err.Error()
			c.Data["title"] = "Error !"
			c.Data["type"] = "error"
			c.TplName = "index.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Change Password Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Change Password Successfully"
			c.Data["title"] = "Success !"
			c.Data["type"] = "success"
			c.TplName = "index.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Change Password Page Success")
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
		return
	}
	// defer func() {
	// 	utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
	// 	sess.SessionRelease(c.Ctx.ResponseWriter)
	// }()

	uname := sess.Get("uname").(string)
	c.Data["uname"] = uname

	uid := sess.Get("uid").(string)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "uid - ", uid)

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

	input_oldpassword := createvalues.Oldpassword
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Old Login Password - ", input_oldpassword)

	input_newpassword := createvalues.Newpassword
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "New Login Password - ", input_newpassword)

	input_confirmnewpassword := createvalues.Confirmpassword
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Old Login Password - ", input_confirmnewpassword)

	if input_newpassword != input_confirmnewpassword {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "New Password Mismatch")
		err = errors.New("New password and Confirm password can't be different")
		sendFailureResponse(c, "New password and Confirm password can't be different")
		return
	}

	err = CheckOldPassword(input_oldpassword, uid)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User  old password Incorrect")
		sendFailureResponse(c, "User  old password Incorrect")
		return
	}

	err = UpdatePassword(uid, input_newpassword)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Change Password Failed")
		sendFailureResponse(c, "Change Password Failed")
		return
	}

	c.Ctx.Output.JSON(map[string]interface{}{
		"message": "Change password created successfully",
	}, true, false)

	responseMap := map[string]interface{}{
		"success": true, // Indicate success in the response
		"message": "Login successful",
	}

	responseData, _ := json.Marshal(responseMap)

	c.Ctx.Output.Status = http.StatusOK
	c.Ctx.Output.Header("Content-Type", "application/json")
	c.Ctx.Output.Body(responseData)

}

func CheckOldPassword(oldpassword, userid string) (err error) {

	err = nil

	row, err := db.Db.Query(`select password from sysuser WHERE id = $1`, userid)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch user")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch user")
		return
	}

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	cp, err := base64.Decode([]byte(data[0][0]))
	if err != nil {
		err = errors.New("Unable to authenticate user")
		return
	}

	var pbkdf pbkdf2.Pbkdf2
	pbkdf.Itr = 32
	pbkdf.KeyLen = 32
	pbkdf.Plain = []byte(oldpassword)
	pbkdf.Salt = cp[:32]
	pbkdf.Cipher = cp[32:]
	result, err := pbkdf.Compare()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User old password incorrect")
		return
	}
	if !result {
		err = errors.New("User old password incorrect")
		return
	}

	return

}

func UpdatePassword(userid, password string) (err error) {

	pass, err := services.EncryptPassword(password)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	result, err := db.Db.Exec("update sysuser set password=$1 ,password_set=true,password_updated_date=now(),updated_at=now() where id=$2 ", pass, userid)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Password Update Fail")
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Password Update Fail")
		return
	}

	if n != 1 {
		err = errors.New("System User Password Update Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pass", password)
	return
}

func sendFailureResponse(c *ChangePassword, message string) {
	c.Ctx.Output.JSON(map[string]interface{}{
		"title":   "FAILURE",
		"message": message,
		"status":  false,
	}, false, false)
}
