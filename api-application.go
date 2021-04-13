// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

//Package tingyun3 听云性能采集探针(sdk)
package tingyun3

//面向api使用者的接口实现部分

import (
	"os"

	"git.codemonky.net/TingYunGo/goagent/utils/logger"
)

//AppInit : 初始化听云探针
//参数:
//    jsonFile: 听云配置文件路径，文件格式为json格式
func tingyunAppInit(jsonFile string) error {
	if app == nil {
		app = new(application)
	}
	_, err := app.init(jsonFile)
	if err != nil {
		app = nil
	}
	return err
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

func init() {
	//取配置文件
	//
	configFile := os.Getenv("TINGYUN_GO_APP_CONFIG")
	if len(configFile) == 0 {
		configFile = "/etc/tingyun/go_app_config.json"
	}
	tingyunAppInit(configFile)
}
