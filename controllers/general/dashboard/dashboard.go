package dashboard

import (
	"SuperEsbAdminWeb/session"
	//	"SuperEsbAdminWeb/utils"
	"errors"
	"runtime/debug"

	"SuperEsbAdminWeb/utils/database/sql"

	"SuperEsbAdminWeb/model/db"
	"strconv"

	"SuperEsbAdminWeb/model/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	//"time"

	"github.com/astaxie/beego"

	log "github.com/sirupsen/logrus"
)

type Dashboard struct {
	beego.Controller
}

type Display struct {
	SystemPool       []SystemPool
	SystemCommision  []SystemCommision
	BillerPool       []BillerPool
	BillerCommision  []BillerCommision
	PGWPool          []PGWPool
	PGWPoolCommision []PGWPoolCommision
}
type SystemPool struct {
	UserName string
	Balance  string
}
type SystemCommision struct {
	UserName string
	Balance  string
}
type BillerPool struct {
	UserName string
	Balance  string
}
type BillerCommision struct {
	UserName string
	Balance  string
}
type PGWPool struct {
	UserName string
	Balance  string
}
type PGWPoolCommision struct {
	UserName string
	Balance  string
}

type ProducerToConsumer struct {
	Producer           string   `json:"producer"`
	Consumer           string   `json:"consumer"`
	SubscribedServices []string `json:"subscribedServices"`
}

type ConsumerUsage struct {
	ConsumerName      string `json:"consumer_name"`
	ConsumerColorCode string `json:"consumer_color_code"`
	ServicesData      []struct {
		ServiceName  string `json:"service_name"`
		ServiceCount string `json:"service_count"`
	} `json:"services_data"`
}

// GenerateRandomHexColor generates a random hex color code
func GenerateRandomHexColor() string {
	rand.Seed(time.Now().UnixNano())

	// Generate random values for red, green, and blue components
	red := rand.Intn(128)        // Bias towards darker shades by limiting the range
	green := rand.Intn(128)      // Bias towards darker shades by limiting the range
	blue := rand.Intn(128) + 128 // Bias towards darker shades by starting from a higher base

	// Format the color code
	color := fmt.Sprintf("#%02x%02x%02x", red, green, blue)
	return color
}

type searchData struct {
	CustomStartDate string `json:"customStartDate"`
	CustomEndDate   string `json:"customEndDate"`
}

func (c *Dashboard) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Dashboard Page Start")
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
			c.TplName = "index.html"
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)

			} else {
				c.Data["DisplayMessage"] = err.Error()
			}

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Dashboard Page Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Dashboard Page Page Success")
		}
		return
	}()

	//utils.SetHTTPHeader(c.Ctx)

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

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "UserName - ", sess.Get("uname"))
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Session ID - ", sess.SessionID())

	uname := sess.Get("uname").(string)

	username := sess.Get("username").(string)
	c.Data["uname"] = uname

	language := sess.Get("language").(string)
	c.Data["language"] = language

	fmt.Println("language", language)

	row8, err := db.Db.Query(`SELECT password_set,password_updated_date FROM sysuser where  email=$1`, uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get System User password set value ")
		log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)

		return
	}
	defer sql.Close(row8)
	_, data8, err := sql.Scan(row8)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User password set value scan fail")
		log.Println(beego.AppConfig.String("loglevel"), "ENGLISH Error", err)

		return
	}

	s1 := data8[0][0]
	b1, _ := strconv.ParseBool(s1)

	var sysuser_pwd_set string

	if b1 == false {

		fmt.Println("Password is not set by system user")
		sysuser_pwd_set = "false"

	} else {
		fmt.Println("Password is set by system user")
		sysuser_pwd_set = "true"

	}

	//------------for system user count-----------

	row, err := db.Db.Query(`SELECT count(*)as sysusercount FROM sysuser where status='true'`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get System User Count data")
		return
	}
	defer sql.Close(row)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get  System User Count data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "sysusercount Query Data - ", data, "Data len - ", len(data))

	number, err := strconv.ParseUint(data[0][0], 10, 32)
	finalIntNum := int(number)

	c.Data["activesysusercount"] = finalIntNum

	row1, err := db.Db.Query(`SELECT count(*)as sysusercount FROM sysuser where status='false'`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get System User Count data")
		return
	}
	defer sql.Close(row1)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row1)
	_, data1, err := sql.Scan(row1)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get  System User Count data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "sysusercount Query Data - ", data1, "Data len - ", len(data1))

	number1, err := strconv.ParseUint(data1[0][0], 10, 32)
	finalIntNum1 := int(number1)

	c.Data["inactivesysusercount"] = finalIntNum1

	//------------for producer count-----------
	row2, err := db.Db.Query(`SELECT count(*)as sysusercount FROM producers where status='true'`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get producers  Count data")
		return
	}
	defer sql.Close(row2)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row2)
	_, data2, err := sql.Scan(row2)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get  producers Count data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "producers Query Data - ", data2, "Data len - ", len(data2))

	number2, err := strconv.ParseUint(data2[0][0], 10, 32)
	finalIntNum2 := int(number2)

	c.Data["activeproducercount"] = finalIntNum2

	row12, err := db.Db.Query(`SELECT count(*)as sysusercount FROM producers where status='false'`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get producers  Count data")
		return
	}
	defer sql.Close(row12)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row12)
	_, data12, err := sql.Scan(row12)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get  producers Count data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "producers Query Data - ", data12, "Data len - ", len(data12))

	number12, err := strconv.ParseUint(data12[0][0], 10, 32)
	finalIntNum12 := int(number12)

	c.Data["inactiveproducercount"] = finalIntNum12

	//------------for consumer count-----------

	row3, err := db.Db.Query(`SELECT count(*)as sysusercount FROM consumers where status='true'`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get consumers Count data")
		return
	}
	defer sql.Close(row3)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
	_, data3, err := sql.Scan(row3)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get  consumers Count data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "sysusercount Query Data - ", data3, "Data len - ", len(data3))

	number3, err := strconv.ParseUint(data3[0][0], 10, 32)
	finalIntNum3 := int(number3)

	c.Data["activeconsumercount"] = finalIntNum3

	row13, err := db.Db.Query(`SELECT count(*)as sysusercount FROM consumers where status='false'`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get consumers Count data")
		return
	}
	defer sql.Close(row13)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row13)
	_, data13, err := sql.Scan(row13)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get  consumers Count data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "sysusercount Query Data - ", data13, "Data len - ", len(data13))

	number13, err := strconv.ParseUint(data13[0][0], 10, 32)
	finalIntNum13 := int(number13)

	c.Data["inactiveconsumercount"] = finalIntNum13

	//------------Sending Producer to Consumer Mapping details -----------

	row4, err := db.Db.Query(`SELECT
	  producer_to_consumer.id,
    producers.producer_name,
    consumers.consumer_name,
    producer_to_consumer.producer_subscribed_services,
     consumers.consumer_address
FROM producer_to_consumer
LEFT JOIN producers ON producers.id = producer_to_consumer.producer_id
LEFT JOIN consumers ON consumers.id = producer_to_consumer.consumer_id
 order by producers.producer_name `)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get consumers Count data")
		return
	}
	defer sql.Close(row4)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row4)
	_, data4, err := sql.Scan(row4)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get  Producer to Consumer Mapping details")
		return
	}

	// var producerToConsumerList []ProducerToConsumer

	// // Iterate through data4 and construct the ProducerToConsumer slice
	// for i := range data4 {
	// 	producer := data4[i][0]
	// 	consumer := data4[i][1]

	// 	// Parse the JSON data in data4[i][2]
	// 	var consumerInfo struct {
	// 		ConsumerSubscribedSvcs []struct {
	// 			ServiceName string `json:"service_name"`
	// 		} `json:"subsrcibed_services"`
	// 	}

	// 	if err := json.Unmarshal([]byte(data4[i][2]), &consumerInfo); err != nil {
	// 		log.Printf("Error parsing JSON for row %d: %v", i, err)
	// 		continue // Skip this row if JSON parsing fails
	// 	}

	// 	// Extract the service names and add to SubscribedServices
	// 	var subscribedServices []string
	// 	for _, svc := range consumerInfo.ConsumerSubscribedSvcs {
	// 		subscribedServices = append(subscribedServices, svc.ServiceName)
	// 	}

	// 	// Append the data to the ProducerToConsumer slice
	// 	p2c := ProducerToConsumer{
	// 		Producer:           producer,
	// 		Consumer:           consumer,
	// 		SubscribedServices: subscribedServices,
	// 	}
	// 	producerToConsumerList = append(producerToConsumerList, p2c)
	// }

	// // Marshal the result into JSON and log it
	// p2cJSON, err := json.Marshal(producerToConsumerList)
	// if err != nil {
	// 	log.Println("Error marshaling JSON:", err)
	// 	return
	// }

	// fmt.Println(string(p2cJSON))

	type Row struct {
		Id                    string
		ProducerName          string
		ConsumerName          string
		Timestamp             string
		Status                string
		ProducerServices      string
		ConsumerDomainAddress string
	}

	var result []Row

	for i := range data4 {
		var r Row

		r.Id = data4[i][0]
		r.ProducerName = data4[i][1]
		r.ConsumerName = data4[i][2]
		r.ProducerServices = data4[i][3]
		r.ConsumerDomainAddress = data4[i][4]

		result = append(result, r)

	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Producer to Consumer Mapping details Query Data - ", result, "Data len - ", len(result))

	menu := sess.Get("menu").(string)
	submenu := sess.Get("submenu").(string)

	//The graph data

	row5, err := db.Db.Query(`SELECT COUNT(esb_request_metadata.service), consumers.id, consumers.consumer_name, esb_request_metadata.service
		FROM esb_request_metadata
		LEFT JOIN consumers ON consumers.id = esb_request_metadata.consumer_id
		WHERE consumer_id IS NOT NULL
		GROUP BY consumers.id, consumers.consumer_name, esb_request_metadata.service
		ORDER BY consumers.consumer_name`)

	if err != nil {
		// Handle error
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Consumer Usage Count data")
		return
	}
	defer sql.Close(row5)
	_, data5, err := sql.Scan(row5)
	if err != nil {
		// Handle error
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Consumer Usage details")
		return
	}

	var consumerUsageList []ConsumerUsage
	currentConsumer := ConsumerUsage{} // Track the current consumer during iteration

	for i := range data5 {
		consumerName := data5[i][2]
		serviceName := data5[i][3]
		serviceCount := data5[i][0]

		if currentConsumer.ConsumerName == "" {
			// First iteration, initialize the current consumer
			currentConsumer.ConsumerName = consumerName
			currentConsumer.ConsumerColorCode = GenerateRandomHexColor()
		}

		if consumerName != currentConsumer.ConsumerName {
			// New consumer detected, add the current consumer to the list
			consumerUsageList = append(consumerUsageList, currentConsumer)
			// Reset the current consumer
			currentConsumer = ConsumerUsage{
				ConsumerName:      consumerName,
				ConsumerColorCode: GenerateRandomHexColor(),
			}
		}

		// Append the service data to the current consumer
		currentConsumer.ServicesData = append(currentConsumer.ServicesData, struct {
			ServiceName  string `json:"service_name"`
			ServiceCount string `json:"service_count"`
		}{
			ServiceName:  serviceName,
			ServiceCount: serviceCount,
		})
	}

	// Append the last consumer to the list
	if currentConsumer.ConsumerName != "" {
		consumerUsageList = append(consumerUsageList, currentConsumer)
	}

	// Marshal the result into JSON and log it
	consumerUsageJSON, err := json.Marshal(consumerUsageList)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Consumer Usage Data - ", string(consumerUsageJSON))

	//Response for Forntend

	responseData := map[string]interface{}{
		"ActiveSysuserCount":    c.Data["activesysusercount"],
		"InActiveSysuserCount":  c.Data["inactivesysusercount"],
		"ActiveProducerCount":   c.Data["activeproducercount"],
		"InActiveProducerCount": c.Data["inactiveproducercount"],
		"ActiveConsumercount":   c.Data["activeconsumercount"],
		"InActiveConsumercount": c.Data["inactiveconsumercount"],
		"SysUserPasswordSet":    sysuser_pwd_set,
		"username":              username,
		"CustomerData":          result,
		"menu":                  menu,
		"submenu":               submenu,
		"ConsumerUsageData":     consumerUsageList,
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

func (c *Dashboard) Post() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "node post page")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search SystemUser Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "index.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search SystemUser  Page Success")
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

	language := sess.Get("language").(string)
	c.Data["language"] = language

	fmt.Println("language", language)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "username :- ", username)
	var Searchvalues searchData
	if err := json.NewDecoder(c.Ctx.Request.Body).Decode(&Searchvalues); err != nil {
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

	from := Searchvalues.CustomStartDate
	to := Searchvalues.CustomEndDate

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "From Date - ", from)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "to Date - ", to)

	row5, err := db.Db.Query(`SELECT COUNT(esb_request_metadata.service), 
    consumers.id, 
    consumers.consumer_name, 
    esb_request_metadata.service
	FROM 
	esb_request_metadata
	LEFT JOIN consumers ON consumers.id = esb_request_metadata.consumer_id
	WHERE consumer_id IS NOT NULL
    AND (TO_DATE(esb_request_metadata.created_at::text,'YYYY/MM/DD') >= TO_DATE( $1,'MM/DD/YYYY')) 
    AND (TO_DATE(esb_request_metadata.created_at::text,'YYYY/MM/DD') <= TO_DATE( $2,'MM/DD/YYYY'))
	GROUP BY consumers.id, consumers.consumer_name, esb_request_metadata.service
	ORDER BY consumers.consumer_name`, from, to)

	if err != nil {
		// Handle error
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Consumer Usage Count data")
		return
	}
	defer sql.Close(row5)
	_, data5, err := sql.Scan(row5)
	if err != nil {
		// Handle error
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Consumer Usage details")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "ConsumerUsageDataPostMethod Query Data - ", data5, "Data len - ", len(data5))

	var consumerUsageList []ConsumerUsage
	currentConsumer := ConsumerUsage{} // Track the current consumer during iteration

	for i := range data5 {
		consumerName := data5[i][2]
		serviceName := data5[i][3]
		serviceCount := data5[i][0]

		if currentConsumer.ConsumerName == "" {
			// First iteration, initialize the current consumer
			currentConsumer.ConsumerName = consumerName
			currentConsumer.ConsumerColorCode = GenerateRandomHexColor()
		}

		if consumerName != currentConsumer.ConsumerName {
			// New consumer detected, add the current consumer to the list
			consumerUsageList = append(consumerUsageList, currentConsumer)
			// Reset the current consumer
			currentConsumer = ConsumerUsage{
				ConsumerName:      consumerName,
				ConsumerColorCode: GenerateRandomHexColor(),
			}
		}

		// Append the service data to the current consumer
		currentConsumer.ServicesData = append(currentConsumer.ServicesData, struct {
			ServiceName  string `json:"service_name"`
			ServiceCount string `json:"service_count"`
		}{
			ServiceName:  serviceName,
			ServiceCount: serviceCount,
		})
	}

	// Append the last consumer to the list
	if currentConsumer.ConsumerName != "" {
		consumerUsageList = append(consumerUsageList, currentConsumer)
	}
	fmt.Println("Length of consumerUsageList", len(consumerUsageList))

	// Marshal the result into JSON and log it
	consumerUsageJSON, err := json.Marshal(consumerUsageList)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Consumer Usage Data - ", string(consumerUsageJSON))

	//Response for Forntend

	responseData := map[string]interface{}{
		"ConsumerUsageDataPostMethod": consumerUsageList,
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

func sendFailureResponse(c *Dashboard, message string) {
	c.Ctx.Output.JSON(map[string]interface{}{
		"title":   "FAILURE",
		"message": message,
		"status":  false,
	}, false, false)
}
