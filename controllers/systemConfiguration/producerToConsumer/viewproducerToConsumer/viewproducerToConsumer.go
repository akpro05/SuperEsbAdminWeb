package viewProducerToConsumer

import (
	"SuperEsbAdminWeb/session"

	// "SuperEsbAdminWeb/services"
	//	"SuperEsbAdminWeb/utils"
	"encoding/json"
	"errors"

	// "fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime/debug"
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
type Field3 struct {
	ServiceName string `json:"service_name"`
	ServiceUrl  string `json:"service_url"`
}

type Display3 struct {
	Field3 []Field3
}

type ViewProducerToConsumer struct {
	beego.Controller
}

func (c *ViewProducerToConsumer) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminId", AdminId)
	log.Println(beego.AppConfig.String("loglevel"), "Info", "ViewProducerToConsumer Start")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "ViewProducerToConsumer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "ViewProducerToConsumer  Page Success")
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
	// defer func() {
	// 	utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
	// 	sess.SessionRelease(c.Ctx.ResponseWriter)
	// }()
	username := sess.Get("username").(string)
	username1 := strings.ToUpper(username)
	c.Data["username"] = username1

	uname := sess.Get("uname").(string)
	c.Data["uname"] = uname

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "username :- ", username)
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
LEFT JOIN consumers ON consumers.id = producer_to_consumer.consumer_id where producer_to_consumer.id=$1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get SubscriberToProducer data")
		sendFailureResponse(c, "Unable to get ProducerToConsumer data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get SystemUser data")
		sendFailureResponse(c, "Unable to get ProducerToConsumer data")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	var status string

	s1 := data[0][4]
	b1, _ := strconv.ParseBool(s1)

	if b1 == true {

		status = "ACTIVE"

	} else {

		status = "INACTIVE"

	}

	// data4, err := services.GetConsumerSubscribedServices(data[0][0])
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("SearchConsumers fetch Failed")
	// 	return
	// }
	// var field3Array []Field3

	// for i := 0; i < len(data4); i++ {
	// 	if len(data4[i]) >= 2 {
	// 		field3 := Field3{
	// 			ServiceName: data4[i][0],
	// 			ServiceUrl:  data4[i][1],
	// 		}
	// 		field3Array = append(field3Array, field3)
	// 	} else {
	// 		log.Printf("Insufficient data for row %d", i)
	// 	}
	// }
	// fmt.Println("data[0][0]", data[0][0])

	//Response for Forntend

	responseData := map[string]interface{}{
		"Id":                    data[0][0],
		"ProducerName":          data[0][6],
		"ConsumerName":          data[0][7],
		"ProducerServices":      data[0][3],
		"Status":                status,
		"ConsumerDomainAddress": data[0][8],
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

func sendFailureResponse(c *ViewProducerToConsumer, message string) {
	c.Ctx.Output.JSON(map[string]interface{}{
		"title":   "FAILURE",
		"message": message,
		"status":  false,
	}, false, false)
}
