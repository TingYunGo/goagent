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

var skipTokens = []string{
	"github.com/go-redis/redis",
	"github.com/TingYunGo/goagent",
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

var _ redis.Hook = Hooks{}

func (h Hooks) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if ctx == nil {
		return ctx, nil
	}
	return context.WithValue(ctx, "TingYunProcessCtx", &processContext{time.Now(), cmd.Name()}), nil
}

func (h Hooks) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if ctx == nil {
		return nil
	}
	if pctx := ctx.Value("TingYunProcessCtx"); pctx != nil {
		if info, ok := pctx.(*processContext); ok {
			handleGoRedis(h.host, cmd.Args(), info.begin, cmd.Err(), 2)
		}
	}
	return nil
}

func (h Hooks) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if ctx == nil {
		return ctx, nil
	}
	return context.WithValue(ctx, "TingYunPipeCtx", &processContext{time.Now(), cmds[0].Name()}), nil
}

func (h Hooks) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if ctx == nil {
		return nil
	}
	if pctx := ctx.Value("TingYunPipeCtx"); pctx != nil {
		if info, ok := pctx.(*processContext); ok {
			handleGoRedis(h.host, cmds[0].Args(), info.begin, nil, 2)
		}
	}
	return nil
}

type processContext struct {
	begin time.Time
	cmd   string
}

var objectSkipList = []string{
	"AUTH",
	"ECHO",
	"PING",
	"QUIT",
	"SELECT",
}

func handleGoRedis(host string, args []interface{}, begin time.Time, err error, skip int) {
	action := tingyun3.GetAction()
	if action == nil {
		return
	}
	callerName := getCallName(3)
	cmd, obj := getArgs(args)
	object := ""
	if len(obj) > 0 && tystring.FindString(objectSkipList, cmd) == -1 {
		object = obj
	}
	component := action.CreateRedisComponent(host, cmd, object, callerName)
	component.FixBegin(begin)
	if err != nil {
		component.SetException(err, callerName, 3)
	}
	component.FixStackEnd(skip, func(funcname string) bool {
		token := "github.com/go-redis/redis/"
		return tystring.SubString(funcname, 0, len(token)) == token
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
	if argc > 0 {
		cmd = args[0].(string)
	} else {
		cmd = ""
	}
	if argc > 1 {
		object = args[1].(string)
	} else {
		object = ""
	}
	return
}

type baseClient struct {
	opt  *redis.Options
	pool interface{}
}

//go:noinline
func WrapbaseClientprocess(c *baseClient, ctx context.Context, cmd redis.Cmder) error {
	begin := time.Now()
	req := tingyun3.LocalGet(9)
	var err error = nil
	if req == nil {
		tingyun3.LocalSet(9, 1)
		defer func() {
			tingyun3.LocalDelete(9)
			handleGoRedis(c.opt.Addr, cmd.Args(), begin, err, 2)
		}()
	}
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
	begin := time.Now()
	req := tingyun3.LocalGet(9)
	var e error = nil
	if req == nil {
		tingyun3.LocalSet(9, 1)
		defer func() {
			tingyun3.LocalDelete(9)
			handleGoRedis(c.opt.Addr, cmds[0].Args(), begin, e, 2)
		}()
	}
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
	begin := time.Now()
	req := tingyun3.LocalGet(9)
	var e error = nil
	if req == nil {
		tingyun3.LocalSet(9, 1)
		defer func() {
			tingyun3.LocalDelete(9)
			handleGoRedis(c.opt.Addr, cmds[0].Args(), begin, e, 2)
		}()
	}
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
	req := tingyun3.LocalGet(9)
	if req == nil {
		tingyun3.LocalSet(9, 1)
	}

	err := baseClientProcess(c, cmd)
	if req == nil {
		tingyun3.LocalDelete(9)
		handleGoRedis(c.opt.Addr, cmd.Args(), begin, err, 2)
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
		req := tingyun3.LocalGet(9)
		if req == nil {
			tingyun3.LocalSet(9, 1)
		}
		err := raw(cmd)
		if req == nil {
			tingyun3.LocalDelete(9)
			handleGoRedis(addr, cmd.Args(), begin, err, 2)
		}
		return err
	}
}
func getProcessPipelineWrapper(entry string, addr string, raw func([]redis.Cmder) error) func([]redis.Cmder) error {
	return func(cmds []redis.Cmder) error {
		begin := time.Now()
		req := tingyun3.LocalGet(9)
		if req == nil && len(cmds) > 0 {
			tingyun3.LocalSet(9, 1)
		}
		err := raw(cmds)
		if req == nil && len(cmds) > 0 {
			tingyun3.LocalDelete(9)
			handleGoRedis(addr, cmds[0].Args(), begin, err, 2)
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
	begin := time.Now()
	req := tingyun3.LocalGet(9)
	var e error = nil
	if req == nil {
		tingyun3.LocalSet(9, 1)
		defer func() {
			addr := ""
			if len(c.opt.Addrs) == 0 {
				addr = "[]"
			} else if len(c.opt.Addrs) == 1 {
				addr = "[" + c.opt.Addrs[0] + "]"
			} else {
				addr = "[" + c.opt.Addrs[0] + ",...]"
			}
			tingyun3.LocalDelete(9)
			handleGoRedis(addr, cmds[0].Args(), begin, e, 2)
		}()
	}
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
	begin := time.Now()
	req := tingyun3.LocalGet(9)
	var e error = nil
	if req == nil {
		tingyun3.LocalSet(9, 1)
		defer func() {
			addr := ""
			if len(c.opt.Addrs) == 0 {
				addr = "[]"
			} else if len(c.opt.Addrs) == 1 {
				addr = "[" + c.opt.Addrs[0] + "]"
			} else {
				addr = "[" + c.opt.Addrs[0] + ",...]"
			}
			tingyun3.LocalDelete(9)
			handleGoRedis(addr, cmds[0].Args(), begin, e, 2)
		}()
	}
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
