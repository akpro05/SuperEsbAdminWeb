package services

import (
	"strings"

	"errors"
	"io/ioutil"
	"net/mail"
	"net/smtp"

	"SuperEsbAdminWeb/model/db"
	"fmt"

	"SuperEsbAdminWeb/utils/database/sql"

	"github.com/scorredoira/email"

	"SuperEsbAdminWeb/utils/encoding/base64"

	log "github.com/sirupsen/logrus"

	"SuperEsbAdminWeb/utils/util/pbkdf2"

	"crypto/rand"
	//"encoding/json"
	//"strconv"

	"github.com/astaxie/beego"
	//"proyava.com/encoding/base64"
	//"royava.com/util/pbkdf2"
)

type unencryptedAuth struct {
	smtp.Auth
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
	m.From = mail.Address{Name: "OMINAYA", Address: uname}
	m.To = recipients

	// send it
	//auth := smtp.PlainAuth("", uname, pass, url)

	config := beego.AppConfig.String("EMAIL_AUTH_CONFIG_MODE")

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

func EncryptPassword(pass string) (out []byte, err error) {

	//commenting display of password in logs
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "inside encryption password rec: ", pass)
	b := make([]byte, 32)
	_, err = rand.Read(b)
	var pbkdf pbkdf2.Pbkdf2
	pbkdf.Itr = 32
	pbkdf.KeyLen = 32
	pbkdf.Plain = []byte(pass)
	pbkdf.Salt = b
	err = pbkdf.Encrypt()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to encrypt password")
		return
	}
	var tmp []byte
	tmp = append(tmp, pbkdf.Salt...)
	tmp = append(tmp, pbkdf.Cipher...)

	out, err = base64.Encode(tmp)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to encrypt password")
		return
	}
	//	log.Printf("%s %s \n", "inside encrypt pass after encryption:", out)
	return
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

func SearchSystemUsers() (data [][]string, err error) {

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
	LEFT JOIN roles ON roles.id = sysuser.role_id
	order by created_at desc`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get SystemUser info")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the SystemUser info")
		return

	}

	return
}

func SearchProducers() (data [][]string, err error) {

	row, err := db.Db.Query(`select id, producer_name, status, email, producer_services, created_at,producer_code from public."producers"`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer info")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Producer Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the Producer info")
		return

	}

	return
}
func SearchProducersByFilter(producer_name, email, from, to, status string) (data [][]string, err error) {

	fmt.Println(producer_name)
	fmt.Println(email)
	fmt.Println(from)
	fmt.Println(to)
	fmt.Println("channel got at SearchProducersByFilter is ", status)

	if status == "" {

		row, err := db.Db.Query(`select id,
producer_name,
status,
email,
producer_services,
created_at,
producer_code
from public."producer"
where (producer_name='' OR producer_name like $1) AND (email='' OR email like $2) 
	AND ($3='' OR TO_DATE(created_at::text,'YYYY/MM/DD') >= TO_DATE( $3,'MM/DD/YYYY')) AND ($4='' OR TO_DATE(created_at::text,'YYYY/MM/DD') <= TO_DATE( $4,'MM/DD/YYYY'))`, producer_name, email, from, to)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Channel info")
			return data, err
		}
		defer sql.Close(row)
		_, data, err = sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Channel Scan Fail")
			return data, err
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("Unable to get the Channel info")
			return data, err
		}

	} else {

		var channelstatus bool

		if status == "ACTIVE" {

			channelstatus = true
		} else if status == "INACTIVE" {
			channelstatus = false
		} else {

		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "channelstatus", channelstatus)

		row, err := db.Db.Query(`select id, producer_name, status, email, producer_services,  created_at,producer_code from public."producer"  where (producer_name='' OR producer_name like $1) AND (email='' OR email like $2) AND status= $3 AND ($4='' OR TO_DATE(created_at::text,'YYYY/MM/DD') >= TO_DATE( $4,'MM/DD/YYYY')) AND ($5='' OR TO_DATE(created_at::text,'YYYY/MM/DD') <= TO_DATE( $5,'MM/DD/YYYY'))`, producer_name, email, channelstatus, from, to)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Producer info")
			return data, err
		}
		defer sql.Close(row)
		_, data, err = sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Producer Scan Fail")
			return data, err
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("Unable to get the Producer info")
			return data, err
		}
	}
	return

}

func SearchProducer() (data [][]string, err error) {

	row, err := db.Db.Query(`select id, producer_name, email, access_code, created_at, status from public."producers" order by created_at desc`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producers info")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producers Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Producers Query Data - ", data, " Data len - ", len(data))

	return
}

func GetActiveProducer() (data [][]string, err error) {

	row, err := db.Db.Query(`select id, producer_name, email, access_code, created_at, status from public."producers" where status=true order by created_at desc`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get GetActiveProducer info")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("GetActiveProducer Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "GetActiveProducer Query Data - ", data, " Data len - ", len(data))

	return
}

func GetProducerToConsumer() (data [][]string, err error) {

	row, err := db.Db.Query(`SELECT
    producer_to_consumer.id,
    producer_to_consumer.producer_id,
    producer_to_consumer.consumer_id,
    producer_to_consumer.producer_subscribed_services,
    producer_to_consumer.status,
    producer_to_consumer.created_at,
    producers.producer_name,
    consumers.consumer_name,
    consumers.consumer_address
FROM producer_to_consumer
LEFT JOIN producers ON producers.id = producer_to_consumer.producer_id
LEFT JOIN consumers ON consumers.id = producer_to_consumer.consumer_id
 order by producer_to_consumer.created_at desc`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to  GetProducerToConsumer info")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("GetProducerToConsumer Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "GetProducerToConsumer Query Data - ", data, " Data len - ", len(data))

	return
}
func SearchGetProducerToConsumerByFilter(producername, consumername, from, to, status string) (data [][]string, err error) {

	fmt.Println("channel got at SearchChannelByFilter is ", status)

	if status == "" {

		row, err := db.Db.Query(`SELECT
    producer_to_consumer.id,
    producer_to_consumer.producer_id,
    producer_to_consumer.consumer_id,
    producer_to_consumer.producer_subscribed_services,
    producer_to_consumer.status,
    producer_to_consumer.created_at,
    producers.producer_name,
    consumers.consumer_name,
    consumers.consumer_address
FROM producer_to_consumer
LEFT JOIN producers ON producers.id = producer_to_consumer.producer_id
LEFT JOIN consumers ON consumers.id = producer_to_consumer.consumer_id
	where (producers.producer_name='' OR producers.producer_name like $1) AND (consumers.consumer_name='' OR consumers.consumer_name like $2) 
	AND ($3='' OR TO_DATE(producer_to_consumer.created_at::text,'YYYY/MM/DD') >= TO_DATE( $3,'MM/DD/YYYY')) AND ($4='' OR TO_DATE(producer_to_consumer.created_at::text,'YYYY/MM/DD') <= TO_DATE( $4,'MM/DD/YYYY')) order by producer_to_consumer.created_at desc`, producername, consumername, from, to)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get ProducerToConsumer info")
			return data, err
		}
		defer sql.Close(row)
		_, data, err = sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("ProducerToConsumer Scan Fail")
			return data, err
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("Unable to get the ProducerToConsumer info")
			return data, err
		}

	} else {

		var channelstatus bool

		if status == "ACTIVE" {

			channelstatus = true
		} else if status == "INACTIVE" {
			channelstatus = false
		} else {

		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "channelstatus", channelstatus)

		row, err := db.Db.Query(`SELECT
    producer_to_consumer.id,
    producer_to_consumer.producer_id,
    producer_to_consumer.consumer_id,
    producer_to_consumer.producer_subscribed_services,
    producer_to_consumer.status,
    producer_to_consumer.created_at,
    producers.producer_name,
    consumers.consumer_name,
    consumers.consumer_address
FROM producer_to_consumer
LEFT JOIN producers ON producers.id = producer_to_consumer.producer_id
LEFT JOIN consumers ON consumers.id = producer_to_consumer.consumer_id
    where (producers.producer_name='' OR producers.producer_name like $1) AND (consumers.consumer_name='' OR consumers.consumer_name like $2) AND producer_to_consumer.status= $3 AND ($4='' OR TO_DATE(producer_to_consumer.created_at::text,'YYYY/MM/DD') >= TO_DATE( $4,'MM/DD/YYYY')) AND ($5='' OR TO_DATE(producer_to_consumer.created_at::text,'YYYY/MM/DD') <= TO_DATE( $5,'MM/DD/YYYY')) order by producer_to_consumer.created_at desc`, producername, consumername, channelstatus, from, to)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get ProducerToConsumer info")
			return data, err
		}
		defer sql.Close(row)
		_, data, err = sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("ProducerToConsumer Scan Fail")
			return data, err
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("Unable to get the ProducerToConsumer info")
			return data, err
		}
	}
	return

}

func SearchRole() (data [][]string, err error) {
	row, err := db.Db.Query(`select id,role_name,"privilege",created_at ,status from roles  `)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get the error message info")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the error message info")
		return
	}
	return

}

func SearchRoleByFilter(role, from, to, status string) (data [][]string, err error) {

	fmt.Println(role)
	fmt.Println(from)
	fmt.Println(to)
	fmt.Println("channel got at SearchProducersByFilter is ", status)

	if status == "" {

		row, err := db.Db.Query(`select id,role_name,"privilege",created_at ,status from roles
		where (role_name='' OR role_name like $1)
		AND ($2='' OR TO_DATE(created_at::text,'YYYY/MM/DD') >= TO_DATE( $2,'MM/DD/YYYY')) 
		AND ($3='' OR TO_DATE(created_at::text,'YYYY/MM/DD') <= TO_DATE( $3,'MM/DD/YYYY'))`, role, from, to)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Channel info")
			return data, err
		}
		defer sql.Close(row)
		_, data, err = sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Channel Scan Fail")
			return data, err
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("Unable to get the Channel info")
			return data, err
		}

	} else {

		var channelstatus bool

		if status == "ACTIVE" {

			channelstatus = true
		} else if status == "INACTIVE" {
			channelstatus = false
		} else {

		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "channelstatus", channelstatus)

		row, err := db.Db.Query(`select id,role_name,"privilege",created_at ,status from roles
		where (role_name='' OR role_name like $1) AND status= $2 
		AND ($3='' OR TO_DATE(created_at::text,'YYYY/MM/DD') >= TO_DATE( $3,'MM/DD/YYYY')) 
		AND ($4='' OR TO_DATE(created_at::text,'YYYY/MM/DD') <= TO_DATE( $4,'MM/DD/YYYY'))`, role, channelstatus, from, to)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Producer info")
			return data, err
		}
		defer sql.Close(row)
		_, data, err = sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Producer Scan Fail")
			return data, err
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("Unable to get the Producer info")
			return data, err
		}
	}
	return

}

func GetRole() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,role_name from roles`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Status Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Status  Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "Data len - ", len(data))
	return

}

func GetActiveRoles() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,role_name from roles Where status=true`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Status Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Status  Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "Data len - ", len(data))
	return

}

func SearchSystemUsersByFilter(name, email, from, to, status string) (data [][]string, err error) {

	fmt.Println(name)
	fmt.Println(email)
	fmt.Println(from)
	fmt.Println(to)
	fmt.Println("channel got at SearchChannelByFilter is ", status)

	if status == "" {

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
	LEFT JOIN roles ON roles.id = sysuser.role_id
	where (sysuser.fullname='' OR sysuser.fullname like $1) 
	AND (sysuser.email='' OR sysuser.email like $2) 
	AND ($3='' OR TO_DATE(sysuser.created_at::text,'YYYY/MM/DD') >= TO_DATE( $3,'MM/DD/YYYY')) 
	AND ($4='' OR TO_DATE(sysuser.created_at::text,'YYYY/MM/DD') <= TO_DATE( $4,'MM/DD/YYYY')) 
	order by sysuser.created_at desc`, name, email, from, to)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Channel info")
			return data, err
		}
		defer sql.Close(row)
		_, data, err = sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Channel Scan Fail")
			return data, err
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("Unable to get the Channel info")
			return data, err
		}

	} else {

		var channelstatus bool

		if status == "ACTIVE" {

			channelstatus = true
		} else if status == "INACTIVE" {
			channelstatus = false
		} else {

		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "channelstatus", channelstatus)

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
	LEFT JOIN roles ON roles.id = sysuser.role_id
	where (sysuser.fullname='' OR sysuser.fullname like $1) 
	AND (sysuser.email='' OR sysuser.email like $2) AND sysuser.status= $3 
	AND ($4='' OR TO_DATE(sysuser.created_at::text,'YYYY/MM/DD') >= TO_DATE( $4,'MM/DD/YYYY')) 
	AND ($5='' OR TO_DATE(sysuser.created_at::text,'YYYY/MM/DD') <= TO_DATE( $5,'MM/DD/YYYY'))`, name, email, channelstatus, from, to)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get SystemUser info")
			return data, err
		}
		defer sql.Close(row)
		_, data, err = sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Message Scan Fail")
			return data, err
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("Unable to get the SystemUser info")
			return data, err
		}
	}
	return

}
func SearchProducerByFilter(producer_name, email, from, to, status string) (data [][]string, err error) {

	fmt.Println(producer_name)
	fmt.Println(email)
	fmt.Println(from)
	fmt.Println(to)
	fmt.Println("channel got at SearchProducerByFilter is ", status)

	if status == "" {

		row, err := db.Db.Query(`select id,
		producer_name,
		email,
		access_code,
		status,
		created_at 
		from public."producers"
	    where (producer_name='' OR producer_name like $1) AND
		 (email='' OR email like $2) AND ($3='' OR TO_DATE(created_at::text,'YYYY/MM/DD') >= TO_DATE( $3,'MM/DD/YYYY')) AND ($4='' OR TO_DATE(created_at::text,'YYYY/MM/DD') <= TO_DATE( $4,'MM/DD/YYYY')) order by created_at desc`, producer_name, email, from, to)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Producer info")
			return data, err
		}
		defer sql.Close(row)
		_, data, err = sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Producer Scan Fail")
			return data, err
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("Unable to get the Producer info")
			return data, err
		}

	} else {

		var channelstatus bool

		if status == "ACTIVE" {

			channelstatus = true
		} else if status == "INACTIVE" {
			channelstatus = false
		} else {

		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "channelstatus", channelstatus)

		row, err := db.Db.Query(`select id,producer_name,email,access_code,status,created_at from public."producers"  where (producer_name='' OR producer_name like $1) AND (email='' OR email like $2) AND status= $3 AND ($4='' OR TO_DATE(created_at::text,'YYYY/MM/DD') >= TO_DATE( $4,'MM/DD/YYYY')) AND ($5='' OR TO_DATE(created_at::text,'YYYY/MM/DD') <= TO_DATE( $5,'MM/DD/YYYY'))`, producer_name, email, channelstatus, from, to)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Producer info")
			return data, err
		}
		defer sql.Close(row)
		_, data, err = sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Producer Scan Fail")
			return data, err
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("Unable to get the Producer info")
			return data, err
		}
	}
	return

}

func SearchConsumers() (data [][]string, err error) {

	row, err := db.Db.Query(`select id, consumer_name, status, email, consumer_services, created_at,consumer_code from public."consumers" order by created_at desc`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Consumer info")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Consumer Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Consumer Query Data - ", data, " Data len - ", len(data))

	return
}

func GetActiveConsumer() (data [][]string, err error) {

	row, err := db.Db.Query(`select id, consumer_name, status, email, consumer_services, created_at,consumer_code from public."consumers" where status=true order by created_at desc`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get GetActiveConsumer info")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("GetActiveConsumer Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "GetActiveConsumer Query Data - ", data, " Data len - ", len(data))

	return
}

func GetESBLogs() (data [][]string, err error) {

	row, err := db.Db.Query(`SELECT transaction_id, created_at, request_id, url, service, out_request, out_response, producer_access_code
	FROM public.esb_request_metadata order by created_at desc`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ESB_Logs info")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("ESB_Logs Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the ESB_Logs info")
		return

	}

	return
}

func SearchEsbLogsByFilter(requestid, service, from, to, url, produceraccesscode string) (data [][]string, err error) {

	fmt.Println(requestid)
	fmt.Println(service)
	fmt.Println(from)
	fmt.Println(to)
	fmt.Println(url)
	fmt.Println(produceraccesscode)

	row, err := db.Db.Query(`SELECT 
	transaction_id,
	 created_at, request_id, url, service, out_request, out_response, producer_access_code
    FROM public."esb_request_metadata"
    where (request_id='' OR request_id like $1) AND (url='' OR url like $2) AND (producer_access_code='' OR producer_access_code like $3) AND ($4='' OR TO_DATE(created_at::text,'YYYY/MM/DD') >= TO_DATE( $4,'MM/DD/YYYY')) AND ($5='' OR TO_DATE(created_at::text,'YYYY/MM/DD') <= TO_DATE( $5,'MM/DD/YYYY')) AND (service='' OR service like $6)`, requestid, url, produceraccesscode, from, to, service)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ESB_Logs info")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("ESB_Logs Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the ESB_Logs info")
		return

	}

	return
}

func SearchConsumersByFilter(consumer_name, email, from, to, status string) (data [][]string, err error) {

	fmt.Println(consumer_name)
	fmt.Println(email)
	fmt.Println(from)
	fmt.Println(to)
	fmt.Println("channel got at SearchConsumersByFilter is ", status)

	if status == "" {

		row, err := db.Db.Query(`select id,
consumer_name,
status,
email,
consumer_services,
created_at,
consumer_code
from public."consumers"
where (consumer_name='' OR consumer_name like $1) AND (email='' OR email like $2) 
	AND ($3='' OR TO_DATE(created_at::text,'YYYY/MM/DD') >= TO_DATE( $3,'MM/DD/YYYY')) AND ($4='' OR TO_DATE(created_at::text,'YYYY/MM/DD') <= TO_DATE( $4,'MM/DD/YYYY')) order by created_at desc`, consumer_name, email, from, to)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Channel info")
			return data, err
		}
		defer sql.Close(row)
		_, data, err = sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Channel Scan Fail")
			return data, err
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("Unable to get the Channel info")
			return data, err
		}

	} else {

		var channelstatus bool

		if status == "ACTIVE" {

			channelstatus = true
		} else if status == "INACTIVE" {
			channelstatus = false
		} else {

		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "channelstatus", channelstatus)

		row, err := db.Db.Query(`select id, consumer_name, status, email, consumer_services,  created_at,consumer_code from public."consumers"  where (consumer_name='' OR consumer_name like $1) AND (email='' OR email like $2) AND status= $3 AND ($4='' OR TO_DATE(created_at::text,'YYYY/MM/DD') >= TO_DATE( $4,'MM/DD/YYYY')) AND ($5='' OR TO_DATE(created_at::text,'YYYY/MM/DD') <= TO_DATE( $5,'MM/DD/YYYY'))`, consumer_name, email, channelstatus, from, to)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Consumer info")
			return data, err
		}
		defer sql.Close(row)
		_, data, err = sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Consumer Scan Fail")
			return data, err
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("Unable to get the Consumer info")
			return data, err
		}
	}
	return

}

func GetConsumerServices(input_consumer_id string) (data [][]string, err error) {

	row, err := db.Db.Query(`
	WITH service_data AS 
		( SELECT 
		    consumer_address AS consumer_address,
			(json_array_elements(consumer_services::json->'services_list')->>'service_name') AS service_name, 
			(json_array_elements(consumer_services::json->'services_list')->>'service_url') AS service_url FROM consumers where id=$1 ) 
		SELECT service_name, service_url,consumer_address FROM service_data`, input_consumer_id)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ConsumerServices")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("ConsumerServices Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "GetConsumerServices Query Data - ", data, " Data len - ", len(data))

	return
}
func GetConsumerSubscribedServices(input_producertoconsumer_id string) (data [][]string, err error) {
	row, err := db.Db.Query(`
        WITH service_data AS (
            SELECT
                (json_array_elements(producer_subscribed_services::json->'consumer_subsrcibed_services')->>'service_name') AS service_name,
                (json_array_elements(producer_subscribed_services::json->'consumer_subsrcibed_services')->>'service_url') AS service_url
            FROM producer_to_consumer
            WHERE id = $1
        )
        SELECT service_name, service_url FROM service_data;
    `, input_producertoconsumer_id)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ConsumerSubscribedServices")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("ConsumerSubscribedServices Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "GetConsumerSubscribedServices Query Data - ", data, " Data len - ", len(data))

	return
}
