package createProducerToConsumer

import (
	"SuperEsbAdminWeb/session"

	"encoding/json"
	"net/http"

	"net/smtp"
	"runtime/debug"
	"strings"

	//	"time"

	"SuperEsbAdminWeb/model/db"
	"SuperEsbAdminWeb/utils"
	"errors"

	"github.com/astaxie/beego"

	"SuperEsbAdminWeb/services"

	// "proyava.com/database/sql"
	// "proyava.com/util/log"

	"SuperEsbAdminWeb/utils/database/sql"
	"io/ioutil"
	"path/filepath"

	"fmt"

	//	"SuperEsbAdminWeb/utils/util/password"
	//	"SuperEsbAdminWeb/utils/util/txnno"
	// "SuperEsbAdminWeb/utils/util/password"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type CreateProducerToConsumer struct {
	beego.Controller
}
type GetConsumersServicelist struct {
	beego.Controller
}

type unencryptedAuth struct {
	smtp.Auth
}

type Field struct {
	Id    string
	Name  string
	Email string
}

type Display struct {
	Fields1 []Field1
}

type Display3 struct {
	Field3 []Field3
}
type Field1 struct {
	Id   string
	Name string
}

type Field3 struct {
	ServiceName string `json:"service_name"`
	ServiceUrl  string `json:"service_url"`
}
type Field2 struct {
	Id   string
	Name string
}
type createData struct {
	Producer           string `json:"input_producer"`
	Consumer           string `json:"input_consumer"`
	ProdServices       string `json:"producerservices"`
	Status             string `json:"input_status"`
	SubscribedServices []struct {
		ServiceName string `json:"service_name"`
		ServiceURL  string `json:"service_url"`
	} `json:"subscribed_services"`
}

type List []struct {
	ServiceName string `json:"service_name"`
	ServiceURL  string `json:"service_url"`
}

type getconsumerId struct {
	ConsumerId string `json:"consumer_id"`
}

type Service struct {
	ServiceURL string `json:"service_url"`
}

func (c *CreateProducerToConsumer) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "CreateProducerToConsumer Start")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "CreateProducerToConsumer  Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "CreateProducerToConsumer  Page Success")
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

	uname := sess.Get("uname").(string)
	c.Data["uname"] = uname

	data2, err := services.GetActiveConsumer()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("SearchConsumers fetch Failed")
		return
	}
	var Dis2 Display
	for i := range data2 {
		var d Field1
		d.Id = data2[i][0]
		d.Name = data2[i][1]
		Dis2.Fields1 = append(Dis2.Fields1, d)
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis2)

	data3, err := services.GetActiveProducer()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("SearchProducers  fetch Failed")
		return
	}
	var Dis1 Display
	for i := range data3 {
		var d Field1
		d.Id = data3[i][0]
		d.Name = data3[i][1]
		Dis1.Fields1 = append(Dis1.Fields1, d)
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis1)

	responseData := map[string]interface{}{
		"ProducerData": Dis1,
		"ConsumerData": Dis2,
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
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			c.Ctx.Output.SetStatus(500)
			c.Ctx.Output.Body([]byte("Error loading index.html"))
			return
		}

		c.Ctx.Output.Header("Content-Type", "text/html")
		c.Ctx.Output.Body(content)
	}

	return
}
func (c *CreateProducerToConsumer) Post() {
	//	var systemusermsg string

	log.Println(beego.AppConfig.String("loglevel"), "Info", "CreateProducerToConsumer page")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "CreateProducerToConsumer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "System User created Successfully"
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "CreateProducerToConsumer  Page Success")
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)
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

	uid := sess.Get("uid").(string)
	c.Data["uid"] = uid
	language := sess.Get("language").(string)
	c.Data["language"] = language

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user email :- ", username)

	log.Printf("Request Body: %s", string(c.Ctx.Input.RequestBody))

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

	input_prod := createvalues.Producer
	input_consumer := createvalues.Consumer
	// input_prodservices := createvalues.ProdServices
	input_status := createvalues.Status

	input_subsc_list := createvalues.SubscribedServices

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

	err = CheckUserAlreadyExists(input_prod, input_consumer)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_PRODUCERTOCONSUMER_ALREADY_EXISTS"))
			sendFailureResponse(c, beego.AppConfig.String("EN_PRODUCERTOCONSUMER_ALREADY_EXISTS"))
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_PRODUCERTOCONSUMER_ALREADY_EXISTS"))
			sendFailureResponse(c, beego.AppConfig.String("FN_PRODUCERTOCONSUMER_ALREADY_EXISTS"))
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}

		return
	}
	id := uuid.New()
	fmt.Println(id.String())

	result, err := db.Db.Exec(`INSERT INTO public."producer_to_consumer" (id,
	 producer_id,
	 consumer_id,
	 producer_subscribed_services, 
	status,
	 created_by,
	 created_at)
		VALUES ($1, $2, $3, $4,$5,$6,now())`,
		id,
		input_prod,
		input_consumer,
		string(jsonData),
		channelstatus,
		uid)
	if err != nil {
		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_PRODUCERTOCONSUMER_CREATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("EN_PRODUCERTOCONSUMER_CREATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_PRODUCERTOCONSUMER_CREATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("FN_PRODUCERTOCONSUMER_CREATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}
	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		if language == "english" {
			err = errors.New(beego.AppConfig.String("EN_PRODUCERTOCONSUMER_CREATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("EN_PRODUCERTOCONSUMER_CREATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)
			return
		} else if language == "french" {
			err = errors.New(beego.AppConfig.String("FN_PRODUCERTOCONSUMER_CREATION_FAILED"))
			sendFailureResponse(c, beego.AppConfig.String("FN_PRODUCERTOCONSUMER_CREATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "FRENCH Error", err)
			return
		}
		return
	}

	var pop_msg string

	if language == "french" {

		pop_msg = beego.AppConfig.String("FN_PRODUCERTOCONSUMER_CREATED_SUCCESSFULLY")

	} else {

		pop_msg = beego.AppConfig.String("EN_PRODUCERTOCONSUMER_CREATED_SUCCESSFULLY")

	}

	c.Ctx.Output.JSON(map[string]interface{}{
		"message": pop_msg,
	}, true, false)

	responseMap1 := map[string]interface{}{
		"success": true, // Indicate success in the response
		"message": "Login successful",
	}

	responseData1, _ := json.Marshal(responseMap1)

	c.Ctx.Output.Status = http.StatusOK
	c.Ctx.Output.Header("Content-Type", "application/json")
	c.Ctx.Output.Body(responseData1)

	return
}

func CheckUserAlreadyExists(prodid, consid string) (err error) {

	err = nil

	row, err := db.Db.Query(`SELECT count(*) FROM public."producer_to_consumer" WHERE producer_id = $1 AND consumer_id = $2`, prodid, consid)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch ProducerToConsumer")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch ProducerToConsumer")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	countlen := data[0][0]

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

	if countlen != "0" {
		err = errors.New("ProducerToConsumer already exists")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}

func sendFailureResponse(c *CreateProducerToConsumer, message string) {
	c.Ctx.Output.JSON(map[string]interface{}{
		"title":   "FAILURE",
		"message": message,
		"status":  false,
	}, false, false)
}

func (c *GetConsumersServicelist) Post() {
	//	var systemusermsg string

	log.Println(beego.AppConfig.String("loglevel"), "Info", "GetConsumersServicelist page")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "GetConsumersServicelist Page Fail")
		} else {
			c.Data["DisplayMessage"] = "System User created Successfully"
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "GetConsumersServicelist  Page Success")
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)
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

	uid := sess.Get("uid").(string)
	c.Data["uid"] = uid

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user email :- ", username)

	log.Printf("Request Body: %s", string(c.Ctx.Input.RequestBody))

	var createvalues getconsumerId
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

	input_consumer_id := createvalues.ConsumerId
	fmt.Println("Consumer id got", input_consumer_id)

	data4, err := services.GetConsumerServices(input_consumer_id)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("SearchConsumers fetch Failed")
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
	c.Ctx.Output.JSON(map[string]interface{}{
		"message":               "GetConsumersServicelist  got successfully",
		"services_list":         field3Array,
		"ConsumerDomainAddress": data4[0][2],
	}, true, false)

	responseMap1 := map[string]interface{}{
		"success": true, // Indicate success in the response
		"message": "Login successful",
	}

	responseData1, _ := json.Marshal(responseMap1)

	c.Ctx.Output.Status = http.StatusOK
	c.Ctx.Output.Header("Content-Type", "application/json")
	c.Ctx.Output.Body(responseData1)

}
