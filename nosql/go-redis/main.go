// +build linux
// +build amd64

package redis

import (
	"context"
	"fmt"
	"reflect"
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

type hook struct {
	host string
}

func (h *hook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	tingyun3.LocalSet(9, &processContext{time.Now(), cmd.Name()})
	return ctx, nil
}

func (h *hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if c := tingyun3.LocalDelete(9); c != nil {
		if info, ok := c.(*processContext); ok {
			handleGoRedis(h.host, cmd.Args(), info.begin)
		}
	}
	return nil
}

func (h *hook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if len(cmds) > 0 {
		tingyun3.LocalSet(9, &processContext{time.Now(), cmds[0].Name()})
	}
	return ctx, nil
}

func (h *hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if c := tingyun3.LocalDelete(9); c != nil {
		if info, ok := c.(*processContext); ok {
			handleGoRedis(h.host, cmds[0].Args(), info.begin)
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

func handleGoRedis(host string, args []interface{}, begin time.Time) {
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
	component.End(1)
}

//go:noinline
func GoRedisNewClient(opt *redis.Options) *redis.Client {
	fmt.Println(opt.Addr)
	return nil
}

//go:noinline
func WrapGoRedisNewClient(opt *redis.Options) *redis.Client {
	r := GoRedisNewClient(opt)
	host := opt.Addr
	if r != nil {
		r.AddHook(&hook{host: host})
	}
	return r
}

//v7 && v8
//go:noinline
func baseClientprocess(c *baseClient, ctx context.Context, cmd redis.Cmder) error {
	fmt.Println("host: ", c.opt.Addr, ", cmd: ", cmd.Name())
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

type baseClientV6 struct {
	pool interface{}
	opt  *redis.Options
}

//go:noinline
func WrapbaseClientprocess(c *baseClient, ctx context.Context, cmd redis.Cmder) error {
	begin := time.Now()
	req := tingyun3.LocalGet(9)
	if req == nil {
		tingyun3.LocalSet(9, 1)
	}
	// fmt.Println("Routine:", tingyun3.GetGID(), " Wrap (*baseCLient).process Enter, ", begin)
	err := baseClientprocess(c, ctx, cmd)
	// fmt.Println("Routine:", tingyun3.GetGID(), " Wrap (*baseCLient).process Leave, ", time.Now())
	if req == nil {
		tingyun3.LocalDelete(9)
		handleGoRedis(c.opt.Addr, cmd.Args(), begin)
	}
	return err
}

//go:noinline
func baseClientprocessPipeline(c *baseClient, ctx context.Context, cmds []redis.Cmder) error {
	fmt.Println(c.opt.Addr)
	return nil
}

//go:noinline
func WrapbaseClientprocessPipeline(c *baseClient, ctx context.Context, cmds []redis.Cmder) error {
	begin := time.Now()
	req := tingyun3.LocalGet(9)
	if req == nil {
		tingyun3.LocalSet(9, 1)
	}
	// fmt.Println("Routine:", tingyun3.GetGID(), " Wrap (*baseCLient).processPipeline Enter, ", begin)
	e := baseClientprocessPipeline(c, ctx, cmds)
	// fmt.Println("Routine:", tingyun3.GetGID(), " Wrap (*baseCLient).processPipeline Leave, ", time.Now())
	if req == nil {
		tingyun3.LocalDelete(9)
		handleGoRedis(c.opt.Addr, cmds[0].Args(), begin)
	}
	return e
}

//v6
//go:noinline
func baseClientProcess(c *baseClientV6, cmd redis.Cmder) error {
	fmt.Println(cmd.Name())
	return nil
}

//go:noinline
func WrapbaseClientProcess(c *baseClientV6, cmd redis.Cmder) error {
	begin := time.Now()
	req := tingyun3.LocalGet(9)
	if req == nil {
		tingyun3.LocalSet(9, 1)
	}

	// fmt.Println("baseClientProcess Start cmd: ", cmd.Name())
	err := baseClientProcess(c, cmd)
	// fmt.Println("baseClientProcess End cmd: ", cmd.Name())
	if req == nil {
		tingyun3.LocalDelete(9)
		handleGoRedis(c.opt.Addr, cmd.Args(), begin)
	}
	return err
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
	tingyun3.Register(reflect.ValueOf(WrapGoRedisNewClient).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapbaseClientprocess).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapbaseClientProcess).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapbaseClientprocessPipeline).Pointer())
}