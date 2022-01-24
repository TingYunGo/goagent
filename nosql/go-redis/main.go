// +build linux
// +build amd64
// +build cgo

package redis

import (
	"context"
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

type Hooks struct {
	host string
}

func (h Hooks) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if ctx == nil {
		return ctx, nil
	}
	if pctx := ctx.Value("TingYunGoRedisCtx"); pctx != nil {
		if info, ok := pctx.(*processContext); ok {
			if info.cmd == cmdProcessHook {
				info.layer++
			}
		}
		return ctx, nil
	}
	return context.WithValue(ctx, "TingYunGoRedisCtx", &processContext{time.Now(), cmdProcessHook, 0}), nil
}

func (h Hooks) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if ctx == nil {
		return nil
	}
	if pctx := ctx.Value("TingYunGoRedisCtx"); pctx != nil {
		if info, ok := pctx.(*processContext); ok {
			if info.cmd != cmdProcessHook {
				return nil
			}
			if info.layer != 0 {
				info.layer--
				return nil
			}
			c, object := parseCmder(cmd)
			handleGoRedis(ctx, h.host, c, object, info.begin, cmd.Err(), 2)
		}
	}
	return nil
}

func (h Hooks) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if ctx == nil {
		return ctx, nil
	}
	if pctx := ctx.Value("TingYunGoRedisCtx"); pctx != nil {
		if info, ok := pctx.(*processContext); ok {
			if info.cmd == cmdProcessPipeHook {
				info.layer++
			}
		}
		return ctx, nil
	}
	configFlag := readConfigInt("go-redis.flag", 0)

	if (configFlag&1) == 0 && tingyun3.MatchCallerName(3, "github.com/go-redis/redis/v8.(*clusterStateHolder).LazyReload.func1") {
		return ctx, nil
	}
	return context.WithValue(ctx, "TingYunGoRedisCtx", &processContext{time.Now(), cmdProcessPipeHook, 0}), nil
}

func (h Hooks) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if ctx == nil {
		return nil
	}
	if pctx := ctx.Value("TingYunGoRedisCtx"); pctx != nil {
		if info, ok := pctx.(*processContext); ok {
			if info.cmd != cmdProcessPipeHook {
				return nil
			}
			if info.layer != 0 {
				info.layer--
				return nil
			}
			cmd, object := parseCmders(cmds)
			handleGoRedis(ctx, h.host, cmd, object, info.begin, nil, 2)
		}
	}
	return nil
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

func handleGoRedis(ctx context.Context, host, cmd, object string, begin time.Time, err error, skip int) {
	action, _ := tingyun3.FindAction(ctx)
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

//go:noinline
func baseClientprocess(c *baseClient, ctx context.Context, cmd redis.Cmder) error {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
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

//go:noinline
func WrapbaseClientprocess(c *baseClient, ctx context.Context, cmd redis.Cmder) error {
	if ctx != nil {
		if pctx := ctx.Value("TingYunGoRedisCtx"); pctx != nil {
			if _, ok := pctx.(*processContext); ok {
				return baseClientprocess(c, ctx, cmd)
			}
		}
	}
	begin := time.Now()
	req := tingyun3.LocalGet(StorageIndexRedis)
	var err error = nil
	if req == nil {
		tingyun3.LocalSet(StorageIndexRedis, 1)
		defer func() {
			command, object := parseCmder(cmd)
			handleGoRedis(ctx, c.opt.Addr, command, object, begin, err, 2)
			tingyun3.LocalDelete(StorageIndexRedis)
			if tingyun3.GetAction() == nil {
				tingyun3.LocalClear()
			}
		}()
	}
	ctx = context.WithValue(ctx, "TingYunGoRedisCtx", &processContext{time.Now(), cmdbaseClientprocess, 0})

	err = baseClientprocess(c, ctx, cmd)
	return err
}

//go:noinline
func baseClientprocessPipeline(c *baseClient, ctx context.Context, cmds []redis.Cmder) error {
	trampoline.arg2 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapbaseClientprocessPipeline(c *baseClient, ctx context.Context, cmds []redis.Cmder) error {
	if ctx != nil {
		if pctx := ctx.Value("TingYunGoRedisCtx"); pctx != nil {
			if _, ok := pctx.(*processContext); ok {
				return baseClientprocessPipeline(c, ctx, cmds)
			}
		}
	}
	begin := time.Now()
	req := tingyun3.LocalGet(StorageIndexRedis)
	var e error = nil
	if req == nil {
		tingyun3.LocalSet(StorageIndexRedis, 1)
		defer func() {
			tingyun3.LocalDelete(StorageIndexRedis)
			cmd, object := parseCmders(cmds)
			handleGoRedis(ctx, c.opt.Addr, cmd, object, begin, e, 2)
			if tingyun3.GetAction() == nil {
				tingyun3.LocalClear()
			}
		}()
	}
	ctx = context.WithValue(ctx, "TingYunGoRedisCtx", &processContext{time.Now(), cmdbaseClientprocessPipeline, 0})

	e = baseClientprocessPipeline(c, ctx, cmds)
	return e
}

type pipelineProcessor func(context.Context, uintptr, []redis.Cmder) (bool, error)

//go:noinline
func baseClientgeneralProcessPipeline(c *baseClient, ctx context.Context, cmds []redis.Cmder, p pipelineProcessor) error {
	trampoline.arg3 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapbaseClientgeneralProcessPipeline(c *baseClient, ctx context.Context, cmds []redis.Cmder, p pipelineProcessor) error {
	if ctx != nil {
		if pctx := ctx.Value("TingYunGoRedisCtx"); pctx != nil {
			if _, ok := pctx.(*processContext); ok {
				return baseClientgeneralProcessPipeline(c, ctx, cmds, p)
			}
		}
	}

	begin := time.Now()
	req := tingyun3.LocalGet(StorageIndexRedis)
	var e error = nil
	if req == nil {
		tingyun3.LocalSet(StorageIndexRedis, 1)
		defer func() {
			tingyun3.LocalDelete(StorageIndexRedis)
			cmd, object := parseCmders(cmds)
			handleGoRedis(ctx, c.opt.Addr, cmd, object, begin, e, 2)
			if tingyun3.GetAction() == nil {
				tingyun3.LocalClear()
			}
		}()
	}
	ctx = context.WithValue(ctx, "TingYunGoRedisCtx", &processContext{time.Now(), cmdbaseClientgeneralProcessPipeline, 0})

	e = baseClientgeneralProcessPipeline(c, ctx, cmds, p)
	return e
}

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
		handleGoRedis(nil, c.opt.Addr, command, object, begin, err, 2)
		if tingyun3.GetAction() == nil {
			tingyun3.LocalClear()
		}
	}
	return err
}

//go:noinline
func redisClientWrapProcess(c *redis.Client, fn func(func(redis.Cmder) error) func(redis.Cmder) error) {
	trampoline.arg5 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
}

//go:noinline
func WrapredisClientWrapProcess(c *redis.Client, fn func(func(redis.Cmder) error) func(redis.Cmder) error) {
	redisClientWrapProcess(c, fn)
}

//go:noinline
func redisClientWrapProcessPipeline(c *redis.Client, fn func(oldProcess func([]redis.Cmder) error) func([]redis.Cmder) error) {
	trampoline.arg6 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
}

//go:noinline
func WrapredisClientWrapProcessPipeline(c *redis.Client, fn func(oldProcess func([]redis.Cmder) error) func([]redis.Cmder) error) {
	redisClientWrapProcessPipeline(c, fn)
}

//go:noinline
func redisClusterClientWrapProcess(c *redis.ClusterClient, fn func(func(redis.Cmder) error) func(redis.Cmder) error) {
	trampoline.arg7 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
}

//go:noinline
func WrapredisClusterClientWrapProcess(c *redis.ClusterClient, fn func(func(redis.Cmder) error) func(redis.Cmder) error) {
	redisClusterClientWrapProcess(c, fn)
}

//go:noinline
func redisClusterClientWrapProcessPipeline(c *redis.ClusterClient, fn func(oldProcess func([]redis.Cmder) error) func([]redis.Cmder) error) {
	trampoline.arg8 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
}

//go:noinline
func WrapredisClusterClientWrapProcessPipeline(c *redis.ClusterClient, fn func(oldProcess func([]redis.Cmder) error) func([]redis.Cmder) error) {
	redisClusterClientWrapProcessPipeline(c, fn)
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
	redisClientWrapProcess(r, func(raw func(redis.Cmder) error) func(redis.Cmder) error {
		return getProcessWrapper("redis.Client.process", addr, raw)
	})
	redisClientWrapProcessPipeline(r, func(raw func([]redis.Cmder) error) func([]redis.Cmder) error {
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
			handleGoRedis(nil, addr, command, object, begin, err, 2)
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
			handleGoRedis(nil, addr, cmd, object, begin, err, 2)
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
	redisClusterClientWrapProcess(r, func(raw func(redis.Cmder) error) func(redis.Cmder) error {
		return getProcessWrapper("redis.ClusterClient.process", addr, raw)
	})
	redisClusterClientWrapProcessPipeline(r, func(raw func([]redis.Cmder) error) func([]redis.Cmder) error {
		return getProcessPipelineWrapper("redisClusterClient.Pipline", addr, raw)
	})
	return r
}

type clusterClient struct {
	opt           *redis.ClusterOptions
	nodes         *uint64
	state         *uint64
	cmdsInfoCache *uint64
}
type ClusterClient struct {
	*clusterClient
	p   func()
	ctx context.Context
}

//go:noinline
func ClusterClient_processPipeline(c *ClusterClient, ctx context.Context, cmds []redis.Cmder) error {
	trampoline.arg2 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20

	return nil
}

//go:noinline
func WrapClusterClient_processPipeline(c *ClusterClient, ctx context.Context, cmds []redis.Cmder) error {
	if ctx != nil {
		if pctx := ctx.Value("TingYunGoRedisCtx"); pctx != nil {
			if _, ok := pctx.(*processContext); ok {
				return ClusterClient_processPipeline(c, ctx, cmds)
			}
		}
	}

	begin := time.Now()
	req := tingyun3.LocalGet(StorageIndexRedis)
	var e error = nil
	if req == nil {
		tingyun3.LocalSet(StorageIndexRedis, 1)
		defer func() {
			addr := ""
			if len(c.opt.Addrs) == 0 {
				addr = "[]"
			} else if len(c.opt.Addrs) == 1 {
				addr = "[" + c.opt.Addrs[0] + "]"
			} else {
				addr = "[" + c.opt.Addrs[0] + ",...]"
			}
			tingyun3.LocalDelete(StorageIndexRedis)
			cmd, object := parseCmders(cmds)
			handleGoRedis(ctx, addr, cmd, object, begin, e, 2)
			if tingyun3.GetAction() == nil {
				tingyun3.LocalClear()
			}
		}()
	}
	ctx = context.WithValue(ctx, "TingYunGoRedisCtx", &processContext{time.Now(), cmdClusterClient_processPipeline, 0})
	e = ClusterClient_processPipeline(c, ctx, cmds)
	return e
}

//go:noinline
func ClusterClient_processTxPipeline(c *ClusterClient, ctx context.Context, cmds []redis.Cmder) error {
	trampoline.arg2 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20

	return nil
}

//go:noinline
func WrapClusterClient_processTxPipeline(c *ClusterClient, ctx context.Context, cmds []redis.Cmder) error {
	if ctx != nil {
		if pctx := ctx.Value("TingYunGoRedisCtx"); pctx != nil {
			if _, ok := pctx.(*processContext); ok {
				return ClusterClient_processTxPipeline(c, ctx, cmds)
			}
		}
	}
	begin := time.Now()
	req := tingyun3.LocalGet(StorageIndexRedis)
	var e error = nil
	if req == nil {
		tingyun3.LocalSet(StorageIndexRedis, 1)
		defer func() {
			addr := ""
			if len(c.opt.Addrs) == 0 {
				addr = "[]"
			} else if len(c.opt.Addrs) == 1 {
				addr = "[" + c.opt.Addrs[0] + "]"
			} else {
				addr = "[" + c.opt.Addrs[0] + ",...]"
			}
			tingyun3.LocalDelete(StorageIndexRedis)
			cmd, object := parseCmders(cmds)
			handleGoRedis(ctx, addr, cmd, object, begin, e, 2)
			if tingyun3.GetAction() == nil {
				tingyun3.LocalClear()
			}
		}()
	}
	ctx = context.WithValue(ctx, "TingYunGoRedisCtx", &processContext{time.Now(), cmdClusterClient_processTxPipeline, 0})
	e = ClusterClient_processTxPipeline(c, ctx, cmds)
	return e
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
	tingyun3.Register(reflect.ValueOf(WrapbaseClientprocess).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapbaseClientProcess).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapbaseClientprocessPipeline).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapbaseClientgeneralProcessPipeline).Pointer())

	tingyun3.Register(reflect.ValueOf(WrapredisNewClient).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapredisNewClusterClient).Pointer())

	tingyun3.Register(reflect.ValueOf(WrapredisClientWrapProcess).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapredisClientWrapProcessPipeline).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapredisClusterClientWrapProcess).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapredisClusterClientWrapProcessPipeline).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapClusterClient_processPipeline).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapClusterClient_processTxPipeline).Pointer())
	tingyun3.Register(reflect.ValueOf(initTrampoline).Pointer())
}
