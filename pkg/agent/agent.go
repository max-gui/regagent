package agent

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/max-gui/consulagent/pkg/consulsets"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/regagent/pkg/ragcli"
	"github.com/max-gui/regagent/pkg/regagentsets"
)

func Consulagent(c *gin.Context) {

	appname := c.Param("appname")
	env := c.Param("env")
	uri := c.Request.RequestURI
	logger := logagent.Inst(c)
	logger.Print(uri)
	fixhealthstr := "/consulagent/" + appname + "/" + env + "/:8500/v1/health/service/" //, *regagentsets.Port)

	// if strings.Contains("/kv/")
	if !strings.Contains(uri, "/kv/") && !ragcli.IsSkip(appname, c) && strings.HasPrefix(uri, fixhealthstr) && c.Request.Method == http.MethodGet {
		service := strings.Split(strings.TrimPrefix(uri, fixhealthstr), "?")[0]
		consulapps := ragcli.GetConsulapps(appname, env, service, c)

		c.JSON(http.StatusOK, consulapps)

		// c.Redirect(http.StatusMovedPermanently, eutarget)
	} else {
		fixstr := "/consulagent/" + appname + "/" + env + "/:8500"
		target := *consulsets.Consul_host + strings.ReplaceAll(uri, fixstr, "")
		log.Print(target)
		log.Print("redirected")
		remote, err := url.Parse(*consulsets.Consul_host)
		if err != nil {
			fmt.Println(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(remote)
		c.Request.Host = remote.Host
		// c.Request.RemoteAddr = "user:eureka@eureka.kube.com"
		// c.Request.RequestURI = "/eureka" + strings.ReplaceAll(uri, "/eurekaagent", "")
		// c.Request.URL.User = remote.User
		c.Request.URL.Path = strings.Split(strings.ReplaceAll(uri, fixstr, ""), "?")[0]
		// c.Request.URL.RawQuery = strings.ReplaceAll(c.Request.URL.RawQuery, "token=", "nekot=")
		// auth := "user:eureka"
		// basicAuth := "Bearer " + base64.StdEncoding.EncodeToString([]byte(*regagentsets.Acltoken))
		// c.Request.Header.Add("Authorization", basicAuth)
		// if c.Request.Method == http.MethodPut {
		log.Print(c.Request)
		// }
		proxy.ServeHTTP(c.Writer, c.Request)
	}

}

func Eurekaagent(c *gin.Context) {

	appname := c.Param("appname")
	env := c.Param("env")

	uri := c.Request.RequestURI

	logger := logagent.Inst(c)
	logger.Info(uri)

	// fixstr := fmt.Sprintf("/consulagent:%s", *regagentsets.Port)
	if ragcli.IsSkip(appname, c) || (strings.TrimSuffix(uri, "/") != "/eurekaagent/"+appname+"/"+env+"/apps" && strings.TrimSuffix(uri, "/") != "/eurekaagent/apps/delta") || c.Request.Method != http.MethodGet {
		// eutarget := *regagentsets.Eu_host + strings.ReplaceAll(uri, "/eurekaagent", "")
		// logger.Print(eutarget)
		logger.Print("redirected")
		remote, err := url.Parse(strings.TrimSuffix(*regagentsets.Eu_host, "/eureka"))
		if err != nil {
			fmt.Println(err)
		}
		// remote.User = c.Request.URL.User
		// remote.Host = "user:eureka@" + remote.Host
		// remote.Path = "/eureka" + strings.ReplaceAll(uri, "/eurekaagent/"+appname, "")
		proxy := httputil.NewSingleHostReverseProxy(remote)
		// // http://user:eureka@eureka.kube.com/eureka
		// c.Request.Host = remote.Host // "eureka.kube.com"
		// // c.Request.RemoteAddr = "user:eureka@eureka.kube.com"
		// // c.Request.RequestURI = "/eureka" + strings.ReplaceAll(uri, "/eurekaagent", "")
		// // c.Request.URL.User = remote.User
		c.Request.URL.Path = strings.Split(strings.ReplaceAll(uri, "/eurekaagent/"+appname+"/"+env, "/eureka"), "?")[0]
		// c.Request.RequestURI = ""
		// pwd, _ := remote.User.Password()

		// auth := remote.User.Username() + ":" + pwd
		auth := "user:eureka"
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		for k := range c.Request.Header {

			c.Request.Header.Del(k)
		}
		c.Request.RemoteAddr = ""
		c.Request.RequestURI = ""
		c.Request.Host = ""
		// c.Request.URL.User = remote.User
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Add("Authorization", basicAuth)

		// if c.Request.Method == http.MethodPost {
		// log.Print(c.Request)
		logger.Print(c.Request)
		// }
		proxy.ServeHTTP(c.Writer, c.Request)
		// c.Redirect(http.StatusMovedPermanently, eutarget)
	} else {
		euapps := ragcli.GetEuapps(appname, env, c)

		c.JSON(http.StatusOK, euapps)
	}
}
