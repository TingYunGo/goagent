// +build linux
// +build amd64 arm64
// +build cgo

package redis

import (
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/TingYunGo/goagent/libs/tystring"

	"github.com/TingYunGo/goagent"

	"github.com/go-redis/redis"
)

const (
	StorageIndexRedis                   = tingyun3.StorageIndexRedis
	cmdbaseClientprocess                = 16
	cmdbaseClientprocessPipeline        = 17
	cmdbaseClientgeneralProcessPipeline = 18
	cmdClusterClient_processPipeline    = 19
	cmdClusterClient_processTxPipeline  = 20
	cmdProcessHook                      = 128
	cmdProcessPipeHook                  = 136
)

var skipTokens = []string{
	"github.com/go-redis/redis",
	"github.com/TingYunGo/goagent",
}

func readConfigInt(name string, defaultValue int) int {
	v, exist := tingyun3.ConfigRead(name)
	if !exist {
		return defaultValue
	}
	if value, ok := v.(int64); ok {
		return int(value)
	}
	return defaultValue
}

//go:noinline
func getCallName(skip int) (callerName string) {
	skip++
	callerName = tingyun3.GetCallerName(skip)
	isSkipName := func(name string) bool {
		for _, token := range skipTokens {
			if tystring.SubString(name, 0, len(token)) == token {
				return true
			}
		}
		return false
	}
	for isSkipName(callerName) {
		skip++
		callerName = tingyun3.GetCallerName(skip)
	}
	return
}

type processContext struct {
	begin time.Time
	cmd   int
	layer int
}

var objectSkipList = []string{
	"AUTH",
	"ECHO",
	"PING",
	"QUIT",
	"SELECT",
}

func handleGoRedis(host, cmd, object string, begin time.Time, err error, skip int) {
	action := tingyun3.GetAction()
	callerName := ""
	if action == nil {
		callerName = getCallName(3)
		if action, _ = tingyun3.CreateTask(callerName); action != nil {
			action.FixBegin(begin)
			defer func() {
				action.Finish()
				tingyun3.LocalClear()
			}()
		}
	}
	if action == nil {
		return
	}
	if len(callerName) == 0 {
		callerName = getCallName(3)
	}
	component := action.CreateRedisComponent(host, cmd, object, callerName)
	if component == nil {
		return
	}
	component.FixBegin(begin)
	if err != nil {
		component.SetException(err, callerName, 3)
	}
	component.FixStackEnd(skip, func(funcname string) bool {
		token := "github.com/go-redis/redis/"
		token1 := "github.com/TingYunGo/goagent/"
		return tystring.SubString(funcname, 0, len(token)) == token || tystring.SubString(funcname, 0, len(token1)) == token1
	})
}

func getArgs(args []interface{}) (cmd, object string) {
	argc := len(args)
	cmd, object = "", ""
	if argc > 0 {
		cmd = args[0].(string)
	}
	if argc > 1 {
		object = args[1].(string)
	}
	if tystring.FindString(objectSkipList, cmd) != -1 {
		object = ""
	}
	return
}
func parseCmder(cmd redis.Cmder) (string, string) {
	return getArgs(cmd.Args())
}
func parseCmders(cmds []redis.Cmder) (string, string) {
	c := "["
	o := "["
	for i, v := range cmds {
		if i > 0 {
			c = c + ","
			o = o + ","
		}
		args := v.Args()
		argc := len(args)
		cmd := ""
		obj := ""
		if argc > 0 {
			cmd = args[0].(string)
		}
		if argc > 1 {
			obj = args[1].(string)
		}
		if tystring.FindString(objectSkipList, cmd) != -1 {
			obj = ""
		}
		if len(c)+len(cmd)+len(o)+len(obj) > 180 && i > 0 {
			c = c + "..."
			o = o + "..."
			break
		}
		c = c + cmd
		o = o + obj
	}
	c = c + "]"
	o = o + "]"
	return c, o
}

type baseClient struct {
	opt  *redis.Options
	pool interface{}
}

type pipelineProcessor func(uintptr, []redis.Cmder) (bool, error)

//go:noinline
func baseClientProcess(c *baseClient, cmd redis.Cmder) error {
	trampoline.arg4 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapbaseClientProcess(c *baseClient, cmd redis.Cmder) error {
	begin := time.Now()
	req := tingyun3.LocalGet(StorageIndexRedis)
	if req == nil {
		tingyun3.LocalSet(StorageIndexRedis, 1)
	}

	err := baseClientProcess(c, cmd)
	if req == nil {
		tingyun3.LocalDelete(StorageIndexRedis)
		command, object := parseCmder(cmd)
		handleGoRedis(c.opt.Addr, command, object, begin, err, 2)
		if tingyun3.GetAction() == nil {
			tingyun3.LocalClear()
		}
	}
	return err
}

//go:noinline
func redisNewClient(opt *redis.Options) *redis.Client {
	trampoline.arg9 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapredisNewClient(opt *redis.Options) *redis.Client {
	r := redisNewClient(opt)
	if r == nil {
		return nil
	}
	addr := opt.Addr
	r.WrapProcess(func(raw func(redis.Cmder) error) func(redis.Cmder) error {
		return getProcessWrapper("redis.Client.process", addr, raw)
	})
	r.WrapProcessPipeline(func(raw func([]redis.Cmder) error) func([]redis.Cmder) error {
		return getProcessPipelineWrapper("redisClient.Pipline", addr, raw)
	})
	return r
}

func getProcessWrapper(entry, addr string, raw func(redis.Cmder) error) func(redis.Cmder) error {
	return func(cmd redis.Cmder) error {
		begin := time.Now()
		req := tingyun3.LocalGet(StorageIndexRedis)
		if req == nil {
			tingyun3.LocalSet(StorageIndexRedis, 1)
		}
		err := raw(cmd)
		if req == nil {
			tingyun3.LocalDelete(StorageIndexRedis)
			command, object := parseCmder(cmd)
			handleGoRedis(addr, command, object, begin, err, 2)
			if tingyun3.GetAction() == nil {
				tingyun3.LocalClear()
			}
		}
		return err
	}
}
func getProcessPipelineWrapper(entry string, addr string, raw func([]redis.Cmder) error) func([]redis.Cmder) error {
	return func(cmds []redis.Cmder) error {
		begin := time.Now()
		req := tingyun3.LocalGet(StorageIndexRedis)
		if req == nil && len(cmds) > 0 {
			tingyun3.LocalSet(StorageIndexRedis, 1)
		}
		err := raw(cmds)
		if req == nil && len(cmds) > 0 {
			tingyun3.LocalDelete(StorageIndexRedis)
			cmd, object := parseCmders(cmds)
			handleGoRedis(addr, cmd, object, begin, err, 2)
			if tingyun3.GetAction() == nil {
				tingyun3.LocalClear()
			}
		}
		return err
	}
}

//go:noinline
func redisNewClusterClient(opt *redis.ClusterOptions) *redis.ClusterClient {
	trampoline.arg10 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapredisNewClusterClient(opt *redis.ClusterOptions) *redis.ClusterClient {
	r := redisNewClusterClient(opt)
	if r == nil {
		return r
	}
	addr := strings.Join(opt.Addrs, ",")
	r.WrapProcess(func(raw func(redis.Cmder) error) func(redis.Cmder) error {
		return getProcessWrapper("redis.ClusterClient.process", addr, raw)
	})
	r.WrapProcessPipeline(func(raw func([]redis.Cmder) error) func([]redis.Cmder) error {
		return getProcessPipelineWrapper("redisClusterClient.Pipline", addr, raw)
	})
	return r
}

//go:noinline
func redisNewFailoverClient(failoverOpt *redis.FailoverOptions) *redis.Client {
	trampoline.arg10 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapredisNewFailoverClient(failoverOpt *redis.FailoverOptions) *redis.Client {

	r := redisNewFailoverClient(failoverOpt)
	if r == nil {
		return nil
	}
	addr := strings.Join(failoverOpt.SentinelAddrs, ",")
	r.WrapProcess(func(raw func(redis.Cmder) error) func(redis.Cmder) error {
		return getProcessWrapper("redis.Client.process", addr, raw)
	})
	r.WrapProcessPipeline(func(raw func([]redis.Cmder) error) func([]redis.Cmder) error {
		return getProcessPipelineWrapper("redisClient.Pipline", addr, raw)
	})
	return r
}

type instanceSet struct {
	lock  sync.RWMutex
	items map[*redis.Client]string
}

func (d *instanceSet) init() {
	d.items = map[*redis.Client]string{}
}
func (d *instanceSet) get(c *redis.Client) string {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if i, found := d.items[c]; found {
		return i
	}
	return ""
}
func (d *instanceSet) Set(c *redis.Client, address string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.items[c] = address
}
func (d *instanceSet) remove(c *redis.Client) {
	d.lock.Lock()
	defer d.lock.Unlock()
	delete(d.items, c)
}

var dbs instanceSet

func init() {
	dbs.init()

	tingyun3.Register(reflect.ValueOf(WrapbaseClientProcess).Pointer())

	tingyun3.Register(reflect.ValueOf(WrapredisNewClient).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapredisNewClusterClient).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapredisNewFailoverClient).Pointer())

	tingyun3.Register(reflect.ValueOf(initTrampoline).Pointer())
}
