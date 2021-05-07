// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/TingYunGo/goagent/libs/configfile"
	"github.com/TingYunGo/goagent/utils/cache_config"
	log "github.com/TingYunGo/goagent/utils/logger"
	"github.com/TingYunGo/goagent/utils/service"
)

const (
	configLocalStringNbsHost             = 1
	configLocalStringNbsLicenseKey       = 2
	configLocalStringNbsAppName          = 3
	configLocalStringNbsLevel            = log.ConfigStringNBSLevel
	configLocalStringNbsLogFileName      = log.ConfigStringNBSLogFileName
	configLocalStringMax                 = 8
	configLocalBoolAgentEnable           = 1
	configLocalBoolSSL                   = 2
	configLocalBoolAudit                 = log.ConfigBoolNBSAudit
	configLocalBoolMax                   = 8
	configLocalIntegerNbsPort            = 1
	configLocalIntegerNbsSaveCount       = 2
	configLocalIntegerNbsMaxLogSize      = log.ConfigIntegerNBSMaxLogSize
	configLocalIntegerNbsMaxLogCount     = log.ConfigIntegerNBSMaxLogCount
	configLocalIntegerNbsActionCacheMax  = 5
	configLocalIntegerNbsActionReportMax = 6
	configLocalIntegerMax                = 8

	configServerStringAppSessionKey     = 1
	configServerStringTingyunIDSecret   = 2
	configServerStringApplicationID     = 3
	configServerStringMax               = 8
	configServerBoolEnabled             = 1
	configServerBoolMax                 = 8
	configServerIntegerDataSentInterval = 1
	configServerIntegerApdexT           = 2
	configServerIntegerMax              = 8

	configServerConfigStringActionTracerRecordSQL            = 1
	configServerConfigStringRumScript                        = 2
	configServerConfigStringExternalURLParamsCaptured        = 3
	configServerConfigStringWebActionURIParamsCaptured       = 4
	configServerConfigStringInstrumentationCustom            = 5
	configServerConfigStringQuantile                         = 6
	configServerConfigStringErrorCollectorIgnoredStatusCodes = 7
	configServerConfigStringDataItemRules                    = 8
	configServerConfigStringMax                              = 16

	configServerConfigBoolAgentEnabled                   = 1
	configServerConfigBoolAutoActionNaming               = 2
	configServerConfigBoolCaptureParams                  = 3
	configServerConfigBoolErrorCollectorEnabled          = 4
	configServerConfigBoolErrorCollectorRecordDBErrors   = 5
	configServerConfigBoolActionTracerEnabled            = 6
	configServerConfigBoolActionTracerSlowSQL            = 7
	configServerConfigBoolActionTracerExplainEnabled     = 8
	configServerConfigBoolTransactionTracerEnabled       = 9
	configServerConfigBoolActionTracerNbsua              = 10
	configServerConfigBoolRumEnabled                     = 11
	configServerConfigBoolIgnoreStaticResources          = 12
	configServerConfigBoolActionTracerRemoveTrailingPath = 13
	configServerConfigBoolHotspotEnabled                 = 14
	configServerConfigBoolRumMixEnabled                  = 15
	configServerConfigBoolTransactionTracerThrift        = 16
	configServerConfigBoolMQEnabled                      = 17
	configServerConfigBoolResourceEnabled                = 18
	configServerConfigBoolLogTracking                    = 19
	configServerConfigBoolActionTracerStackTraceEnabled  = 20
	configServerConfigBoolMax                            = 24

	configServerConfigIntegerActionTracerActionThreshold     = 1
	configServerConfigIntegerActionTracerSlowSQLThreshold    = 2
	configServerConfigIntegerActionTracerExplainThreshold    = 3
	configServerConfigIntegerActionTracerStacktraceThreshold = 4
	configServerConfigIntegerRumSampleRatio                  = 5
	configServerConfigIntegerResourceLow                     = 6
	configServerConfigIntegerResourceHigh                    = 7
	configServerConfigIntegerResourceSafe                    = 8
	configServerConfigIntegerApdexThreshold                  = 9
	configServerConfigIntegerMTime                           = 10
	configServerConfigIntegerMax                             = 11

	configServerConfigIArrayIgnoredStatusCodes = 1
	configServerConfigIArrayMax                = 8
)

var localStringKeyMap = map[string]int{
	"nbs.host":          configLocalStringNbsHost,
	"collectors":        configLocalStringNbsHost,
	"nbs.license_key":   configLocalStringNbsLicenseKey,
	"license_key":       configLocalStringNbsLicenseKey,
	"nbs.app_name":      configLocalStringNbsAppName,
	"nbs.level":         configLocalStringNbsLevel,
	"nbs.log_file_name": configLocalStringNbsLogFileName,
}
var localBoolKeyMap = map[string]int{
	"nbs.agent_enabled": configLocalBoolAgentEnable,
	"nbs.ssl":           configLocalBoolSSL,
	"nbs.audit":         configLocalBoolAudit,
}

var localIntegerKeyMap = map[string]int{
	"nbs.port":              configLocalIntegerNbsPort,
	"nbs.savecount":         configLocalIntegerNbsSaveCount,
	"nbs.max_log_size":      configLocalIntegerNbsMaxLogSize,
	"nbs.max_log_count":     configLocalIntegerNbsMaxLogCount,
	"nbs.action_cache_max":  configLocalIntegerNbsActionCacheMax,
	"nbs.action_report_max": configLocalIntegerNbsActionReportMax,
}

var serverStringKeyMap = map[string]int{
	"sessionKey":      configServerStringAppSessionKey,
	"tingyunIdSecret": configServerStringTingyunIDSecret,
	"idSecret":        configServerStringTingyunIDSecret,
	"applicationId":   configServerStringApplicationID,
}

var serverBoolKeyMap = map[string]int{
	"enabled": configServerBoolEnabled,
}

var serverIntegerKeyMap = map[string]int{
	"dataSentInterval": configServerIntegerDataSentInterval,
	"apdex_t":          configServerIntegerApdexT,
}

var serverConfigStringKeyMap = map[string]int{
	"nbs.action_tracer.record_sql":         configServerConfigStringActionTracerRecordSQL,
	"nbs.rum.script":                       configServerConfigStringRumScript,
	"nbs.external_url_params_captured":     configServerConfigStringExternalURLParamsCaptured,
	"nbs.web_action_uri_params_captured":   configServerConfigStringWebActionURIParamsCaptured,
	"nbs.instrumentation_custom":           configServerConfigStringInstrumentationCustom,
	"nbs.quantile":                         configServerConfigStringQuantile,
	"error_collector.ignored_status_codes": configServerConfigStringErrorCollectorIgnoredStatusCodes,
	"data_item.rules":                      configServerConfigStringDataItemRules,
}

var serverConfigBoolKeyMap = map[string]int{
	"agent_enabled":                          configServerConfigBoolAgentEnabled,
	"auto_action_naming":                     configServerConfigBoolAutoActionNaming,
	"action_tracer.capture_params":           configServerConfigBoolCaptureParams,
	"error_collector.enabled":                configServerConfigBoolErrorCollectorEnabled,
	"nbs.error_collector.record_db_errors":   configServerConfigBoolErrorCollectorRecordDBErrors,
	"action_tracer.enabled":                  configServerConfigBoolActionTracerEnabled,
	"action_tracer.stack_enabled":            configServerConfigBoolActionTracerStackTraceEnabled,
	"nbs.action_tracer.slow_sql":             configServerConfigBoolActionTracerSlowSQL,
	"nbs.action_tracer.explain_enabled":      configServerConfigBoolActionTracerExplainEnabled,
	"nbs.action_tracer.nbsua":                configServerConfigBoolActionTracerNbsua,
	"nbs.rum.enabled":                        configServerConfigBoolRumEnabled,
	"nbs.ignore_static_resources":            configServerConfigBoolIgnoreStaticResources,
	"nbs.action_tracer.remove_trailing_path": configServerConfigBoolActionTracerRemoveTrailingPath,
	"nbs.hotspot.enabled":                    configServerConfigBoolHotspotEnabled,
	"nbs.rum.mix_enabled":                    configServerConfigBoolRumMixEnabled,
	"nbs.transaction_tracer.thrift":          configServerConfigBoolTransactionTracerThrift,
	"mq.enabled":                             configServerConfigBoolMQEnabled,
	"nbs.resource.enabled":                   configServerConfigBoolResourceEnabled,
	"nbs.log_tracking":                       configServerConfigBoolLogTracking,
}

var serverConfigIntegerKeyMap = map[string]int{

	"mTime":                                configServerConfigIntegerMTime,
	"apdex.threshold":                      configServerConfigIntegerApdexThreshold,
	"action_tracer.action_threshold":       configServerConfigIntegerActionTracerActionThreshold,
	"nbs.action_tracer.slow_sql_threshold": configServerConfigIntegerActionTracerSlowSQLThreshold,
	"nbs.action_tracer.explain_threshold":  configServerConfigIntegerActionTracerExplainThreshold,
	"action_tracer.stack_trace_threshold":  configServerConfigIntegerActionTracerStacktraceThreshold,
	"nbs.rum.sample_ratio":                 configServerConfigIntegerRumSampleRatio,
	"nbs.resource.low":                     configServerConfigIntegerResourceLow,
	"nbs.resource.high":                    configServerConfigIntegerResourceHigh,
	"nbs.resource.safe":                    configServerConfigIntegerResourceSafe,
}

type configKeyMaps struct {
	strings  map[string]int
	bools    map[string]int
	integers map[string]int
}

var localKeyMaps = configKeyMaps{localStringKeyMap, localBoolKeyMap, localIntegerKeyMap}
var serverKeyMaps = configKeyMaps{serverStringKeyMap, serverBoolKeyMap, serverIntegerKeyMap}

type dataItemRules struct {
	postParam      []string
	requestHeader  []string
	responseHeader []string
}
type dataItemRulesConfig struct {
	dataItemRules [4]dataItemRules
	currentIndex  int
}

func (d *dataItemRulesConfig) Update(config string) error {
	jsonData := make(map[string]interface{})
	if err := json.Unmarshal([]byte(config), &jsonData); err != nil {
		return err
	}
	newIndex := (d.currentIndex + 1) % 4
	items := &d.dataItemRules[newIndex]
	items.requestHeader = []string{}
	items.responseHeader = []string{}
	items.postParam = []string{}

	parseItems := func(name string) []string {
		r := []string{}
		if array, err := jsonReadArray(jsonData, name); err == nil {
			for _, v := range array {
				if item, ok := v.(string); ok {
					r = append(r, item)
				}
			}
		}
		return r
	}
	items.postParam = parseItems("postParam")
	items.requestHeader = parseItems("requestHeader")
	items.responseHeader = parseItems("responseHeader")
	return nil
}
func (d *dataItemRulesConfig) Get() *dataItemRules {
	return &d.dataItemRules[d.currentIndex]
}
func (d *dataItemRulesConfig) Commit() {
	d.currentIndex = (d.currentIndex + 1) % 4
}

type configurations struct {
	local         cache_config.Configuration
	server        cache_config.Configuration
	serverExt     cache_config.Configuration
	serverArrays  cache_config.Arrays
	svc           service.Service
	apdexs        apdexActionMap
	dataItemRules dataItemRulesConfig
	loginCount    int64
	configfile    string
	appID         int64
	lock          sync.RWMutex
	started       bool
	loginError    bool
	reported      bool
}

func parseConfig(filenames string, c *cache_config.Configuration) error {
	files := strings.Split(filenames, ":")
	nameFound := false
	for _, filename := range files {
		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		extName := fileExtName(filename)
		if extName == "json" {
			jsonData := map[string]interface{}{}
			if err = json.Unmarshal(bytes, &jsonData); err != nil {
				return err
			}
			for k, v := range jsonData {
				if k == "nbs.app_name" {
					nameFound = true
				}
				c.Update(localStringKeyMap, localBoolKeyMap, localIntegerKeyMap, k, v)
			}
		} else if extName == "conf" {

			conf := configfile.Parse(bytes)
			if conf == nil {
				return errors.New("conf parse failed")
			}
			session := conf.Session("")
			if session != nil {
				for k, v := range session {
					if !v.IsArray() {
						value := v.Get()
						if isBool(value) {
							bvalue, _ := toBool(value)
							c.Update(localStringKeyMap, localBoolKeyMap, localIntegerKeyMap, k, bvalue)
						} else if isNumber(value) {
							nvalue, _ := strconv.Atoi(value)
							c.Update(localStringKeyMap, localBoolKeyMap, localIntegerKeyMap, k, nvalue)
						} else {
							if k == "nbs.app_name" {
								nameFound = true
							} else if k == "collectors" {
								collectors := strings.Split(value, ",")
								if len(collectors) > 0 {
									value = collectors[0]
								}
							}
							c.Update(localStringKeyMap, localBoolKeyMap, localIntegerKeyMap, k, value)
						}
					}
				}
			}
		}
	}
	if !nameFound {
		c.Update(localStringKeyMap, localBoolKeyMap, localIntegerKeyMap, "nbs.app_name", getDefaultAppName())
	}
	c.Commit()
	return nil
}
func (c *configurations) update(configfile string) error {
	if err := parseConfig(c.configfile, &c.local); err != nil {
		return err
	}
	c.configfile = configfile
	return nil
}
func (c *configurations) Init(configfile string) error {
	if len(c.configfile) > 0 {
		return c.update(configfile)
	}
	c.local.Init(configLocalStringMax, configLocalBoolMax, configLocalIntegerMax)
	c.server.Init(configServerStringMax, configServerBoolMax, configServerIntegerMax)
	c.serverExt.Init(configServerConfigStringMax, configServerConfigBoolMax, configServerConfigIntegerMax)
	c.serverArrays.Init(configServerConfigIArrayMax)
	c.apdexs.Init()
	err := parseConfig(configfile, &c.local)
	c.started = err == nil
	c.reported = false
	c.loginCount = 0
	c.loginError = false
	if c.started {
		c.configfile = configfile
		c.svc.Start(func(running func() bool) {
			lastTime := time.Now()
			lastModify, err := modifyTime(c.configfile)
			if err != nil {
				lastModify = lastTime
			}
			for running() {
				time.Sleep(time.Second)
				if now := time.Now(); now.Sub(lastTime) < 60*time.Second {
					continue
				}
				if modTime, err := modifyTime(c.configfile); err == nil {
					if modTime.Equal(lastModify) {
						continue
					}
					if parseConfig(c.configfile, &c.local) == nil {
						lastModify = modTime
					}
				}
			}

		})
	}
	return err
}
func (c *configurations) NeverLogin() bool { return c.loginCount == 0 }
func (c *configurations) HasLogin() bool   { return c.loginCount > 0 && !c.loginError }
func (c *configurations) Release() {
	if c.started {
		c.started = false
		c.svc.Stop()
	}
}
func (c *configurations) UpdateServerConfig(result map[string]interface{}, login bool) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	for {
		if _, err := jsonReadString(result, "sessionKey"); err != nil {
			return err
		} else if _, err = jsonToString(result, "appId"); err != nil {
			return err
		} else {
			c.appID, _ = jsonReadInt64(result, "appId")
			for k, v := range result {
				c.server.Update(serverStringKeyMap, serverBoolKeyMap, serverIntegerKeyMap, k, v)
			}

			if config, err := jsonReadObjects(result, "config"); err == nil {
				for k, v := range config {
					c.serverExt.Update(serverConfigStringKeyMap, serverConfigBoolKeyMap, serverConfigIntegerKeyMap, k, v)
				}
			}
			c.apdexs.apdexT = readServerConfigInt(configServerConfigIntegerApdexThreshold, 500)
			if actionApdex, err := jsonReadObjects(result, "actionApdex"); err == nil {
				for k, v := range actionApdex {
					if val, err := readInt(v); err == nil {
						c.apdexs.Update(k, val)
					}
				}
			}

			c.server.Commit()
			c.serverExt.Commit()
			ignoreStatus := c.serverExt.CStrings.Read(configServerConfigStringErrorCollectorIgnoredStatusCodes, "")
			if ignoreStatus != "" {
				statusArray := strings.Split(ignoreStatus, ",")
				if len(statusArray) > 0 {
					intArray := make([]int64, len(statusArray))
					count := 0
					for _, v := range statusArray {
						if b, err := strconv.Atoi(v); err == nil {
							intArray[count] = int64(b)
							count++
						}
					}
					if count > 0 {
						c.serverArrays.Update(configServerConfigIArrayIgnoredStatusCodes, intArray[0:count])
					}
				}
			}
			c.dataItemRules.Update(readServerConfigString(configServerConfigStringDataItemRules, "{}"))
			c.dataItemRules.Commit()
			c.serverArrays.Commit()
			c.apdexs.Commit()
			if login {
				c.loginCount++
			}
		}
		return nil
	}
}
func (c *configurations) UpdateConfig(onFirst func()) {
	if c.loginCount > 0 {
		if !c.reported {
			c.reported = true
			onFirst()
		}
	}
}
func configValue(config *cache_config.Configuration, key string, maps *configKeyMaps) (interface{}, bool) {
	if v, found := maps.strings[key]; found {
		return config.CStrings.Find(v)
	} else if v, found := maps.bools[key]; found {
		return config.CBools.Find(v)
	} else if v, found := maps.integers[key]; found {
		return config.CIntegers.Find(v)
	}
	return nil, false
}
func (c *configurations) Value(name string) (interface{}, bool) {
	if v, found := configValue(&c.server, name, &serverKeyMaps); found {
		return v, found
	}
	return configValue(&c.local, name, &localKeyMaps)
}

func modifyTime(filename string) (time.Time, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

type apdexActionMap struct {
	current int
	apdexT  int
	arrays  [4]map[string]int
}

func (s *apdexActionMap) Init() *apdexActionMap {
	s.current = 3
	s.apdexT = 500
	for i := 0; i < 4; i++ {
		s.arrays[i] = make(map[string]int)
	}
	return s
}
func (s *apdexActionMap) Read(key string) int {
	if r, ok := s.arrays[s.current][key]; ok {
		return int(r)
	}
	return s.apdexT
}
func (s *apdexActionMap) Update(key string, value int) {
	s.arrays[(s.current+1)%4][key] = value
}
func (s *apdexActionMap) Commit() {
	s.current = (s.current + 1) % 4
}
