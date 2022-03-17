package ragcli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/regagent/pkg/regagentsets"
)

func GetConsulapps(appname, env, service string, c context.Context) []*api.ServiceEntry {

	logger := logagent.Inst(c)

	services := getAppsList(appname, "",
		func(servicesmap map[string]interface{}) {
			servicesmap[strings.ToLower(service)] = nil
			// return map[string]interface{}{strings.ToLower(service): nil}
		}, c)

	// services := getAppsList(appname, service, c)
	// service := strings.Split(strings.TrimPrefix(uri, fixhealthstr), "?")[0]

	consulapps := []*api.ServiceEntry{}
	if _, ok := services[strings.ToLower(service)]; ok {
		entry := api.ServiceEntry{}
		agservice := api.AgentService{}
		if *regagentsets.AgentPort == "80" {
			agservice.Address = fmt.Sprintf("127.0.0.1/proxy/%s/%s/", service, env)
		} else {
			agservice.Address = fmt.Sprintf("127.0.0.1:%s/proxy/%s/%s/", *regagentsets.AgentPort, service, env)
		}
		// agservice.Address = fmt.Sprintf("127.0.0.1:%s/proxy/%s/%s/", *regagentsets.AgentPort, service, env)
		port, err := strconv.Atoi(*regagentsets.AgentPort)
		if err != nil {
			logger.Panic(err)
		}
		agservice.Port = port
		agservice.Service = service
		// x-baggage-AF-env"))
		// c.Set("region", c.Request.Header.Get("x-baggage-AF-region"))
		agservice.Tags = []string{
			fmt.Sprintf("x-baggage-AF-env:%s=x-baggage-AF-env:%s",
				env,  //c.Value("env"),
				env), //c.Value("env")),
			"x-baggage-AF-region:default=x-baggage-AF-region:default"}
		agservice.Meta = map[string]string{
			"x-baggage-AF-env":    env, //c.Value("env").(string),
			"x-baggage-AF-region": "default",
		}
		entry.Service = &agservice

		check := api.HealthCheck{}
		check.Node = "LFB-L0515490"
		check.CheckID = "serfHealth"
		check.Status = "passing"
		entry.Checks = api.HealthChecks{&check}
		consulapps = append(consulapps, &entry)
		// consulapps = append() []*api.ServiceEntry{&entry}
		// {
		// 	"Node": "LFB-L0515490",
		// 	"CheckID": "serfHealth",
		// 	"Name": "Serf Health Status",
		// 	"Status": "passing",
		// 	"Notes": "",
		// 	"Output": "Agent alive and reachable",
		// 	"ServiceID": "",
		// 	"ServiceName": "",
		// 	"ServiceTags": [],
		// 	"Type": "",
		// 	"Definition": {},
		// 	"CreateIndex": 236385951,
		// 	"ModifyIndex": 236385951
		//   }

		// consulapps := consulhelp.GetHealthService(service)

		// if len(consulapps) > 0 {
		// 	entry := consulapps[0]
		// 	entry.Service.Address = "127.0.0.1:7979/" + service + "/"
		// 	entry.Service.Port = 7979
		// 	consulapps = []*api.ServiceEntry{entry}
		// }

		// euapps := eurekaapps()
		// for index, entry := range consulapps {
		// 	entry.Service.Address = "127.0.0.1:7979/" + service + "/"
		// 	entry.Service.Port = 7979
		// 	consulapps[index] = entry
		// 	// http://127.0.0.1:9999/eurekaagent
		// 	// for instindex, instval := range appval.Instance {
		// 	// 	instval.Port.Realport = 7979
		// 	// 	instval.HomePageUrl = "http://127.0.0.1:7979/eurekaagent"
		// 	// 	instval.HostName = "127.0.0.1:7979/" + instval.App + "/" //"127.0.0.1"
		// 	// 	instval.IpAddr = "127.0.0.1:7979/" + instval.App + "/"   //"127.0.0.1"
		// 	// 	euapps.Applications.Application[appindex].Instance[instindex] = instval
		// 	// }
		// }

		// var region, env string
		// if val, ok := c.Request.Header[strings.Title(strings.ToLower("x-baggage-AF-region"))]; ok {
		// 	region = val[0]
		// }
		// if val, ok := c.Request.Header[strings.Title(strings.ToLower("x-baggage-AF-env"))]; ok {
		// 	env = val[0]
		// }
		// log.Printf("region:%s", region)
		// log.Printf("env:%s", env)

		// for index, val := range euapps.Applications.Application {
		// 	euapps.Applications.Application[index].Instance = dede(val.Instance, region, env)
		// }
		// }
	}

	return consulapps
}
