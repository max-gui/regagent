package ragcli

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/max-gui/consulagent/pkg/consulhelp"
	"github.com/max-gui/consulagent/pkg/consulsets"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/regagent/pkg/regagentsets"
	// "github.com/max-gui/regagent/internal/pkg/regagentsets"
)

type mutexKV struct {
	sync.RWMutex
	kvs map[string]interface{}
}

var kvmap = mutexKV{kvs: make(map[string]interface{})}

func (v *mutexKV) help(tricky func(map[string]interface{}) (bool, interface{})) (bool, interface{}) {
	v.Lock()
	ok, res := tricky(v.kvs)
	v.Unlock()
	return ok, res
}

func Getall(c context.Context) map[string]interface{} {
	if ok, value := kvmap.help(func(kvs map[string]interface{}) (bool, interface{}) {
		if val, ok := kvs["allservices"]; ok {
			return ok, val
		} else {
			return ok, nil
		}
	}); ok {
		realvalue := value.(struct {
			services  map[string]interface{}
			lastCheck time.Time
		})
		if time.Duration(*consulsets.Cacheminutes)*time.Minute > time.Since(realvalue.lastCheck) {
			return realvalue.services
		}
	}

	servicesmap := map[string]interface{}{}

	servicearrmap := consulhelp.GetServices(c)
	for k := range servicearrmap {
		servicesmap[k] = struct{}{}
	}

	euapps := eurekapps(c)
	for _, v := range euapps.Applications.Application {

		servicesmap[strings.ToLower(v.Name)] = nil
	}

	return servicesmap
}

// type getFromReg func(servicesmap map[string]interface{})

func GetAproveServices(appname string, c context.Context) (services map[string]interface{}, allpath bool) {
	logger := logagent.InstArch(c)

	bytes := consulhelp.GetConfigFull(*regagentsets.ConfArchPrefix+appname, c)
	var resmap = map[string]interface{}{}
	json.Unmarshal(bytes, &resmap)

	// logger.Print(resmap)

	if val, ok := resmap["Application"]; ok {
		if val != nil {
			if servicesval, ok := val.(map[string]interface{})["Service"]; ok {
				// logger.Print(servicesval)
				if services, ok := servicesval.([]interface{}); ok {
					// appyml := resmap["Deploy"].(map[string]interface{})["Build"].(map[string]interface{})["Appyml"].(string)
					servicesmap := map[string]interface{}{}
					for _, v := range services {
						servicesmap[strings.ToLower(v.(string))] = nil
					}
					return servicesmap, false
				} else if allpathval, ok := val.(map[string]interface{})["Allpath"]; ok {
					// logger.Print(allpathval)
					if allpath, ok := allpathval.(bool); ok && allpath {
						// if allpath {
						return nil, true
						// }
					} else {
						return map[string]interface{}{}, false
					}
				}
			}
		}
	}

	logger.WithField("misservice", appname).Panic("service field is empty")
	return nil, false
}

func getAppsList(appname, service string, c context.Context) map[string]interface{} {
	// logger := logagent.Inst(c)
	// servicesmap, isallpath := GetAproveServices(appname, c)
	if servicesmap, isallpath := GetAproveServices(appname, c); isallpath {
		return Getall(c)
	} else {
		servicesmap["defualt"] = nil
		return servicesmap
	}
}

// func getAppsListOld(appname, service string, getReg getFromReg, c context.Context) (servicesmap map[string]interface{}) {
// 	servicesmap = map[string]interface{}{}
// 	logger := logagent.InstArch(c)
// 	defer func() {
// 		if e := recover(); e != nil {

// 			getReg(servicesmap)
// 			// euapps := eurekaapp(c)
// 			// for _, v := range euapps.Applications.Application {

// 			// 	servicesmap[strings.ToLower(v.Name)] = nil
// 			// }
// 			// servicesmap[strings.ToLower(service)] = nil
// 			logger := logagent.InstArch(c)
// 			logger.WithField("appmisservice", appname).
// 				WithField("misservice", service).
// 				Info("miss iac or service field")
// 		}
// 	}()

// 	bytes := consulhelp.GetConfigFull(*regagentsets.ConfArchPrefix+appname, c)
// 	var resmap = map[string]interface{}{}
// 	json.Unmarshal(bytes, &resmap)

// 	// logger.Print(resmap)

// 	if val, ok := resmap["Application"]; ok {
// 		if val != nil {
// 			if servicesval, ok := val.(map[string]interface{})["Service"]; ok {
// 				// logger.Print(servicesval)
// 				if services, ok := servicesval.([]interface{}); ok {
// 					// appyml := resmap["Deploy"].(map[string]interface{})["Build"].(map[string]interface{})["Appyml"].(string)
// 					for _, v := range services {
// 						servicesmap[strings.ToLower(v.(string))] = nil
// 					}
// 					// return servicesmap //, appyml
// 				}
// 			}
// 			if allpathval, ok := val.(map[string]interface{})["Allpath"]; ok {
// 				// logger.Print(allpathval)
// 				if allpath, ok := allpathval.(bool); ok && allpath {
// 					// if allpath {
// 					getReg(servicesmap)
// 					// }
// 				}
// 			}
// 		}
// 	}

// 	if len(servicesmap) <= 0 {
// 		logger.WithField("misservice", appname).Panic("service field is empty")
// 	}

// 	return servicesmap //, ""
// }

func IsSkip(appname string, c context.Context) bool {
	bytes := consulhelp.GetConfigFull(*regagentsets.ConfSkipAppPrefix+"skiplist", c)
	var resmap = map[string]interface{}{}
	json.Unmarshal(bytes, &resmap)

	_, ok := resmap[appname]

	return ok

}
