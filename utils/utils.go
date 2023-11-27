package utils

import (
	"encoding/json"
	"errors"
	"net/url"

	"SuperEsbAdminWeb/model/db"

	"SuperEsbAdminWeb/utils/database/sql"
	// "fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/session"

	// "proyava.com/encoding/json/PID"
	// "proyava.com/encoding/json/ReportApi"
	// "proyava.com/encoding/json/ServiceApi"
	//	"proyava.com/net/http"
	// "proyava.com/util/datetime"
	//"proyava.com/util/log"
	// "proyava.com/util/pgp"
	// "proyava.com/util/txnno"
	log "github.com/sirupsen/logrus"
)

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

func IsAuthorized(role, menu1, submenu1 string) (result bool, err error) {

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "role", role)

	result = false

	row, err := db.Db.Query("SELECT role_name,privilege::json->'Menus' as menu,privilege::json->'Submenus' as submenu from roles where role_name=$1 limit 1", role)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get roles ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User Login Count Scan Fail")
		return
	}
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "role count Query Data - ", data, "\nData len - ", len(data))

	role_name := data[0][0]
	menu := data[0][1]
	submenu := data[0][2]
	// fmt.Println("role_name", role_name)
	//fmt.Println(submenu)

	var menusjson, submenusjson string

	// menusjson = ""
	// submenusjson = ""

	if role == role_name {
		menusjson = menu
		submenusjson = submenu
		//	fmt.Println(menusjson)
		//	fmt.Println(submenusjson)
	} else {

	}

	if role == "" {
		err = errors.New("Role invalid")
		return
	}

	var menus []string
	err = json.Unmarshal([]byte(menusjson), &menus)
	if err != nil {
		return
	}

	// for i := range menus {

	// 	fmt.Println("unmarshaled menus :", menus[i])
	// }

	var submenus []string
	err = json.Unmarshal([]byte(submenusjson), &submenus)
	if err != nil {
		return
	}

	for _, men := range menus {
		if strings.EqualFold(men, menu1) {
			result = true
		}
	}

	for _, men := range submenus {
		if strings.EqualFold(men, submenu1) {
			result = true
		}
	}
	return
}

// func EncodeChannel() (ch ServiceApi.Channel, err error) {
// 	ch.Name = beego.AppConfig.String("ChannelName")
// 	ch.Type = beego.AppConfig.String("ChannelType")
// 	var hd ServiceApi.HostIP
// 	hd.Type = beego.AppConfig.String("ChannelIPType")
// 	hd.Value = beego.AppConfig.String("ChannelIP")
// 	return
// }

// func EncodeDevice(pip string) (dv ServiceApi.Device, err error) {
// 	var ip ServiceApi.IP
// 	ip.Type = "IPV4"
// 	ip.Value = pip
// 	return
// }

// func EncodeSecurity(pid []byte) (sc ServiceApi.Security, err error) {
// 	var p pgp.PGP
// 	p.Plain = pid
// 	p.PubKey = beego.AppConfig.String("PGPPublicKey")
// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Public Key Path - ", p.PubKey)
// 	err = p.Encrypt()
// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", "PGP Encryption Fail - ", err)
// 		err = errors.New("Unable to process your request.")
// 		return
// 	}
// 	sc.Data = p.Data
// 	sc.Skey = p.Skey
// 	sc.Hmac = p.Hmac
// 	return
// }

// func Communicate(req []byte, url string, timeout int) (obj ServiceApi.Root, err error) {
// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Url - ", url)
// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "TimeOut - ", timeout)
// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Request - ", string(req))
// 	var net http.HTTP
// 	net.IsHTTPS = false
// 	net.Timeout = timeout
// 	net.Method = "POST"
// 	net.CharacterSet = "text/xml; charset=utf-8"
// 	net.Url = url
// 	rsp, err := net.Communicate(string(req))
// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
// 		err = errors.New("Communication with server fail")
// 		return
// 	}
// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Response - ", string(rsp))
// 	err = json.Unmarshal([]byte(rsp), &obj)
// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", "Service Response Unmarshal Fail- ", err)
// 		err = errors.New("Unable to process your request")
// 		return
// 	}
// 	return
// }

// func CommunicateSwitchReport(req []byte, url string, timeout int) (obj ReportApi.Root, err error) {
// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Url - ", url)
// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "TimeOut - ", timeout)
// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Request - ", string(req))
// 	var net http.HTTP
// 	net.IsHTTPS = false
// 	net.Timeout = timeout
// 	net.Method = "POST"
// 	net.CharacterSet = "text/xml; charset=utf-8"
// 	net.Url = url
// 	rsp, err := net.Communicate(string(req))
// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
// 		err = errors.New("Communication with server fail")
// 		return
// 	}
// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Response - ", string(rsp))
// 	err = json.Unmarshal([]byte(rsp), &obj)
// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", "Service Response Unmarshal Fail- ", err)
// 		err = errors.New("Unable to process your request")
// 		return
// 	}
// 	return
// }

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

// func GetVoucher(sc *ServiceApi.Security) (pid PID.Root, err error) {
// 	var p pgp.PGP
// 	p.Data = sc.Data
// 	p.Skey = sc.Skey
// 	p.Hmac = sc.Hmac
// 	p.PriKey = beego.AppConfig.String("PGPPrivateKey")
// 	err = p.Decrypt()
// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", "PID Decrypt Fail - ", err)
// 		err = errors.New("Unable to process your request")
// 		return
// 	}

// 	err = json.Unmarshal(p.Plain, &pid)
// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", "PID Unmarshal Fail - ", err)
// 		err = errors.New("Unable to process your request")
// 		return
// 	}
// 	return
// }

// func SendSMS(uname, msisdn, text, pip string) (err error) {
// 	msg, err := encodeSendSMSJson(uname, msisdn, text, pip)
// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
// 		return
// 	}

// 	to, err := beego.AppConfig.Int("SendSMSServiceTimeOut")
// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
// 		return
// 	}

// 	obj, err := Communicate(msg, beego.AppConfig.String("SendSMSServiceUrl"), to)
// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
// 		return
// 	}

// 	SI := obj.Tksp.ServiceInterface
// 	res := SI.Response
// 	if res == nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", "Response Node Nil")
// 		return
// 	}
// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Response Code - ", res.Code)
// 	if res.Code != "0" {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", res.Message)
// 		return
// 	}
// 	return
// }

// func SendEmail(uname, emailID, body, subject, pip string) (err error) {
// 	msg, err := encodeSendEmailJson(uname, emailID, body, subject, pip)
// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
// 		return
// 	}

// 	to, err := beego.AppConfig.Int("SendEmailServiceTimeOut")
// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
// 		return
// 	}

// 	obj, err := Communicate(msg, beego.AppConfig.String("SendEmailServiceUrl"), to)
// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
// 		return
// 	}

// 	SI := obj.Tksp.ServiceInterface
// 	res := SI.Response
// 	if res == nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", "Response Node Nil")
// 		return
// 	}
// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Response Code - ", res.Code)
// 	if res.Code != "0" {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", res.Message)
// 		return
// 	}
// 	return
// }

// func encodeSendSMSJson(uname, msisdn, text, pip string) (msg []byte, err error) {
// 	var obj ServiceApi.Root

// 	obj.Tksp.ServiceInterface.Name = "SendSMSViaAdmin"
// 	obj.Tksp.ServiceInterface.Version = "1.0"
// 	obj.Tksp.ServiceInterface.Type = "SendSMS"
// 	ch, err := EncodeChannel()
// 	if err != nil {
// 		return
// 	}
// 	obj.Tksp.ServiceInterface.Channel = &ch

// 	dv, err := EncodeDevice(pip)
// 	if err != nil {
// 		return
// 	}
// 	obj.Tksp.ServiceInterface.Device = &dv

// 	var usrs ServiceApi.SysUsers
// 	usrs.Username = uname
// 	usrs.Type = "SYS_USER"
// 	obj.Tksp.ServiceInterface.SysUsers = append(obj.Tksp.ServiceInterface.SysUsers, &usrs)

// 	var td ServiceApi.TransactionData
// 	var sd ServiceApi.ServiceDetail
// 	sd.BusinessCategory = "B2B"
// 	sd.Partner = "NA"
// 	sd.Category = "NOTIFICATION"
// 	sd.SubCategory = "SMS"
// 	sd.UserType = "SYSTEM"
// 	sd.Denomination = "0"
// 	td.ServiceDetail = &sd
// 	td.CountryCode = "IND"
// 	td.TimeStamp, _ = datetime.Get("", "", "Africa/Harare")
// 	td.RequestID = txnno.Generate()
// 	td.TransactionDescription = "Send SMS for user" + uname
// 	obj.Tksp.ServiceInterface.TransactionData = &td

// 	var not ServiceApi.Notification
// 	not.Dest = msisdn
// 	not.Type = "MOBILE"
// 	not.Message = text
// 	obj.Tksp.ServiceInterface.Notification = &not
// 	msg, err = json.Marshal(&obj)
// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
// 		err = errors.New("Transaction History Request Encoding Fail")
// 		return
// 	}
// 	return
// }

// func encodeSendEmailJson(uname, emailID, body, subject, pip string) (msg []byte, err error) {
// 	var obj ServiceApi.Root

// 	obj.Tksp.ServiceInterface.Name = "SendEmailViaAdmin"
// 	obj.Tksp.ServiceInterface.Version = "1.0"
// 	obj.Tksp.ServiceInterface.Type = "SendEmail"
// 	ch, err := EncodeChannel()
// 	if err != nil {
// 		return
// 	}
// 	obj.Tksp.ServiceInterface.Channel = &ch

// 	dv, err := EncodeDevice(pip)
// 	if err != nil {
// 		return
// 	}
// 	obj.Tksp.ServiceInterface.Device = &dv

// 	var usrs ServiceApi.SysUsers
// 	usrs.Username = uname
// 	usrs.Type = "SYS_USER"
// 	obj.Tksp.ServiceInterface.SysUsers = append(obj.Tksp.ServiceInterface.SysUsers, &usrs)

// 	var td ServiceApi.TransactionData
// 	var sd ServiceApi.ServiceDetail
// 	sd.BusinessCategory = "SYSTEM"
// 	sd.Partner = "NA"
// 	sd.Category = "NOTIFICATION"
// 	sd.SubCategory = "EMAIL"
// 	sd.Denomination = "0"
// 	sd.UserType = "SYSTEM"
// 	td.ServiceDetail = &sd
// 	td.CountryCode = "IND"
// 	td.TimeStamp, _ = datetime.Get("", "", "Africa/Harare")
// 	td.RequestID = txnno.Generate()
// 	td.TransactionDescription = "Send email for user " + uname
// 	obj.Tksp.ServiceInterface.TransactionData = &td

// 	var not ServiceApi.Notification
// 	not.Dest = emailID
// 	not.Type = subject
// 	not.Message = body
// 	obj.Tksp.ServiceInterface.Notification = &not
// 	msg, err = json.Marshal(&obj)
// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
// 		err = errors.New("Transaction History Request Encoding Fail")
// 		return
// 	}
// 	return
// }
