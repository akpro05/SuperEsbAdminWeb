package updatesysUser

import (
	"SuperEsbAdminWeb/model/db"
	"SuperEsbAdminWeb/services"
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
type UpdatesysUser struct {
	beego.Controller
}

type updateData struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Status   string `json:"input_status"`
	Address  string `json:"address"`
	Rolename string `json:"input_role"`
	Language string `json:"input_language"`
}

func (c *UpdatesysUser) Get() {
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update System User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update System User Page Success")
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

	language := sess.Get("language").(string)
	c.Data["language"] = language

	data1, err := services.GetRole()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("GetRole fetch Failed")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}
	var Dis1 Display
	for i := range data1 {
		var d Field1
		d.Id = data1[i][0]
		d.Name = data1[i][1]
		Dis1.Fields1 = append(Dis1.Fields1, d)
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "username :- ", username)
	row, err := db.Db.Query(`select 
	sysuser.id,
	sysuser.fullname,
	sysuser.mobile,
	sysuser.email,
	sysuser.address,
	sysuser.status,
	sysuser.created_at,
	sysuser.role_id,
	sysuser.language from sysuser
	LEFT JOIN roles ON roles.id = sysuser.role_id where sysuser.id=$1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		//err = errors.New("Unable to get SystemUser data")
		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_UNABLE_TO_GET_SYSTEMUSER_DATA"))
			sendFailureResponse(c, "Unable to get SystemUser data")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_UNABLE_TO_GET_SYSTEMUSER_DATA"))
			sendFailureResponse(c, "Impossible d'obtenir les données de l'utilisateur système")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		//err = errors.New("Unable to get SystemUser data")
		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_UNABLE_TO_GET_SYSTEMUSER_DATA"))
			sendFailureResponse(c, "Unable to get SystemUser data")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_UNABLE_TO_GET_SYSTEMUSER_DATA"))
			sendFailureResponse(c, "Impossible d'obtenir les données de l'utilisateur système")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		//err = errors.New("SystemUser data not found")

		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_SYSTEMUSER_DATA_NOT_FOUND"))
			sendFailureResponse(c, "SystemUser data not found")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_SYSTEMUSER_DATA_NOT_FOUND"))
			sendFailureResponse(c, "Données utilisateur système introuvables")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
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

		"RoleData": Dis1,
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
func (c *UpdatesysUser) Post() {
	//var systemusermsg string

	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminId - ", AdminId)

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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search System User Page Fail")
		} else {
			c.Data["DisplayMessage"] = "System User updated successfully"
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search System User  Page Success")
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

	uid := sess.Get("uid").(string)
	c.Data["uid"] = uid

	language := sess.Get("language").(string)
	c.Data["language"] = language

	var Updatevalues updateData
	if err := json.NewDecoder(c.Ctx.Request.Body).Decode(&Updatevalues); err != nil {
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

	input_full_name := Updatevalues.FullName
	input_mobile := Updatevalues.Mobile
	input_email := Updatevalues.Email
	input_address := Updatevalues.Address
	input_status := Updatevalues.Status
	input_role := Updatevalues.Rolename
	input_language := Updatevalues.Language

	fmt.Println("Full name", input_full_name)
	fmt.Println("Mobile", input_mobile)
	fmt.Println("email", input_email)
	fmt.Println("Status", input_status)
	fmt.Println("Role", input_role)
	fmt.Println("Language", input_language)

	var channelstatus bool

	if input_status == "ACTIVE" {

		channelstatus = true

		err = ResetLoginCount(AdminId)
		if err != nil {
			// err = errors.New("Admin User Login Count Reset Failed")
			// sendFailureResponse(c, "Admin User Login Count Reset Failed")
			if language == "english" {
				err = errors.New(beego.AppConfig.String("EN_SYSTEMUSER_USER_LOGIN_COUNT_RESET_FAILED"))
				sendFailureResponse(c, "Admin User Login Count Reset Failed")
				log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
				return
			} else if language == "french" {
				err = errors.New(beego.AppConfig.String("EN_SYSTEMUSER_USER_LOGIN_COUNT_RESET_FAILED"))
				sendFailureResponse(c, "Échec de la réinitialisation du nombre de connexions de l'utilisateur administrateur")
				log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
				return
			}
			return
		}

	} else {
		channelstatus = false
	}

	if input_full_name == "" || input_email == "" || input_mobile == "" || input_address == "" ||
		input_status == "" {
		// err = errors.New("No fields can be empty")
		// sendFailureResponse(c, "No fields can be empty")

		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_ANY_FIELD_CANNOT_BE_EMPTY"))
			sendFailureResponse(c, "Any Field Cannot be empty")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_ANY_FIELD_CANNOT_BE_EMPTY"))
			sendFailureResponse(c, "N'importe quel champ ne peut pas être vide")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}

	fmt.Println(input_address)

	res, err := db.Db.Exec(`UPDATE public."sysuser" SET fullname=$1,mobile=$2,address=$3,updated_by=$4,status=$5,language=$6,updated_at=now(),role_id=$7 WHERE id = $8`, input_full_name, input_mobile, input_address, uid, channelstatus, input_language, input_role, AdminId)
	if err != nil {
		//err = errors.New("SystemUser updation failed")
		//sendFailureResponse(c, "SystemUser updation failed")
		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_SYSTEMUSER_UPDATION_FAILED"))
			sendFailureResponse(c, "SystemUser updation failed")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_SYSTEMUSER_UPDATION_FAILED"))
			sendFailureResponse(c, "Échec de la mise à jour de l'utilisateur système")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}
	i, err := res.RowsAffected()
	if err != nil || i == 0 {
		// err = errors.New("SystemUser updation failed")
		// sendFailureResponse(c, "SystemUser updation failed")
		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_SYSTEMUSER_UPDATION_FAILED"))
			sendFailureResponse(c, "SystemUser updation failed")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_SYSTEMUSER_UPDATION_FAILED"))
			sendFailureResponse(c, "Échec de la mise à jour de l'utilisateur système")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}

	var pop_msg string

	if language == "french" {

		pop_msg = beego.AppConfig.String("FN_SYSTEMUSER_UPDATESUCCESSFULLY")

	} else {

		pop_msg = beego.AppConfig.String("EN_SYSTEMUSER_UPDATESUCCESSFULLY")

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
	result, err := db.Db.Exec("UPDATE sysuser set pass_count=$1 where id=$2 ", count, uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Count Update Fail")
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Count Update Fail")
		return
	}

	if n != 1 {
		err = errors.New("System User Count Update Fail")
		return
	}

	return
}

func sendFailureResponse(c *UpdatesysUser, message string) {
	c.Ctx.Output.JSON(map[string]interface{}{
		"title":   "FAILURE",
		"message": message,
		"status":  false,
	}, false, false)
}
