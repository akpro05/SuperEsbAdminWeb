package updateProducerToConsumer

import (
	"SuperEsbAdminWeb/model/db"
	"SuperEsbAdminWeb/services"
	"SuperEsbAdminWeb/session"
	"SuperEsbAdminWeb/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime/debug"

	"SuperEsbAdminWeb/utils/database/sql"
	// "fmt"
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
type Field3 struct {
	ServiceName string `json:"service_name"`
	ServiceUrl  string `json:"service_url"`
}
type Field4 struct {
	ServiceName string `json:"service_name"`
	ServiceUrl  string `json:"service_url"`
}
type Display4 struct {
	Field4 []Field4
}
type Display3 struct {
	Field3 []Field3
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
type UpdateProducerToConsumer struct {
	beego.Controller
}

type updateData struct {
	ProducerServices   string `json:"producerservices"`
	Status             string `json:"input_status"`
	SubscribedServices []struct {
		ServiceName string `json:"service_name"`
		ServiceURL  string `json:"service_url"`
	} `json:"subscribed_services"`
}

func (c *UpdateProducerToConsumer) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminId", AdminId)
	log.Println(beego.AppConfig.String("loglevel"), "Info", "UpdateProducerToConsumer Start")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", " UpdateProducerToConsumer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", " UpdateProducerToConsumer Page Success")
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

	row, err := db.Db.Query(`SELECT
    producer_to_consumer.id,
    producer_to_consumer.producer_id,
    producer_to_consumer.consumer_id,
    producer_to_consumer.producer_subscribed_services,
    producer_to_consumer.status,
    producer_to_consumer.created_at,
    producers.producer_name,
    consumers.consumer_name
FROM producer_to_consumer
LEFT JOIN producers ON producers.id = producer_to_consumer.producer_id
LEFT JOIN consumers ON consumers.id = producer_to_consumer.consumer_id where producer_to_consumer.id=$1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get SubscriberToProducer data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerToConsumer data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("ProducerToConsumer data not found")
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

	data4, err := services.GetConsumerServices(data[0][2])
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("GetConsumerServices fetch Failed")
		return
	}
	var field3Array []Field3

	for i := 0; i < len(data4); i++ {
		if len(data4[i]) >= 2 {
			field3 := Field3{
				ServiceName: data4[i][0],
				ServiceUrl:  data4[i][1],
			}
			field3Array = append(field3Array, field3)
		} else {
			log.Printf("Insufficient data for row %d", i)
		}
	}

	//Response for Forntend

	responseData := map[string]interface{}{
		"Id":                    data[0][0],
		"ProducerName":          data[0][6],
		"ConsumerName":          data[0][7],
		"ProducerServices":      data[0][3],
		"Status":                status,
		"services_list":         field3Array,
		"ConsumerDomainAddress": data4[0][2],
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
func (c *UpdateProducerToConsumer) Post() {
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "UpdateProducerToConsumer update  Page Fail")
		} else {
			c.Data["DisplayMessage"] = "System User updated successfully"
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "UpdateProducerToConsumer update Page Success")
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

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "username :- ", username)

	uid := sess.Get("uid").(string)
	c.Data["uid"] = uid

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

	input_status := Updatevalues.Status

	input_subsc_list := Updatevalues.SubscribedServices

	fmt.Println("input_subsc_list :", input_subsc_list)

	type SubscribedService struct {
		ServiceName string `json:"service_name"`
		ServiceURL  string `json:"service_url"`
	}

	// Convert input_subsc_list to []SubscribedService
	var subscribedServices []SubscribedService
	for _, item := range input_subsc_list {
		subscribedServices = append(subscribedServices, SubscribedService{
			ServiceName: item.ServiceName,
			ServiceURL:  item.ServiceURL,
		})
	}

	// Create the struct to hold the "subsrcibed_services" key and the SubscribedService value
	jsonresult := struct {
		SubsrcibedServices []SubscribedService `json:"subsrcibed_services"`
	}{subscribedServices}

	// Marshal the struct to a JSON string
	jsonData, err := json.Marshal(jsonresult)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Print the JSON string

	fmt.Println(string(jsonData))

	var channelstatus bool

	if input_status == "ACTIVE" {

		channelstatus = true

	} else {
		channelstatus = false
	}

	res, err := db.Db.Exec(`UPDATE public."producer_to_consumer" SET producer_subscribed_services=$1,updated_by=$2,status=$3,updated_at=now() WHERE id = $4`, string(jsonData), uid, channelstatus, AdminId)
	if err != nil {
		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_PRODUCERTOCONSUMER_UPDATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("EN_PRODUCERTOCONSUMER_UPDATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_PRODUCERTOCONSUMER_UPDATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("FN_PRODUCERTOCONSUMER_UPDATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}

		return
	}
	i, err := res.RowsAffected()
	if err != nil || i == 0 {
		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_PRODUCERTOCONSUMER_UPDATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("EN_PRODUCERTOCONSUMER_UPDATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_PRODUCERTOCONSUMER_UPDATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("FN_PRODUCERTOCONSUMER_UPDATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}

		return
	}

	var pop_msg string

	if language == "french" {

		pop_msg = beego.AppConfig.String("FN_PRODUCERTOCONSUMER_UPDATED_SUCCESSFULLY")

	} else {

		pop_msg = beego.AppConfig.String("EN_PRODUCERTOCONSUMER_UPDATED_SUCCESSFULLY")

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
	result, err := db.Db.Exec("UPDATE public.'sysuser' set pass_count=$1 where id=$2 ", count, uname)
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

func sendFailureResponse(c *UpdateProducerToConsumer, message string) {
	c.Ctx.Output.JSON(map[string]interface{}{
		"title":   "FAILURE",
		"message": message,
		"status":  false,
	}, false, false)
}
