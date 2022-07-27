package ragcli

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/regagent/pkg/regagentsets"
)

type EurekaApplications struct {
	Applications struct {
		Versions__delta string          `json:"versions__delta"`
		Apps__hashcode  string          `json:"apps__hashcode"`
		Application     []Eurekaappinfo `json:"application"`
	} `json:"applications"`
}

func (euinst EurekaInstance) String() string {
	return fmt.Sprintf("name:%s,x-baggage-AF-env:%s,x-baggage-AF-region:%s,dc:%s,Address:%s,Port:%d,extaddress:%s,extport:%s",
		euinst.App,
		euinst.Metadata["x-baggage-AF-env"],
		euinst.Metadata["x-baggage-AF-region"],
		euinst.Metadata["dc"],
		euinst.IpAddr,
		euinst.Port.Realport,
		euinst.Metadata["extaddress"],
		euinst.Metadata["extport"])
}

// func (euinst EurekaInstance) Print() string {
// 	var logmess = fmt.Sprintf("x-baggage-AF-env:%s,x-baggage-AF-region:%s,dc:%s,Address:%s,Port:%d,extaddress:%s,extport:%s",
// 		euinst.Metadata["x-baggage-AF-env"],
// 		euinst.Metadata["x-baggage-AF-region"],
// 		euinst.Metadata["dc"],
// 		euinst.IpAddr,
// 		euinst.Port.Realport,
// 		euinst.Metadata["extaddress"],
// 		euinst.Metadata["extport"])

// 	return logmess
// }

type EurekaInstance struct {
	InstanceId       string `json:"instanceId"`
	HostName         string `json:"hostName"`
	App              string `json:"app"`
	Status           string `json:"status"`
	Overriddenstatus string `json:"overriddenstatus"`
	IpAddr           string `json:"ipAddr"`
	Port             struct {
		Realport int    `json:"$"`
		Enabled  string `json:"@enabled"`
	} `json:"port"`
	SecurePort struct {
		Realport int    `json:"$"`
		Enabled  string `json:"@enabled"`
	} `json:"securePort"`
	CountryId      int `json:"countryId"`
	DataCenterInfo struct {
		Class string `json:"@class"`
		Name  string `json:"name"`
	} `json:"dataCenterInfo"`
	LeaseInfo struct {
		RenewalIntervalInSecs int   `json:"renewalIntervalInSecs"`
		DurationInSecs        int   `json:"durationInSecs"`
		RegistrationTimestamp int64 `json:"registrationTimestamp"`
		LastRenewalTimestamp  int64 `json:"lastRenewalTimestamp"`
		EvictionTimestamp     int64 `json:"evictionTimestamp"`
		ServiceUpTimestamp    int64 `json:"serviceUpTimestamp"`
	} `json:"leaseInfo"`
	Metadata                      map[string]string `json:"metadata"`
	HomePageUrl                   string            `json:"homePageUrl"`
	StatusPageUrl                 string            `json:"statusPageUrl"`
	HealthCheckUrl                string            `json:"healthCheckUrl"`
	VipAddress                    string            `json:"vipAddress"`
	SecureVipAddress              string            `json:"secureVipAddress"`
	IsCoordinatingDiscoveryServer string            `json:"isCoordinatingDiscoveryServer"`
	LastUpdatedTimestamp          string            `json:"lastUpdatedTimestamp"`
	LastDirtyTimestamp            string            `json:"lastDirtyTimestamp"`
	ActionType                    string            `json:"actionType"`
}

type EurekaApplication struct {
	Application Eurekaappinfo `json:"application"`
}

type Eurekaappinfo struct {
	Name     string           `json:"name"`
	Instance []EurekaInstance `json:"instance"`
}

func GetEuapps(appname, env string, c context.Context) EurekaApplications {
	services := getAppsList(appname, "",
		func(servicesmap map[string]interface{}) {
			euapps := eurekapps(c)
			for _, v := range euapps.Applications.Application {

				servicesmap[strings.ToLower(v.Name)] = nil
			}
		}, c)

	// services := getAppsList(appname, "fake", c)

	logger := logagent.Inst(c)
	euapps := EurekaApplications{}
	euapps.Applications.Versions__delta = "1"
	euapps.Applications.Apps__hashcode = "DOWN_2_STARTING_11_UP_618_"
	agentport, err := strconv.Atoi(*regagentsets.AgentPort)
	if err != nil {
		logger.Panic(err)
	}
	for kname := range services {
		// kname := k.(string)
		eurekainst := EurekaInstance{}
		eurekainst.InstanceId = kname + strconv.Itoa(7979)
		eurekainst.Status = "UP"
		eurekainst.Overriddenstatus = "UP"
		eurekainst.ActionType = "ADDED"
		eurekainst.App = kname
		eurekainst.VipAddress = kname
		eurekainst.SecureVipAddress = kname
		eurekainst.Port.Realport = agentport
		eurekainst.HomePageUrl = "http://127.0.0.1:7979/eurekaagent"

		if *regagentsets.AgentPort == "80" {
			eurekainst.HostName = fmt.Sprintf("127.0.0.1/proxy/%s/%s/", kname, env)
			eurekainst.IpAddr = fmt.Sprintf("127.0.0.1/proxy/%s/%s/", kname, env)
		} else {
			eurekainst.HostName = fmt.Sprintf("127.0.0.1:%s/proxy/%s/%s/", *regagentsets.AgentPort, kname, env)
			eurekainst.IpAddr = fmt.Sprintf("127.0.0.1:%s/proxy/%s/%s/", *regagentsets.AgentPort, kname, env)
		}

		// eurekainst.HostName = fmt.Sprintf("127.0.0.1:%s/proxy/%s/%s/", *regagentsets.AgentPort, kname, env) // "127.0.0.1:7979/" + kname + "/"
		// eurekainst.IpAddr = fmt.Sprintf("127.0.0.1:%s/proxy/%s/%s/", *regagentsets.AgentPort, kname, env)   //"127.0.0.1:7979/" + kname + "/"
		eurekainst.DataCenterInfo = struct {
			Class string `json:"@class"`
			Name  string `json:"name"`
		}{Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo", Name: "MyOwn"}
		eurekainst.Metadata = map[string]string{
			"x-baggage-AF-env":    env, //c.Value("env").(string),
			"x-baggage-AF-region": "default",
			"source":              "consul",
		}

		euapps.Applications.Application = append(euapps.Applications.Application, Eurekaappinfo{Name: kname, Instance: []EurekaInstance{eurekainst}})

	}
	// consulapps := consulhelp.GetServices(c)
	// for k := range consulapps {
	// 	appname := strings.ToUpper(k)
	// 	if _, ok := servicemap[appname]; !ok {

	// 		eurekainst := EurekaInstance{}
	// 		eurekainst.InstanceId = k + strconv.Itoa(7979)
	// 		eurekainst.Status = "UP"
	// 		eurekainst.App = appname
	// 		eurekainst.VipAddress = appname
	// 		eurekainst.SecureVipAddress = appname
	// 		eurekainst.Port.Realport = 7979
	// 		eurekainst.HomePageUrl = "http://127.0.0.1:7979/eurekaagent"
	// 		eurekainst.HostName = "127.0.0.1:7979/" + k + "/"
	// 		eurekainst.IpAddr = "127.0.0.1:7979/" + k + "/"
	// 		eurekainst.DataCenterInfo = struct {
	// 			Class string `json:"@class"`
	// 			Name  string `json:"name"`
	// 		}{Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo", Name: "MyOwn"}
	// 		eurekainst.Metadata = map[string]string{
	// 			"x-baggage-AF-env":    c.Value("env").(string),
	// 			"x-baggage-AF-region": "default",
	// 			"source":              "consul",
	// 		}

	// 		euapps.Applications.Application = append(euapps.Applications.Application, Eurekaappinfo{Name: appname, Instance: []EurekaInstance{eurekainst}})
	// 	}
	// }
	return euapps
}

func tocall(method, url string, heads map[string]string, c context.Context) *http.Response {
	logger := logagent.Inst(c)
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	req, _ := http.NewRequest(method, url, nil)
	for k, v := range heads {
		req.Header.Add(k, v)
	}
	response, err := netClient.Do(req)
	if err != nil {
		logger.Panic(err)
	}
	return response
}
func eurekapps(c context.Context) EurekaApplications {
	logger := logagent.Inst(c)
	resp := tocall("GET", *regagentsets.Eu_host+"/apps/", map[string]string{"Accept": "application/json"}, c)
	// resp := tocall("GET", "http://user:eureka@eureka.kube.com/eureka/apps/"+servicename, map[string]string{"Accept": "application/json"})
	resbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	var resjson = EurekaApplications{}
	err = json.Unmarshal(resbody, &resjson)
	// logger.Info(resjson.Application.Instance[0].HomePageUrl)
	// logger.Info(resjson.Application.Instance[0].IpAddr)
	// logger.Info(resjson.Application.Instance[0].Metadata)
	// logger.Info(resjson.Application.Instance[0].Port.Realport)
	if err != nil {
		logger.Panic(err)
	}
	return resjson
}

func Eurekapp(servicename string, c context.Context) EurekaApplication {
	logger := logagent.Inst(c)
	resp := tocall("GET", *regagentsets.Eu_host+"/apps/"+servicename, map[string]string{"Accept": "application/json"}, c)
	// resp := tocall("GET", "http://user:eureka@eureka.kube.com/eureka/apps/"+servicename, map[string]string{"Accept": "application/json"})
	resbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	var resjson = EurekaApplication{}
	err = json.Unmarshal(resbody, &resjson)
	//logger.Info(resjson.Application.Instance[0].HomePageUrl)
	//logger.Info(resjson.Application.Instance[0].IpAddr)
	//logger.Info(resjson.Application.Instance[0].Metadata)
	//logger.Info(resjson.Application.Instance[0].Port.Realport)
	if err != nil {
		logger.Print(err)
	}
	return resjson
}
