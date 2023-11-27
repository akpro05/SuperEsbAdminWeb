package login

import (
	"SuperEsbAdminWeb/model/db"
	"SuperEsbAdminWeb/utils"

	"errors"

	"runtime/debug"
	"strconv"

	"encoding/json"
	"fmt"

	"net/http"
	"time"

	"SuperEsbAdminWeb/session"

	"SuperEsbAdminWeb/utils/database/sql"

	log "github.com/sirupsen/logrus"

	//"proyava.com/database/sql"
	"SuperEsbAdminWeb/utils/encoding/base64"
	//"proyava.com/util/log"
	"SuperEsbAdminWeb/utils/util/pbkdf2"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type Login struct {
	beego.Controller
}

type LoginData struct {
	Email string `form:"email" valid:"Required"`
	Pass  string `form:"password" valid:"Required;MinSize(6);MaxSize(16)"`
}

type LoginReactRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Language string `json:"language"`
}

type roles struct {

	// defining struct variables
	Menu    string          `json:"menu,omitempty"`
	Submenu string          `json:"submenu,omitempty"`
	Data    json.RawMessage `json:"data"`
}

func (c *Login) Get() {
	defer func() {
		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.Data["DisplayMessage"] = "Something went wrong.Please Contact CustomerCare."
			c.TplName = "error/error.html"
		}
		return
	}()
	//	utils.SetHTTPHeader(c.Ctx)

	session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", c.Ctx.Input.IP())
	c.TplName = "index.html"

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

	// defer func() {
	// 	utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
	// 	sess.SessionRelease(c.Ctx.ResponseWriter)
	// }()
	return
}

func (c *Login) Post() {

	msg := ""

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Login Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)

	var err error

	defer func() {

		if err != nil {
			if l_exception := recover(); l_exception != nil {
				stack := debug.Stack()
				log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
				session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
				c.TplName = "error/error.html"
			} else {
				c.Data["DisplayMessage"] = err.Error()
			}
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Login Fail")

		} else {

			c.Data["DisplayMessage"] = msg
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Login Success")
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

	// defer func() {
	// 	utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
	// 	sess.SessionRelease(c.Ctx.ResponseWriter)
	// }()

	// var l LoginData
	// if err := c.ParseForm(&l); err != nil {
	// 	err = errors.New("Invalid Request Received")
	// 	return
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Form Data - ", l)
	// c.Data["FormData"] = l

	var loginvalues LoginReactRequest
	if err := json.NewDecoder(c.Ctx.Request.Body).Decode(&loginvalues); err != nil {
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

	email := loginvalues.Email
	password := loginvalues.Password
	frontendlanguage := loginvalues.Language
	fmt.Println("Language received from the frontend:", frontendlanguage)

	fmt.Println("email", email)
	fmt.Println("password", password)

	valid := validation.Validation{}
	b, err := valid.Valid(&loginvalues)
	if err != nil {
		err = errors.New("Parameter validation failed")
		if frontendlanguage == "english" {
			errMessage := beego.AppConfig.String("EN_PARAMETER_VALIDATION_FAILED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "Parameter validation failed")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if frontendlanguage == "french" {
			errMessage := beego.AppConfig.String("FN_PARAMETER_VALIDATION_FAILED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "La validation des paramètres a échoué")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}

	if !b {
		for _, err := range valid.Errors {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", err.Key, ":", err.Message, ":", err.Field, ":", err.LimitValue, ":", err.Name, ":", err.Tmpl, ":", err.Value)
		}
		err = errors.New("Invalid Input values")
		if frontendlanguage == "english" {
			errMessage := beego.AppConfig.String("EN_INVALID_INPUT_VALUES")
			err = errors.New(errMessage)
			sendFailureResponse(c, "Invalid Input values")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if frontendlanguage == "french" {
			errMessage := beego.AppConfig.String("FN_INVALID_INPUT_VALUES")
			err = errors.New(errMessage)
			sendFailureResponse(c, "Valeurs d'entrée invalides")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}

	err = session.CheckUserSession(email)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	}

	get_status_details, err := db.Db.Query(`SELECT SYSUSER.STATUS, ROLES.STATUS FROM SYSUSER LEFT JOIN ROLES ON ROLES.ID = SYSUSER.ROLE_ID WHERE SYSUSER.EMAIL = $1`, loginvalues.Email)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Details fetch Failed")
		if frontendlanguage == "english" {
			errMessage := beego.AppConfig.String("EN_SYSTEM_USERS_NOT_DETAILS_SCAN_FAILED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "System User Details scan Failed")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if frontendlanguage == "french" {
			errMessage := beego.AppConfig.String("FN_SYSTEM_USERS_NOT_DETAILS_SCAN_FAILED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "Échec de l'analyse des détails de l'utilisateur système")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}
	defer sql.Close(get_status_details)
	_, get_status_details_data, err := sql.Scan(get_status_details)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Details scan Failed")
		if frontendlanguage == "english" {
			errMessage := beego.AppConfig.String("EN_SYSTEM_USERS_NOT_DETAILS_SCAN_FAILED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "System User Details scan Failed")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if frontendlanguage == "french" {
			errMessage := beego.AppConfig.String("FN_SYSTEM_USERS_NOT_DETAILS_SCAN_FAILED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "Échec de l'analyse des détails de l'utilisateur système")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}

	fmt.Println("Count of user with this email :", len(get_status_details_data))

	if len(get_status_details_data) == 0 {

		err = errors.New("System User Not Found - No account linked with this email ,Please verify your details")
		if frontendlanguage == "english" {
			errMessage := beego.AppConfig.String("EN_SYSTEM_USERS_NOT_FOUND_PLEASE_VERIFY")
			err = errors.New(errMessage)
			sendFailureResponse(c, "System User Not Found - No account linked with this email ,Please verify your details")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if frontendlanguage == "french" {
			errMessage := beego.AppConfig.String("FN_SYSTEM_USERS_NOT_FOUND_PLEASE_VERIFY")
			err = errors.New(errMessage)
			sendFailureResponse(c, "Utilisateur système introuvable - Aucun compte lié à cet e-mail, veuillez vérifier vos coordonnées")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return

	}

	sys_user_status, err := strconv.ParseBool(get_status_details_data[0][0])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	sys_user_role_status, err := strconv.ParseBool(get_status_details_data[0][1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("System User Status : ", sys_user_status)
	fmt.Println("System User Role Status : ", sys_user_role_status)

	if sys_user_status == false {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User's account has been deactivated. Please contact support for assistance")
		if frontendlanguage == "english" {
			errMessage := beego.AppConfig.String("EN_SYSTEM_USERS_ACCOUNT_DEACTIVATED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "System User's account has been deactivated. Please contact support for assistance")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if frontendlanguage == "french" {
			errMessage := beego.AppConfig.String("FN_SYSTEM_USERS_ACCOUNT_DEACTIVATED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "Le compte de l'utilisateur système a été désactivé. Veuillez contacter le support pour obtenir de l'aide")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}

		return
	}

	if sys_user_role_status == false {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User's associated Role  has been deactivated. Please contact support for assistance")
		if frontendlanguage == "english" {
			errMessage := beego.AppConfig.String("EN_SYSTEM_USERS_BEEN_DEACTIVATED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "System User's associated Role has been deactivated. Please contact support for assistance")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if frontendlanguage == "french" {
			errMessage := beego.AppConfig.String("FN_SYSTEM_USERS_BEEN_DEACTIVATED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "Le rôle associé à l'utilisateur système a été désactivé. Veuillez contacter le support pour obtenir de l'aide")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}

	count, err := getpasscount(loginvalues.Email)
	if err != nil {

		log.Println(beego.AppConfig.String("loglevel"), "getpasscount Error", err)
		err = errors.New("System User Details fetch Failed")
		if frontendlanguage == "english" {
			errMessage := beego.AppConfig.String("EN_SYSTEM_USERS_DETAILS_FETCH_FAILED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "System User Details fetch Failed")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if frontendlanguage == "french" {
			errMessage := beego.AppConfig.String("FN_SYSTEM_USERS_DETAILS_FETCH_FAILED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "Échec de la récupération des détails de l'utilisateur système")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}

	loginCount, _ := beego.AppConfig.Int("LOGIN_COUNT")

	row, err := db.Db.Query(`SELECT password_set,password_updated_date,language FROM sysuser where  email=$1`, loginvalues.Email)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "password_updated_date Error", err)
		err = errors.New("System User Details fetch Failed")
		if frontendlanguage == "english" {
			errMessage := beego.AppConfig.String("EN_SYSTEM_USERS_DETAILS_FETCH_FAILED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "System User Details fetch Failed")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if frontendlanguage == "french" {
			errMessage := beego.AppConfig.String("FN_SYSTEM_USERS_DETAILS_FETCH_FAILED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "Échec de la récupération des détails de l'utilisateur système")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "password_updated_date Error", err)
		err = errors.New("System User Details fetch Failed")
		if frontendlanguage == "english" {
			errMessage := beego.AppConfig.String("EN_SYSTEM_USERS_DETAILS_FETCH_FAILED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "System User Details fetch Failed")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if frontendlanguage == "french" {
			errMessage := beego.AppConfig.String("FN_SYSTEM_USERS_DETAILS_FETCH_FAILED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "Échec de la récupération des détails de l'utilisateur système")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Password Set Value - ", data[0][0])

	//lang := data[0][2]

	//Password need to have expiry time of 30 calendar days, post expiry force Admin to change

	var ts, year, month, day string

	//Createdat1 := data8[0][1]
	//fmt.Println(Createdat1)
	//for refer 2023-09-29 12:32:46.725791+05:30
	ts = data[0][1]
	year = ts[0:4]
	month = ts[5:7]
	day = ts[8:10]
	Createdat := year + "-" + month + "-" + day
	fmt.Println("Createdat:=", Createdat)

	y1, _ := strconv.Atoi(year)
	m1, _ := strconv.Atoi(month)
	d1, _ := strconv.Atoi(day)

	t := time.Now()
	y2 := t.Year()  // type int
	m2 := t.Month() // type time.Month
	d2 := t.Day()

	t1 := Date(y1, m1, d1)
	t2 := Date(y2, int(m2), d2)
	days := t2.Sub(t1).Hours() / 24
	fmt.Println(days)

	if days >= 30 {

		fmt.Println("Password expiry time of 30 calendar days")
		if frontendlanguage == "english" {

			fmt.Println("Your Password is expired after 30 days  , Please Reset it ")

			msg = "Your Password is expired after 30 days  , Please Reset it  "
			sendFailureResponse(c, "Your Password is expired after 30 days  , Please Reset it")
		} else {

			fmt.Println("Votre mot de passe a expiré après 30 jours, veuillez le réinitialiser")

			msg = "Votre mot de passe a expiré après 30 jours, veuillez le réinitialiser"

		}
		return

	}

	if count >= loginCount {

		result, err := db.Db.Exec(`UPDATE sysuser SET status=$1,updated_at=now() WHERE email=$2`, false, loginvalues.Email)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("System User update details failed")
			if frontendlanguage == "english" {
				errMessage := beego.AppConfig.String("EN_SYSTEM_USERS_UPDATE_DETAILS_FAILED")
				err = errors.New(errMessage)
				sendFailureResponse(c, "System User update details failed")
				log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
				return
			} else if frontendlanguage == "french" {
				errMessage := beego.AppConfig.String("EN_SYSTEM_USERS_UPDATE_DETAILS_FAILED")
				err = errors.New(errMessage)
				sendFailureResponse(c, "Échec de la mise à jour des détails de l'utilisateur système")
				log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
				return
			}
			return
		}

		i, err := result.RowsAffected()
		if err != nil || i == 0 {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("System User update details failed")
			if frontendlanguage == "english" {
				errMessage := beego.AppConfig.String("EN_SYSTEM_USERS_UPDATE_DETAILS_FAILED")
				err = errors.New(errMessage)
				sendFailureResponse(c, "System User update details failed")
				log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
				return
			} else if frontendlanguage == "french" {
				errMessage := beego.AppConfig.String("EN_SYSTEM_USERS_UPDATE_DETAILS_FAILED")
				err = errors.New(errMessage)
				sendFailureResponse(c, "Échec de la mise à jour des détails de l'utilisateur système")
				log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
				return
			}
			return
		}

		msg = "User Auntentication Fail exceded the Limit 3 times , System User is Suspended"
		if frontendlanguage == "english" {
			errMessage := beego.AppConfig.String("EN_USER_AUNTENTICATIONA_FAIL_EXCEDED_LIMIT")
			err = errors.New(errMessage)
			sendFailureResponse(c, "User Auntentication Fail exceded the Limit 3 times , System User is Suspended")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if frontendlanguage == "french" {
			errMessage := beego.AppConfig.String("FN_USER_AUNTENTICATIONA_FAIL_EXCEDED_LIMIT")
			err = errors.New(errMessage)
			sendFailureResponse(c, "L'échec de l'authentification de l'utilisateur a dépassé la limite 3 fois, l'utilisateur du système est suspendu")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}

		return

	}

	name, id, language, role, menu, submenu, err := authinticate(loginvalues.Email, loginvalues.Password)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User Authentication Failed")
		if frontendlanguage == "english" {
			errMessage := beego.AppConfig.String("EN_USER_AUNTENTICATIONA_FAILED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "User Authentication Failed")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if frontendlanguage == "french" {
			errMessage := beego.AppConfig.String("FN_USER_AUNTENTICATIONA_FAILED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "L'identification de l'utilisateur a échoué")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}

		// increment the count
		count++
		// update count
		r_err := PasswordMismatch(count, loginvalues.Email)
		if r_err != nil {

			return
		}

		return
	}

	if count > 0 {
		log.Println(beego.AppConfig.String("loglevel"), "Count Reset Error", err)
		err = ResetLoginCount(loginvalues.Email)
		if err != nil {
			err = errors.New("User Authentication Failed")
			if frontendlanguage == "english" {
				errMessage := beego.AppConfig.String("EN_USER_AUNTENTICATIONA_FAILED")
				err = errors.New(errMessage)
				sendFailureResponse(c, "User Authentication Failed")
				log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
				return
			} else if frontendlanguage == "french" {
				errMessage := beego.AppConfig.String("FN_USER_AUNTENTICATIONA_FAILED")
				err = errors.New(errMessage)
				sendFailureResponse(c, "L'identification de l'utilisateur a échoué")
				log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
				return
			}
			return
		}
	}

	if data[0][1] == "" {
		fmt.Println("Password Update date value got null")
		err = errors.New("User Authentication Failed")
		if frontendlanguage == "english" {
			errMessage := beego.AppConfig.String("EN_USER_AUNTENTICATIONA_FAILED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "User Authentication Failed")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if frontendlanguage == "french" {
			errMessage := beego.AppConfig.String("FN_USER_AUNTENTICATIONA_FAILED")
			err = errors.New(errMessage)
			sendFailureResponse(c, "L'identification de l'utilisateur a échoué")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return

	}
	sess.Set("uname", email)
	sess.Set("username", name)
	sess.Set("uid", id)
	// sess.Set("user_type", string(user_type))
	// sess.Set("mobile", mobile)
	sess.Set("language", language)
	sess.Set("role", role)
	sess.Set("menu", menu)
	sess.Set("submenu", submenu)

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "User id :- ", id)
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "User Name :- ", name)

	err = session.SetUserSession(sess.SessionID(), email)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
		sendFailureResponse(c, "System is unable to process your request.Please contact customer care")
		if frontendlanguage == "english" {
			errMessage := beego.AppConfig.String("EN_SYSTEM_UNABLE_TO_PROCESS_REQUEST")
			err = errors.New(errMessage)
			sendFailureResponse(c, "System is unable to process your request.Please contact customer care")
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if frontendlanguage == "french" {
			errMessage := beego.AppConfig.String("FN_SYSTEM_UNABLE_TO_PROCESS_REQUEST")
			err = errors.New(errMessage)
			sendFailureResponse(c, "Le système n'est pas en mesure de traiter votre demande. Veuillez contacter le service client.")
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}

	c.Ctx.Output.JSON(map[string]interface{}{
		"title":    "SUCCESS",
		"message":  "Login successful",
		"status":   true,
		"menu":     menu,
		"submenu":  submenu,
		"language": language,
	}, false, false)

	c.Redirect("/Dashboard", 302)
	return
}

func getpasscount(uname string) (count int, err error) {
	row, err := db.Db.Query("select pass_count from public.\"sysuser\" where email=$1 and status=true limit 1", uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User Login Count Not Found")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User Login Count Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("User Login Count Not Found")
		return
	}

	if data[0][0] == "" {
		data[0][0] = "0"
	}

	count_str := data[0][0]

	count, _ = strconv.Atoi(count_str)
	return

}

func authinticate(email, pass string) (name, id, language, role, menu, submenu string, err error) {
	row, err := db.Db.Query("select id,email,password,fullname,mobile,language,role_id from public.\"sysuser\" where email=$1", email)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to authenticate user")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to authenticate user")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", " User Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("User not registered")
		return
	}

	cp, err := base64.Decode([]byte(data[0][2]))
	if err != nil {
		err = errors.New("Unable to authenticate user")
		return
	}

	var pbkdf pbkdf2.Pbkdf2
	pbkdf.Itr = 32
	pbkdf.KeyLen = 32
	pbkdf.Plain = []byte(pass)
	pbkdf.Salt = cp[:32]
	pbkdf.Cipher = cp[32:]
	result, err := pbkdf.Compare()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User password incorrect")
		return
	}
	if !result {
		err = errors.New("User password incorrect")
		return
	}

	id = data[0][0]
	email = data[0][1]
	//	password = data[0][2]
	name = data[0][3]
	//	mobile = data[0][4]
	language = data[0][5]

	row1, err := db.Db.Query("select id,role_name,privilege::json->'Menus' as menu,privilege::json->'Submenus' as submenu from roles where id=$1", data[0][6])
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to authenticate user")
		return
	}
	defer sql.Close(row1)
	_, data1, err := sql.Scan(row1)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to authenticate user")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "Data len - ", len(data1))

	if len(data1) <= 0 {
		err = errors.New("User not registered")
		return
	}

	role = data1[0][1]
	menu = data1[0][2]
	submenu = data1[0][3]

	return
}

func PasswordMismatch(count int, uname string) (err error) {
	result, err := db.Db.Exec("update sysuser set pass_count=$1 where email=$2 ", count, uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(" User Count Update Fail")
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(" User Count Update Fail")
		return
	}

	if n != 1 {
		err = errors.New(" User Count Update Fail")
		return
	}

	return
}

func ResetLoginCount(uname string) (err error) {
	count := 0
	result, err := db.Db.Exec("UPDATE sysuser set pass_count=$1 where email=$2 ", count, uname)
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

func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func sendFailureResponse(c *Login, message string) {
	c.Ctx.Output.JSON(map[string]interface{}{
		"title":   "FAILURE",
		"message": message,
		"status":  false,
	}, false, false)
}
