package regagentsets

import (
	"flag"
	"os"
	"strings"

	"github.com/max-gui/consulagent/pkg/consulsets"
	"gopkg.in/yaml.v2"
)

const PthSep = string(os.PathSeparator)

var (
	AgentPort, Eu_host, ConfWatchPrefix, ConfArchPrefix, ConfSkipAppPrefix *string
)

func StartupInit(bytes []byte) {

	confmap := map[string]interface{}{}
	yaml.Unmarshal(bytes, confmap)
	// *consulsets.Acltoken = confmap["af-arch"].(map[interface{}]interface{})["resource"].(map[interface{}]interface{})["private"].(map[interface{}]interface{})["acl-token"].(string)
	acltoken := confmap["af-arch"].(map[interface{}]interface{})["resource"].(map[interface{}]interface{})["private"].(map[interface{}]interface{})["acl-token"].(string)
	consulsets.StartupInit(acltoken)
	eu_hosts := confmap["af-arch"].(map[interface{}]interface{})["resource"].(map[interface{}]interface{})["agent"].(map[interface{}]interface{})["client"].(map[interface{}]interface{})["serviceUrl"].(map[interface{}]interface{})["defaultZone"].(string)
	*Eu_host = strings.Split(eu_hosts, "/,")[0]

}

func init() {

	// // logsets. Appname = "regagent"
	// Apppath = flag.String("apppath", "/Users/max/Downloads/regagent", "app root path")

	Eu_host = flag.String("euhost", "http://user:eureka@eureka.kube.com/eureka", "eureka url")
	// Port = flag.String("port", "9999", "this app's port")
	AgentPort = flag.String("agentport", "7979", "call agent's port")
	ConfWatchPrefix = flag.String("ConfWatchPrefix", "ops/", "watch prefix for consul")
	ConfArchPrefix = flag.String("confArchPrefix", "ops/iac/arch/", "arch prefix for consul")
	ConfSkipAppPrefix = flag.String("confSkipAppPrefix", "ops/charon/skip/", "arch prefix for consul")
	// Appenv = flag.String("appenv", "prod", "this application's working env")
	// Jsonlog = flag.Bool("jsonlog", false, "jsonlog or not")

	// log.Print("in")
	// flag.Parse()

}

// var Reppath = func() string {
// 	return *Apppath + PthSep + *Repopathname + PthSep
// }
