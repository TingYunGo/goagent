// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/TingYunGo/goagent/libs/pool"

	log "github.com/TingYunGo/goagent/utils/logger"
)

/*
 * 日志级别/审计模式控制
 */
const (
	LevelOff      = log.LevelOff
	LevelCritical = log.LevelCritical
	LevelError    = log.LevelError
	LevelWarning  = log.LevelWarning
	LevelInfo     = log.LevelInfo
	LevelVerbos   = log.LevelVerbos
	LevelDebug    = log.LevelDebug
	LevelMask     = log.LevelMask
	Audit         = log.Audit
)

//actionPool里的请求归纳处理
func (a *application) parseActions(parseMax int) int {
	actionParsed := 0
	for a.actionPool.Size() > 0 && actionParsed < parseMax {
		if action := a.actionPool.Get(); action != nil {
			actionParsed++
			if a.configs.HasLogin() {
				a.GetReportBlock().Append(action.(*Action))
			}
			action.(*Action).destroy()
		} else {
			break
		}
	}
	return actionParsed
}

const (
	serverUnInited     = 0
	serverInLogin      = 1
	serverLoginSuccess = 2
	serverLoginFaild   = 3
	uploadSuccess      = 0
	uploadError        = 1
)

type serverControl struct {
	loginResultTime   time.Time //前一次login返回时间(无论成功失败)
	lastUploadTime    time.Time //前一次上传时间，失败重传需要等待几秒钟间隔，据此计算
	reportTime        time.Time
	finishedPool      pool.SerialReadPool "完成任务通知队列"
	uploadCount       int32
	loginState        uint8
	postIsReturn      bool
	requestLoginReset bool
}

func (s *serverControl) Pushback(eventCallback func()) {
	s.finishedPool.Put(eventCallback)
}
func (s *serverControl) doEvent() {
	for t := s.finishedPool.Get(); t != nil; t = s.finishedPool.Get() {
		t.(func())()
	}
}
func (s *serverControl) Reset() {
	if s.requestLoginReset {
		return
	} else if s.loginState == serverInLogin {
		s.requestLoginReset = true
	} else if s.loginState == serverLoginSuccess {
		s.loginState = serverUnInited
	}
}
func (s *serverControl) OnReturn() {
	s.postIsReturn = true
	s.lastUploadTime = time.Now()
}

func (s *serverControl) init() {
	s.loginState = serverUnInited
	s.postIsReturn = false
	s.finishedPool.Init()
}

//redirect, login, 处理
func (a *application) checkLogin() bool {
	//需要添加如下变量

	//appId
	//report计时器

	//处理过程

	//若login状态为serverLoginSuccess,则返回true
	//若login状态为serverUnInited,则开始login过程,置状态为serverInLogin, 返回false
	//若login状态为serverInLogin,则返回false
	//若状态为serverLoginFaild, 5 second waiting,则等待时间=time.Now()-loginResultTime
	//   若等待时间小于5秒,则返回false
	//   若等待时间大于等于5秒,则置login状态为serverUnInited,返回false

	//login完成回调(另一个routine)
	//    1.写日志
	//    2.更新loginResultTime
	//    2.成功,设置application Id, 设置配置项参数,设置login 状态为serverLoginSuccess,置位clear标志
	//    3.失败,设置login状态为 5 second waiting
	if a.serverCtrl.requestLoginReset && a.serverCtrl.loginState != serverInLogin {
		a.serverCtrl.loginState = serverUnInited
		a.serverCtrl.requestLoginReset = false
	}
	switch a.serverCtrl.loginState {
	case serverUnInited:
		a.startLogin()
		return false
	case serverInLogin:
		return false
	case serverLoginFaild:
		a.checkFaildState()
		return false
	case serverLoginSuccess:
		return true
	default:
		//未定义行为, 不该出现 log error
		return false
	}
}

func (a *application) ReadApdex(name string) int32 {
	return int32(a.configs.apdexs.Read(name))
}

func (a *application) setNextReportTime(value time.Time, reportInterval int) {
	unix := value.Unix()
	unix = int64(reportInterval) + unix - unix%int64(reportInterval)
	a.serverCtrl.reportTime = time.Unix(unix, 0)
}

//login返回, 验证application Id, session key
func (a *application) parseLoginResult(jsonData map[string]interface{}) error {
	//验证status
	//如果status不为success,则从result生成错误信息,返回error
	for {
		if status, err := jsonReadString(jsonData, "status"); err != nil {
			break
		} else if status != "success" {
			break
		} else if result, err := jsonReadObjects(jsonData, "result"); err != nil {
			break
		} else if err = a.configs.UpdateServerConfig(result, true); err != nil {
			return err
		}
		enabled := readServerConfigBool(configServerConfigBoolAgentEnabled, true)
		if enabled {
			enabled = a.configs.serverExt.CBools.Read(configServerConfigBoolAgentEnabled, true)
		}
		if enabled {
			now := time.Now()
			if a.configs.NeverLogin() {
				a.logger.Println(LevelInfo, "ApplicationStart...")
				a.setNextReportTime(now, int(a.configs.server.CIntegers.Read(configServerIntegerDataSentInterval, 60)))
			}
		}

		return nil
	}
	b, _ := json.Marshal(jsonData)
	return errors.New("server result json error : " + string(b))
}
func (a *application) startLogin() {
	a.serverCtrl.loginState = serverInLogin
	err := a.server.Login(func(err error, result map[string]interface{}) {
		defer a.serverCtrl.OnReturn()
		a.serverCtrl.loginResultTime = time.Now()
		if err != nil { //login过程有错误,置状态,写日志
			a.logger.Println(log.LevelError, err)
			a.serverCtrl.loginState = serverLoginFaild
		} else if err := a.parseLoginResult(result); err != nil { //login结果有错误,写日志,置状态
			a.logger.Println(log.LevelError, err)
			a.serverCtrl.loginState = serverLoginFaild
		} else { //login成功
			if !a.configs.HasLogin() {
				a.serverCtrl.lastUploadTime = time.Now()
			}
			a.server.lastAlive = time.Now()
			a.serverCtrl.loginState = serverLoginSuccess
		}
	})
	if err != nil {
		a.serverCtrl.loginResultTime = time.Now()
		a.serverCtrl.loginState = serverLoginFaild
	}
}
func (a *application) checkFaildState() {
	if now := time.Now(); now.Sub(a.serverCtrl.loginResultTime) < 5*time.Second {
		return
	}
	a.serverCtrl.loginState = serverUnInited
}

func (a *application) timerCheck() {
	if a.serverCtrl.postIsReturn {
		a.server.ReleaseRequest()
		a.serverCtrl.postIsReturn = false
	}
	//不在LoginSuccess状态,返回
	if !a.checkLogin() {
		return
	}
	a.configs.UpdateConfig(a.Runtime.Init)
	a.upload()
	a.server.keepAlive(func(err error, jsonData map[string]interface{}) {
		if err != nil {
			a.logger.Println(LevelError, "keepAlive error:", err.Error(), ": Reset Login")
			a.serverCtrl.Reset()
			return
		}
		if status, err := jsonReadString(jsonData, "status"); err != nil || status != "success" {
			a.logger.Println(LevelError, "keepAlive result status not success: ", status)
			a.serverCtrl.Reset()
			return
		}
		if result, err := jsonReadObjects(jsonData, "result"); err == nil {
			if id, _ := jsonReadString(result, "id"); id == "ConfigChanged" {
				if args, err := jsonReadObjects(result, "args"); err == nil {
					if sessionKey := a.configs.server.CStrings.Read(configServerStringAppSessionKey, ""); len(sessionKey) > 0 {
						args["sessionKey"] = sessionKey
						args["appId"] = a.configs.appID
						a.configs.UpdateServerConfig(args, false)
					}
				}
			}
		}
	})
}

var dropTrace int = 0

func (a *application) upload() {
	if a.server.request != nil {
		return
	}
	if a.reportQueue.Size() > 0 {
		if time.Now().Sub(a.serverCtrl.lastUploadTime) < time.Second && a.reportQueue.Size() == 1 {
			return
		}
		saveCount := a.configs.local.CIntegers.Read(configLocalIntegerNbsSaveCount, 10)
		if saveCount < 1 {
			saveCount = 1
		}
		for a.reportQueue.Size() > int(saveCount) {
			data, _ := a.reportQueue.PopFront()
			t := data.(*structAppData)
			if len(t.traces.Traces) > 0 {
				dropTrace += len(t.traces.Traces)
			}
			t.destroy()
		}

		data, _ := a.reportQueue.PopFront()
		traceData := data.(*structAppData)
		if len(traceData.traces.Traces) == 0 {
			return
		}
		b, _ := traceData.Serialize()
		traceData.destroy()
		var err error
		a.server.request, err = a.server.Upload(b, func(err error, rcode int, statusCode int) {
			if rcode == -2 {
				a.logger.Println(LevelError, "Upload error: status:", statusCode, ",rcode: ", rcode, ": ", err)
			}
			a.serverCtrl.Pushback(func() {
				a.serverCtrl.OnReturn()
			})
		})
		if err != nil {
			a.logger.Println(LevelError, "App.", "Upload Error :", err)
			a.serverCtrl.lastUploadTime = time.Now()
		}
	}
}

func (a *application) loop(running func() bool) {
	init_delay := app.configs.local.CIntegers.Read(configLocalIntegerAgentInitDelay, 1)
	time.Sleep(time.Second * time.Duration(init_delay))
	lastParsed := 1
	sleepDuration := time.Millisecond
	for running() {
		a.serverCtrl.doEvent()
		//处理采集到的事务
		actionParsed := a.parseActions(10000)

		//发送到server
		a.timerCheck()
		if actionParsed == 0 {
			if lastParsed == 0 && sleepDuration < 100*time.Millisecond {
				sleepDuration *= 2
			}
			time.Sleep(sleepDuration)
		} else {
			sleepDuration = time.Millisecond
		}
		lastParsed = actionParsed
	}
}

var firstRun bool = true

func makeLoginRequest() ([]byte, error) {
	//用参数初始化头信息
	//构造login请求包
	//序列化数据
	getHost := func() string {
		if host, e := os.Hostname(); e == nil {
			return host
		}
		return "Unknown"
	}
	getPath := func() string {
		file, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(file)
		return path
	}
	oneagentUUID := ""
	if os.Getenv("TINGYUN_ONEAGENT_GO") == "enable" {
		oneagentUUID = trimSpace(fileCat("/opt/tingyun-oneagent/conf/oneagent.uuid"))
	}
	app.logger.Println(LevelDebug, "Login:", httpListenAddr.Addr)
	port, _ := parsePort(httpListenAddr.Addr)
	envInfo := map[string]string{}
	for _, item := range os.Environ() {
		if k, v := parseEnv(item); len(k) > 0 && len(v) > 0 {
			envInfo[k] = v
		}
	}
	return json.Marshal(map[string]interface{}{
		"host":         getHost(),
		"appName":      app.configs.local.CStrings.Read(configLocalStringNbsAppName, "TingYunDefault"),
		"language":     "Go",
		"oneAgentUuid": oneagentUUID,
		"port":         port,
		"agentTime":    time.Now().Unix(),
		"agentVersion": TINGYUN_GO_AGENT_VERSION,
		"pid":          os.Getpid(),
		"firstRun":     firstRun,
		"environment": map[string]interface{}{
			"config": map[string]interface{}{
				"license_key":       app.configs.local.CStrings.Read(configLocalStringNbsLicenseKey, ""),
				"nbs.log_file_name": app.configs.local.CStrings.Read(configLocalStringNbsLogFileName, "agent.log"),
				"audit":             app.configs.local.CBools.Read(configLocalBoolAudit, false),
				"nbs.max_log_count": app.configs.local.CIntegers.Read(configLocalIntegerNbsMaxLogCount, 3),
				"nbs.max_log_size":  app.configs.local.CIntegers.Read(configLocalIntegerNbsMaxLogSize, 10),
				"nbs.ssl":           app.configs.local.CBools.Read(configLocalBoolSSL, false),
			},
			"env": envInfo,
			"system": map[string]string{
				"cmdline":    getPath(),
				"OS":         runtime.GOOS,
				"ARCH":       runtime.GOARCH,
				"Compiler":   runtime.Compiler,
				"Go-Version": runtime.Version(),
				"GOROOT":     runtime.GOROOT(),
			},
			"meta": map[string]interface{}{
				"readonly":      true,
				"pid":           os.Getpid(),
				"agentVersion":  TINGYUN_GO_AGENT_VERSION,
				"tingyun.debug": false,
			},
		},
	})
}
