// +build linux
// +build amd64 arm64
// +build cgo

package redisv9

import (
	"context"
	"reflect"
	"runtime"
	"time"

	"github.com/TingYunGo/goagent/libs/tystring"

	"github.com/TingYunGo/goagent"

	redis "github.com/redis/go-redis/v9"
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
	"net/http",
	"github.com/redis/go-redis",
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

func matchMethod(method, matcher string) bool {
	return tystring.SubString(method, 0, len(matcher)) == matcher
}
func isNativeMethod(method string) bool {

	if matchMethod(method, skipTokens[0]) {
		return true
	}
	if matchMethod(method, skipTokens[1]) {
		return true
	}
	if matchMethod(method, skipTokens[2]) {
		return true
	}
	return false
}

//go:noinline
func getCallName(skip int) (callerName string) {

	stackList := make([]uintptr, 8)
	count := runtime.Callers(skip, stackList)

	for i := 0; i < count; i++ {
		callerName = runtime.FuncForPC(stackList[i] - 1).Name()
		if !isNativeMethod(callerName) {
			break
		}
	}
	return
}

type Hooks struct {
	host string
}

var _ redis.Hook = Hooks{}

func (h Hooks) DialHook(next redis.DialHook) redis.DialHook {
	return next
}
func (h Hooks) ProcessHook(next redis.ProcessHook) redis.ProcessHook {

	return func(ctx context.Context, cmd redis.Cmder) error {
		handled := false
		var pctx *processContext = nil
		if ctx != nil {
			if pctx = fetchCtx(ctx); pctx != nil {
				handled = true
			}
		}
		if !handled {
			pctx = &processContext{time.Now(), cmdProcessHook, 0}
			if ctx == nil {
				ctx = context.Background()
			}
			ctx = context.WithValue(ctx, "TingYunGoRedisCtx", pctx)
		}
		err := next(ctx, cmd)
		if !handled {
			c, object := parseCmder(cmd)
			handleGoRedis(ctx, h.host, c, object, pctx.begin, cmd.Err(), 2)
		}
		return err
	}

}
func fetchCtx(ctx context.Context) *processContext {
	if pctx := ctx.Value("TingYunGoRedisCtx"); pctx != nil {
		if info, ok := pctx.(*processContext); ok {
			return info
		}
	}
	return nil

}
func (h Hooks) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {

	return func(ctx context.Context, cmds []redis.Cmder) error {
		handled := false
		var pctx *processContext = nil
		if ctx != nil {
			if pctx = fetchCtx(ctx); pctx != nil {
				handled = true
			}
		}

		if !handled {
			pctx = &processContext{time.Now(), cmdProcessPipeHook, 0}
			if ctx == nil {
				ctx = context.Background()
			}
			ctx = context.WithValue(ctx, "TingYunGoRedisCtx", pctx)
		}

		err := next(ctx, cmds)

		if !handled {
			cmd, object := parseCmders(cmds)
			handleGoRedis(ctx, h.host, cmd, object, pctx.begin, nil, 2)
		}

		return err
	}
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
	if component == nil {
		return
	}
	component.FixBegin(begin)
	if err != nil {
		component.SetException(err, callerName, 3)
	}
	component.FixStackEnd(skip, isNativeMethod)
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
		if pctx := fetchCtx(ctx); pctx != nil {
			return baseClientprocess(c, ctx, cmd)
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
func baseClientprocessPipeline(c *baseClient, ctx context.Context, cmds []redis.Cmder, p pipelineProcessor) error {
	trampoline.arg2 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

type pipelineProcessor func(context.Context, uintptr, []redis.Cmder) (bool, error)

//go:noinline
func WrapbaseClientprocessPipeline(c *baseClient, ctx context.Context, cmds []redis.Cmder, p pipelineProcessor) error {
	if ctx != nil {
		if pctx := fetchCtx(ctx); pctx != nil {
			return baseClientprocessPipeline(c, ctx, cmds, p)
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

	e = baseClientprocessPipeline(c, ctx, cmds, p)
	return e
}

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
		if pctx := fetchCtx(ctx); pctx != nil {
			return baseClientgeneralProcessPipeline(c, ctx, cmds, p)
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
	r.AddHook(Hooks{host: addr})
	return r
}

func seriallizeAddresses(addresses []string) string {
	if len(addresses) == 0 {
		return "[]"
	} else if len(addresses) == 1 {
		return "[" + addresses[0] + "]"
	} else {
		return "[" + addresses[0] + ",...]"
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
	addr := seriallizeAddresses(opt.Addrs)
	r.AddHook(Hooks{host: addr})
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
	addr := seriallizeAddresses(failoverOpt.SentinelAddrs)
	r.AddHook(Hooks{host: addr})

	return r
}

//go:noinline
func redisNewFailoverClusterClient(failoverOpt *redis.FailoverOptions) *redis.ClusterClient {
	trampoline.arg10 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapredisNewFailoverClusterClient(failoverOpt *redis.FailoverOptions) *redis.ClusterClient {

	r := redisNewFailoverClusterClient(failoverOpt)
	if r == nil {
		return nil
	}
	addr := seriallizeAddresses(failoverOpt.SentinelAddrs)
	r.AddHook(Hooks{host: addr})
	return r
}

func init() {
	tingyun3.Register(reflect.ValueOf(WrapbaseClientprocess).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapbaseClientprocessPipeline).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapbaseClientgeneralProcessPipeline).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapredisNewClient).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapredisNewClusterClient).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapredisNewFailoverClient).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapredisNewFailoverClusterClient).Pointer())

	tingyun3.Register(reflect.ValueOf(initTrampoline).Pointer())
}
