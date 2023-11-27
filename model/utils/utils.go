package utils

import (
	"SuperEsbAdminWeb/model/db"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/smtp"
	"net/url"
	"strings"
	"text/template"

	"github.com/astaxie/beego/context"

	log "github.com/sirupsen/logrus"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"

	"SuperEsbAdminWeb/utils/encoding/base64"
	p "SuperEsbAdminWeb/utils/util/password"
	"SuperEsbAdminWeb/utils/util/pbkdf2"
	//"/*proyava.com/encoding/base64"
	//"proyava.com/util/log"
	//p "proyava.com/util/password"
	//"proyava.com/util/pbkdf2"/
)

type MenusStruct struct {
	Menus []string
}

func GeneratePassword() (login_pass, encrypted_pwd string, err error) {
	login_pass, _ = p.Numeric(6)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Login Password", login_pass)

	b := make([]byte, 32)
	_, err = rand.Read(b)
	var pbkdf pbkdf2.Pbkdf2
	pbkdf.Itr = 32
	pbkdf.KeyLen = 32
	pbkdf.Plain = []byte(login_pass)
	pbkdf.Salt = b
	err = pbkdf.Encrypt()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to create password")
		return
	}
	var tmp []byte
	tmp = append(tmp, pbkdf.Salt...)
	tmp = append(tmp, pbkdf.Cipher...)

	out, err := base64.Encode(tmp)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to create password")
		return
	}
	encrypted_pwd = string(out)
	return
}

func IsAuthorized(role, page string) (result bool, err error) {

	result = false

	var menusjson string

	menusjson = ""

	if role == "admin@gmail.com" {
		menusjson = beego.AppConfig.String("ADMIN_MENU_ARRAY")
	} else if role == "admin1@gmail.com" {
		menusjson = beego.AppConfig.String("ADMIN_MENU_ARRAY")
	} else if role == "REPORT" {
		menusjson = beego.AppConfig.String("REPORT_MENU_ARRAY")
	}

	if role == "" {
		err = errors.New("Role invalid")
		return
	}

	var menus MenusStruct
	err = json.Unmarshal([]byte(menusjson), &menus)
	if err != nil {
		return
	}

	for _, men := range menus.Menus {
		if strings.EqualFold(men, page) {
			result = true
		}
	}
	return
}

func SendEmail(host string, port int, userName string, password string, to []string, subject string, message string) (err error) {
	parameters := struct {
		From    string
		To      string
		Subject string
		Message string
	}{
		userName,
		strings.Join([]string(to), ","),
		subject,
		message,
	}

	buffer := new(bytes.Buffer)

	template := template.Must(template.New("emailTemplate").Parse(emailScript()))
	template.Execute(buffer, &parameters)

	auth := smtp.PlainAuth("", userName, password, host)

	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		userName,
		to,
		buffer.Bytes())

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	}
	return err
}
func emailScript() (script string) {
	return `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
MIME-version: 1.0
Content-Type: text/html; charset="UTF-8"

{{.Message}}`
}

func SetHTTPHeader(Ctx *context.Context) {
	Ctx.Output.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	Ctx.Output.Header("Pragma", "no-cache")
	Ctx.Output.Header("Expires", "0")
	Ctx.Output.Header("X-Content-Type-Options", "nosniff")
	Ctx.Output.Header("Strict-Transport-Security", "max-age=31536000 ; includeSubDomains")
	Ctx.Output.Header("X-Frame-Options", "SAMEORIGIN")
	Ctx.Output.Header("X-XSS-Protection", "1; mode=block")
	Ctx.Output.Header("X-Content-Security-Policy", "default-src 'self'")
	//Ctx.Output.Header("X-WebKit-CSP", "default-src 'self'")
}
func EventLogs(c *context.Context, sess session.Store, method string, input url.Values, output map[interface{}]interface{}, err_r error) (res string, err error) {
	m2 := make(map[string]interface{})

	for key, value := range output {
		switch key := key.(type) {
		case string:
			m2[key] = value
		}
	}
	event := make(map[string]interface{})

	event["PIP"] = c.Input.IP()
	event["URL"] = c.Input.URL()
	event["SessionID"] = sess.SessionID()
	if sess.Get("uname") != nil {
		event["UserName"] = sess.Get("uname")
	}
	event["Host"] = c.Input.Host()
	event["Method"] = method
	if err_r != nil {
		event["Status"] = err_r

	} else {
		event["Status"] = "Success"
	}

	event["Input"] = input
	event["Output"] = m2

	jsonString, err := json.Marshal(event)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}
	event_json := string(jsonString)

	_, err = db.Db.Exec("INSERT INTO web_event(event, created_on) VALUES ( $1,now())", event_json)

	return
}
