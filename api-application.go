// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

//Package tingyun3 听云性能采集探针(sdk)
package tingyun3

//面向api使用者的接口实现部分

import (
	"fmt"
	"os"

	"github.com/TingYunGo/goagent/libs/tystring"

	"github.com/TingYunGo/goagent/utils/logger"
)

//AppInit : 初始化听云探针
//参数:
//    jsonFile: 听云配置文件路径，文件格式为json格式
func tingyunAppInit(jsonFile string) error {
	if app == nil {
		new_app, err := new(application).init(jsonFile)
		if new_app == nil {
			return err
		}
		app = new_app
	}
	return nil
}

//Running : 检测探针是否启动(为Frameworks提供接口)
//返回值: bool
func tingyunRunning() bool {
	return app != nil
}

//AppStop : 停止听云探针
func tingyunAppStop() {
	if app == nil {
		return
	}
	app.stop()
	app = nil
}

// ConfigRead : 读配置项
func ConfigRead(name string) (interface{}, bool) {
	if app == nil {
		return nil, false
	}
	return app.configs.Value(name)
}

// Log : 返回日志对象接口
func Log() *log.Logger {
	if app == nil {
		return nil
	}
	return app.logger
}

var configDisabled bool = false
var app *application = nil

var defaultAppName = "GoApp"

func getDefaultAppName() string {
	return defaultAppName
}
func checkOneagent() bool {
	if tystring.CaseCMP(os.Getenv("TINGYUN_ONEAGENT_GO"), "enable") != 0 {
		return false
	}
	if !fileExist("/opt/tingyun-oneagent/conf/oneagent.conf") {
		return false
	}
	if !fileExist("/opt/tingyun-oneagent/conf/go.conf") {
		return false
	}
	return true
}

var isOneagent = false

func oneagentLogPath() string {
	if isOneagent {
		if containerID := getContainerID(); len(containerID) > 0 {
			return fmt.Sprintf("/opt/tingyun-oneagent/logs/agent/golang-agent-%s-%d.log", containerID, os.Getpid())
		} else {
			return fmt.Sprintf("/opt/tingyun-oneagent/logs/agent/golang-agent-%d.log", os.Getpid())
		}
	}
	return ""
}
func envGetAppName() string {
	return os.Getenv("TINGYUN_GO_APP_NAME")
}
func init() {
	listens.init()
	//check user defined
	configFile := os.Getenv("TINGYUN_GO_APP_CONFIG")
	//check oneagent defined
	if len(configFile) == 0 {
		if checkOneagent() {
			isOneagent = true
			configFile = "/opt/tingyun-oneagent/conf/oneagent.conf:/opt/tingyun-oneagent/conf/go.conf"
		}
	}
	//default
	if len(configFile) == 0 {
		configFile = "/etc/tingyun/go_app_config.json"
	}
	if appname := envGetAppName(); len(appname) > 0 {
		defaultAppName = appname
	} else {
		defaultAppName = getExt(readLink("/proc/self/exe"), '/')
	}
	tingyunAppInit(configFile)
}
