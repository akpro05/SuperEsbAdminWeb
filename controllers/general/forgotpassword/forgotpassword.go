package forgotpassword

import (
	"SuperEsbAdminWeb/model/db"
	"SuperEsbAdminWeb/session"
	"SuperEsbAdminWeb/utils"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"path/filepath"
	"runtime/debug"
	"strconv"

	"SuperEsbAdminWeb/utils/database/sql"

	log "github.com/sirupsen/logrus"

	"SuperEsbAdminWeb/utils/encoding/base64"

	"github.com/astaxie/beego"
	"github.com/scorredoira/email"

	//"proyava.com/database/sql"

	//"proyava.com/util/log"
	"SuperEsbAdminWeb/utils/util/password"
	"SuperEsbAdminWeb/utils/util/pbkdf2"

	"io/ioutil"
	"net/mail"
	"strings"
)

type unencryptedAuth struct {
	smtp.Auth
}
type Forgotpassword struct {
	beego.Controller
}
type ForgotReactRequest struct {
	Email    string `json:"email"`
	Language string `json:"language"`
}

func (c *Forgotpassword) Get() {

	var err error
	sessErr := false

	defer func() {
		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.Data["DisplayMessage"] = "Something went wrong.Please Contact CustomerCare."
			c.TplName = "error/error.html"
		}
		if err != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)

			} else {
				c.Data["DisplayMessage"] = err.Error()
			}
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Forgot Password  Page Fail")
		} else {

			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Forgot Password  Page Success")
		}
		return
	}()
	//	utils.SetHTTPHeader(c.Ctx)

	session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", c.Ctx.Input.IP())
	c.TplName = "general/login/login.html"

	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
		return
	}
	defer sess.SessionRelease(c.Ctx.ResponseWriter)
	sessionId := sess.SessionID()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Session ID - ", sessionId)

	// validationpath1 := beego.AppConfig.String("VALIDATION_LANG_PATH")
	// sess.Set("VALIDATION_LANG_PATH", validationpath1)

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "VALIDATION_LANG_PATH - ", validationpath1)

	// vpath := sess.Get("VALIDATION_LANG_PATH").(string)
	// c.Data["vpath"] = vpath
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "vpath", vpath)

	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()

	return
}

func (c *Forgotpassword) Post() {

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Forgot Password post page")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Forgot Password Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Password has been reset successfully."
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Forgot Password  Page Success")
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
	defer sess.SessionRelease(c.Ctx.ResponseWriter)
	sessionId := sess.SessionID()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Session ID - ", sessionId)

	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()

	//Response from frontend
	var Forgotvalues ForgotReactRequest
	if err := json.NewDecoder(c.Ctx.Request.Body).Decode(&Forgotvalues); err != nil {
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

	//utils.SetHTTPHeader(c.Ctx)

	uname := Forgotvalues.Email
	frontendlanguage := Forgotvalues.Language
	fmt.Println("Language received from the frontend:", frontendlanguage)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "UserEmail - ", uname)

	// if uname == "" {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", "Blank User Email")
	// 	err = errors.New("User-email can't be blank.")
	// 	return
	// }

	err = SearchUser(uname)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		//	err = errors.New("User Not Found")
		return
	}

	newPass, _ := password.AlphaNumericSpecial(6)

	err = UpdatePassword(uname, newPass)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System Admin User Update Password Failed")
		if frontendlanguage == "english" {
			errMessage := beego.AppConfig.String("EN_SYSTEM_ADMIN_USER_UPADTE_PASSWORD_FAILED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "System Admin User Update Password Failed")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if frontendlanguage == "french" {
			errMessage := beego.AppConfig.String("FN_SYSTEM_ADMIN_USER_UPADTE_PASSWORD_FAILED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "Échec de la mise à jour du mot de passe de l'utilisateur administrateur système")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "New Password :- ", newPass)

	var sendmailtoErr error

	sendmailto := SendEmail(uname, "", newPass)
	if sendmailto != nil {
		log.Println("error while sending the email:", sendmailto)
		sendmailtoErr = sendmailto
	}

	if sendmailtoErr != nil {
		log.Println("error while sending the email:", sendmailtoErr)

		var errorMessage string

		if frontendlanguage == "english" {
			errorMessage = "Failed to send reset email" // English error message
		} else if frontendlanguage == "french" {
			errorMessage = "Échec de l'envoi du courrier de réinitialisation" // French error message
		}

		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": errorMessage,
		}
	} else {
		// Successful login, redirect to the dashboard or set the success message
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"message": "Password reset successfully", // English message
		}
	}

	// Check Accept header for response type
	acceptHeader := c.Ctx.Input.Header("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		// Prepare JSON response
		responseData := map[string]interface{}{
			"success": c.Data["json"].(map[string]interface{})["success"],
			"message": c.Data["json"].(map[string]interface{})["message"],
		}

		c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
		json.NewEncoder(c.Ctx.ResponseWriter).Encode(responseData)
	} else {
		// Prepare HTML response
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

func SearchUser(uname string) (err error) {

	row, err := db.Db.Query("SELECT id,email,status FROM public.sysuser where email=$1 limit 1", uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System Admin User Not Found")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System Admin User Detail Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("System Admin User Not Found")
		return
	}

	status, _ := strconv.ParseBool(data[0][2])

	if status == false {
		err = errors.New("System Admin User is currently Suspended")
		return
	}

	return
}

func UpdatePassword(uname, password string) (err error) {

	b := make([]byte, 32)
	_, err = rand.Read(b)
	var pbkdf pbkdf2.Pbkdf2
	pbkdf.Itr = 32
	pbkdf.KeyLen = 32
	pbkdf.Plain = []byte(password)
	pbkdf.Salt = b
	err = pbkdf.Encrypt()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(" System Admin User Password Update Fail")
		return
	}
	var tmp []byte
	tmp = append(tmp, pbkdf.Salt...)
	tmp = append(tmp, pbkdf.Cipher...)

	out, err := base64.Encode(tmp)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(" System Admin User Password Update Fail")
		return
	}

	result, err := db.Db.Exec("update sysuser set password=$1,password_updated_date=now() where email=$2 ", string(out), uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(" System Admin User Password Update Fail")
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(" System Admin User Password Update Fail")
		return
	}

	if n != 1 {
		err = errors.New("System Admin User Password Update Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "System Admin User new Password : ", password)
	return
}

func SendEmail(emilid string, name string, password string) (err error) {
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "called - ")

	uname := beego.AppConfig.String("EMAIL_NOTIFY_USERNAME")
	pass := beego.AppConfig.String("EMAIL_NOTIFY_PASSWORD")
	url := beego.AppConfig.String("EMAIL_NOTIFY_URL")
	to := beego.AppConfig.String("EMAIL_NOTIFY_TIMEOUT")
	recipients := strings.Split(emilid, "||")
	loginurl := beego.AppConfig.String("EMAIL_APPLICATION_LOGIN_URL")

	tmpFile := beego.AppConfig.String("EMAIL_TEMPLATE")

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
	m.From = mail.Address{Name: "Proadmin", Address: uname}
	m.To = recipients

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
func sendFailureResponse(c *Forgotpassword, message string) {
	c.Ctx.Output.JSON(map[string]interface{}{
		"title":   "FAILURE",
		"message": message,
		"status":  false,
	}, false, false)
}
