package ragcli

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/max-gui/consulagent/pkg/consulhelp"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/regagent/pkg/regagentsets"
	// "github.com/max-gui/regagent/internal/pkg/regagentsets"
)

type getFromReg func(servicesmap map[string]interface{})

func getAppsList(appname, service string, getReg getFromReg, c context.Context) (servicesmap map[string]interface{}) {
	servicesmap = map[string]interface{}{}
	logger := logagent.Inst(c)
	defer func() {
		if e := recover(); e != nil {

			getReg(servicesmap)
			// euapps := eurekaapp(c)
			// for _, v := range euapps.Applications.Application {

			// 	servicesmap[strings.ToLower(v.Name)] = nil
			// }
			// servicesmap[strings.ToLower(service)] = nil
			logger := logagent.Inst(c)
			logger.WithField("appmisservice", appname).
				WithField("misservice", service).
				Info("miss iac or service field")
		}
	}()

	bytes := consulhelp.GetConfigFull(*regagentsets.ConfArchPrefix+appname, c)
	var resmap = map[string]interface{}{}
	json.Unmarshal(bytes, &resmap)

	log.Print(resmap)

	if val, ok := resmap["Application"]; ok {
		if val != nil {
			if servicesval, ok := val.(map[string]interface{})["Service"]; ok {
				log.Print(servicesval)
				if services, ok := servicesval.([]interface{}); ok {
					// appyml := resmap["Deploy"].(map[string]interface{})["Build"].(map[string]interface{})["Appyml"].(string)
					for _, v := range services {
						servicesmap[strings.ToLower(v.(string))] = nil
					}
					// return servicesmap //, appyml
				}
			}
			if allpathval, ok := val.(map[string]interface{})["Allpath"]; ok {
				log.Print(allpathval)
				if allpath, ok := allpathval.(bool); ok && allpath {
					// if allpath {
					getReg(servicesmap)
					// }
				}
			}
		}
	}

	if len(servicesmap) <= 0 {
		logger.WithField("misservice", appname).Panic("service field is empty")
	}

	return servicesmap //, ""
}

func IsSkip(appname string, c context.Context) bool {
	bytes := consulhelp.GetConfigFull(*regagentsets.ConfSkipAppPrefix+"skiplist", c)
	var resmap = map[string]interface{}{}
	json.Unmarshal(bytes, &resmap)

	_, ok := resmap[appname]

	return ok

}
