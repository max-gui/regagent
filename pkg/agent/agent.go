package agent

import (
	"encoding/base64"
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

func Consulagent8500(c *gin.Context) {

	Consulagent(c, "/consulagent/", "/:8500")

}

func Consulagent80(c *gin.Context) {

	Consulagent(c, "/consul/", "")

}

func Consulagent(c *gin.Context, leadingname, portfix string) {

	appname := c.Param("appname")
	env := c.Param("env")
	uri := c.Request.RequestURI
	logger := logagent.InstArch(c)
	logger.Print(uri)
	fixhealthstr := leadingname + appname + "/" + env + portfix + "/v1/health/service/"           //, *regagentsets.Port)
	fixcatalogservicesstr := leadingname + appname + "/" + env + portfix + "/v1/catalog/services" //, *regagentsets.Port)
	fixstatusleaderstr := leadingname + appname + "/" + env + portfix + "/v1/status/leader"       //, *regagentsets.Port)
	fixproxystr := leadingname + appname + "/" + env + portfix

	hostmode := leadingname == "/consul/"
	// fixhealthstr := "/consulagent/" + appname + "/" + env + "/:8500/v1/health/service/"           //, *regagentsets.Port)
	// fixcatalogservicesstr := "/consulagent/" + appname + "/" + env + "/:8500/v1/catalog/services" //, *regagentsets.Port)
	// fixstatusleaderstr := "/consulagent/" + appname + "/" + env + "/:8500/v1/status/leader"       //, *regagentsets.Port)

	// logger.Print(fixhealthstr)
	// logger.Print(fixcatalogservicesstr)

	// if strings.Contains("/kv/")
	if !strings.Contains(uri, "/kv/") && !ragcli.IsSkip(appname, c) && strings.HasPrefix(uri, fixhealthstr) && c.Request.Method == http.MethodGet {
		service := strings.Split(strings.TrimPrefix(uri, fixhealthstr), "?")[0]
		consulapps := ragcli.GetConsulapps(appname, env, service, hostmode, c)

		c.JSON(http.StatusOK, consulapps)

		// c.Redirect(http.StatusMovedPermanently, eutarget)
	} else if !strings.Contains(uri, "/kv/") && !ragcli.IsSkip(appname, c) && strings.HasPrefix(uri, fixcatalogservicesstr) && c.Request.Method == http.MethodGet {
		catalogapps := ragcli.GetCatalogSevice(appname, env, c)

		c.JSON(http.StatusOK, catalogapps)
	} else if !strings.Contains(uri, "/kv/") && !ragcli.IsSkip(appname, c) && strings.HasPrefix(uri, fixstatusleaderstr) && c.Request.Method == http.MethodGet {
		// catalogapps := ragcli.GetCatalogSevice(appname, env, c)

		c.JSON(http.StatusOK, "127.0.0.1:8300")
	} else {
		if hostmode {
			redirectUri := *consulsets.Consul_host + strings.ReplaceAll(uri, fixproxystr, "") // + "?token=" + *consulsets.Acltoken
			logger.Info(redirectUri)
			c.Redirect(301, redirectUri)
		} else {
			// fixstr := "/consulagent/" + appname + "/" + env + "/:8500"
			target := *consulsets.Consul_host + strings.ReplaceAll(uri, fixproxystr, "")
			logger.Print(target)
			logger.Print("redirected")
			remote, err := url.Parse(*consulsets.Consul_host)
			if err != nil {
				logger.Print(err)
			}
			proxy := httputil.NewSingleHostReverseProxy(remote)
			c.Request.Host = remote.Host
			// c.Request.RemoteAddr = "user:eureka@eureka.kube.com"
			// c.Request.RequestURI = "/eureka" + strings.ReplaceAll(uri, "/eurekaagent", "")
			// c.Request.URL.User = remote.User
			c.Request.URL.Path = strings.Split(strings.ReplaceAll(uri, fixproxystr, ""), "?")[0]
			// c.Request.URL.RawQuery = strings.ReplaceAll(c.Request.URL.RawQuery, "token=", "nekot=")
			// auth := "user:eureka"
			// basicAuth := "Bearer " + base64.StdEncoding.EncodeToString([]byte(*regagentsets.Acltoken))
			// c.Request.Header.Add("Authorization", basicAuth)
			// if c.Request.Method == http.MethodPut {
			logger.Print(c.Request)
			// }
			proxy.ServeHTTP(c.Writer, c.Request)
		}
	}

}

func Eureka8500(c *gin.Context) {

	Eurekaagent(c, "/eurekaagent/")

}

func Eureka80(c *gin.Context) {

	Eurekaagent(c, "/eureka/")

}

func Eurekaagent(c *gin.Context, leadingname string) {

	appname := c.Param("appname")
	env := c.Param("env")

	uri := c.Request.RequestURI

	logger := logagent.InstArch(c)
	logger.Info(uri)

	hostmode := leadingname == "/eureka/"
	// fixstr := fmt.Sprintf("/consulagent:%s", *regagentsets.Port)
	// if ragcli.IsSkip(appname, c) || (strings.TrimSuffix(uri, "/") != "/eurekaagent/"+appname+"/"+env+"/apps" && strings.TrimSuffix(uri, "/") != "/eurekaagent/apps/delta") || c.Request.Method != http.MethodGet {
	if ragcli.IsSkip(appname, c) || (strings.TrimSuffix(uri, "/") != leadingname+appname+"/"+env+"/apps" && strings.TrimSuffix(uri, "/") != leadingname+"apps/delta") || c.Request.Method != http.MethodGet {
		// eutarget := *regagentsets.Eu_host + strings.ReplaceAll(uri, "/eurekaagent", "")
		// logger.Print(eutarget)

		if hostmode {
			redirectUri := strings.TrimSuffix(*regagentsets.Eu_host, "/eureka") + strings.ReplaceAll(uri, leadingname+appname+"/"+env, "/eureka")
			redirectUri = "http://" + strings.Split(redirectUri, "@")[1]
			logger.Info(redirectUri)
			c.Redirect(301, redirectUri)
		} else {
			logger.Print("redirected")
			remote, err := url.Parse(strings.TrimSuffix(*regagentsets.Eu_host, "/eureka"))
			if err != nil {
				logger.Print(err)
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
			// c.Request.URL.Path = strings.Split(strings.ReplaceAll(uri, "/eurekaagent/"+appname+"/"+env, "/eureka"), "?")[0]
			c.Request.URL.Path = strings.Split(strings.ReplaceAll(uri, leadingname+appname+"/"+env, "/eureka"), "?")[0]
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
		}
		// c.Redirect(http.StatusMovedPermanently, eutarget)
	} else {
		euapps := ragcli.GetEuapps(appname, env, hostmode, c)

		c.JSON(http.StatusOK, euapps)
	}
}
