// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"errors"
	"time"

	"github.com/TingYunGo/goagent/utils/logger"

	"github.com/TingYunGo/goagent/libs/list"
	"github.com/TingYunGo/goagent/libs/pool"
	"github.com/TingYunGo/goagent/utils/service"
)

func (a *application) stop() {
	if a == nil {
		return
	}
	a.svc.Stop()
	a.logger.Printf(LevelInfo, "Agent stoped\n")
	a.logger.Release()
	a.configs.Release()
}

//action结束,将Action事务对象抛给app处理
func (a *application) appendAction(action *Action) {
	action.time.End()
	a.actionPool.Put(action)
}

//给Action.Finish调用
func appendAction(action *Action) {
	if app != nil {
		app.appendAction(action)
	} else {
		//释放Action对象
		action.destroy()
	}
}
func readServerConfigInt(id int, defaultValue int) int {
	if app == nil {
		return defaultValue
	}
	return int(app.configs.serverExt.CIntegers.Read(id, int64(defaultValue)))
}
func readServerConfigBool(id int, defaultValue bool) bool {
	if app == nil {
		return defaultValue
	}
	return app.configs.serverExt.CBools.Read(id, defaultValue)
}
func readLocalConfigBool(id int, defaultValue bool) bool {
	if app == nil {
		return defaultValue
	}
	return app.configs.local.CBools.Read(id, defaultValue)
}

func readServerConfigString(id int, defaultValue string) string {
	if app == nil {
		return defaultValue
	}
	return app.configs.serverExt.CStrings.Read(id, defaultValue)
}
func (a *application) init(configfile string) (*application, error) {
	err := a.configs.Init(configfile)
	if err != nil {
		return nil, err
	}
	if enabled := a.configs.local.CBools.Read(configLocalBoolAgentEnable, true); !enabled {
		configDisabled = true
		a.configs.Release()
		return nil, errors.New("Agent Is disabled by config file")
	}
	if appname := a.configs.local.CStrings.Read(configLocalStringNbsAppName, ""); appname == "" {
		return nil, errors.New(configfile + ": nbs.app_name not found")
	}

	if license := a.configs.local.CStrings.Read(configLocalStringNbsLicenseKey, ""); license == "" {
		return nil, errors.New(configfile + ": nbs.license_key not found")
	}
	a.serverCtrl.Reset()
	if a.logger == nil {
		a.logger = log.New(&a.configs.local)
		a.actionPool.Init()
		a.server.init(&a.configs)
		a.serverCtrl.init()
		a.reportQueue.Init()
		a.Runtime.Init()
		a.svc.Start(a.loop)
	}
	a.logger.Println(log.LevelInfo, "App Init by ", configfile)
	return a, nil
}
func (a *application) createAction(category string, method string, istask bool) (*Action, error) {
	if enabled := readServerConfigBool(configServerConfigBoolAgentEnabled, true); !enabled {
		return nil, errors.New("Agent disabled by server config")
	}
	action := &Action{
		category:       category,
		url:            "",
		trackID:        "",
		trackEnable:    false,
		statusCode:     0,
		requestParams:  make(map[string]string),
		responseParams: make(map[string]string),
		customParams:   make(map[string]string),
		stateUsed:      actionUsing,
		tracerIDMaker:  0,
		root: &Component{
			tracerParentID: -1,
			exID:           false,
			callStack:      nil,
			time:           timeRange{time.Now(), -1},
			_type:          ComponentDefault,
		},
		isTask: istask,
	}
	action.current = action.root
	if category == "URI" {
		action.path = method
	} else {
		action.method = method
		action.root.method = method
	}
	action.root.action = action
	action.root.tracerID = action.makeTracerID()
	action.cache.Init()
	action.errors.Init()
	action.time.Init()
	action.cache.Put(action.root)

	return action, nil
}

type application struct {
	configs     configurations      "配置选项集合"
	actionPool  pool.SerialReadPool "完成事务消息池"
	logger      *log.Logger
	svc         service.Service
	server      serviceDC
	serverCtrl  serverControl
	reportQueue list.List
	Runtime     runtimePerf
}
