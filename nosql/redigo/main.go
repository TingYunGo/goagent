// +build linux
// +build amd64

package redis

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/TingYunGo/goagent"
	"github.com/TingYunGo/goagent/libs/tystring"
)

type instanceSet struct {
	lock  sync.RWMutex
	items map[uintptr]string
}

type dbinstanceSet struct {
	dbset [8192]instanceSet
}

func (d *dbinstanceSet) init() {
	for k := range d.dbset {
		db := &(d.dbset[k])
		db.init()
	}
}

func (d *dbinstanceSet) get(conn uintptr) string {
	return d.dbset[conn%8192].get(conn)
}
func (d *dbinstanceSet) Set(conn uintptr, address string) {
	d.dbset[conn%8192].Set(conn, address)
}
func (d *dbinstanceSet) remove(conn uintptr) {
	d.dbset[conn%8192].remove(conn)
}

func interfaceToptr(i interface{}) uintptr {
	return reflect.ValueOf(i).Pointer()
}
func (d *instanceSet) init() {
	d.items = map[uintptr]string{}
}
func (d *instanceSet) get(conn uintptr) string {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if i, found := d.items[conn]; found {
		return i
	}
	return ""
}
func (d *instanceSet) Set(conn uintptr, address string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.items[conn] = address
}
func (d *instanceSet) remove(conn uintptr) {
	d.lock.Lock()
	defer d.lock.Unlock()
	delete(d.items, conn)
}

//go:noinline
func getCallName(skip int) (callerName string) {
	skip++
	callerName = tingyun3.GetCallerName(skip)
	token := "github.com/gomodule/redigo/redis"
	for tystring.SubString(callerName, 0, len(token)) == token {
		skip++
		callerName = tingyun3.GetCallerName(skip)
	}
	return
}

var dbs dbinstanceSet

//go:noinline
func redigoDial(network, address string, options ...redigo.DialOption) (redigo.Conn, error) {
	fmt.Println(network, address, options)
	return nil, nil
}

//go:noinline
func WrapredigoDial(network, address string, options ...redigo.DialOption) (redigo.Conn, error) {
	c, e := redigoDial(network, address, options...)
	if c != nil {
		dbs.Set(interfaceToptr(c), address)
	}
	return c, e
}

//go:noinline
func RedigoDialContext(ctx context.Context, network, address string, options ...redigo.DialOption) (redigo.Conn, error) {
	fmt.Println("network:", network, ", address:", address)
	return nil, nil
}

//go:noinline
func WrapRedigoDialContext(ctx context.Context, network, address string, options ...redigo.DialOption) (redigo.Conn, error) {
	c, e := RedigoDialContext(ctx, network, address, options...)
	if c != nil {
		dbs.Set(interfaceToptr(c), address)
	}
	return c, e
}

//go:noinline
func RedigoConnClose(conn uintptr) error {
	fmt.Println("conn:", conn)
	return nil
}

//go:noinline
func WrapRedigoConnClose(conn uintptr) error {
	dbs.remove(conn)
	e := RedigoConnClose(conn)
	return e
}

//有序列表,大小写不敏感, 放弃抓取对象的命令列表
var objectSkipList = []string{
	"AUTH",
	"ECHO",
	"PING",
	"QUIT",
	"SELECT",
}

func coreRedigoDoWithTimeout(begin time.Time, c uintptr, readTimeout time.Duration, cmd string, args []interface{}, r interface{}, err error) {
	action := tingyun3.GetAction()
	if action == nil {
		return
	}
	host := dbs.get(c)
	if host == "" {
		host = "UNKNOWN"
	}
	callerName := getCallName(3)
	object := ""
	if len(args) > 0 && tystring.FindString(objectSkipList, cmd) == -1 {
		if o, ok := args[0].(string); ok {
			object = o
		}
	}
	component := action.CreateRedisComponent(host, cmd, object, callerName)
	component.FixBegin(begin)
	if err != nil {
		component.SetError(err, callerName, 3)
	}
	component.FixStackEnd(2, func(funcName string) bool {
		token := "github.com/gomodule/redigo/"
		return tystring.SubString(funcName, 0, len(token)) == token
	})
	return
}

//go:noinline
func RedigoDoWithTimeout(c uintptr, readTimeout time.Duration, cmd string, args ...interface{}) (interface{}, error) {
	fmt.Println(cmd)
	return nil, nil
}

//go:noinline
func WrapRedigoDoWithTimeout(c uintptr, readTimeout time.Duration, cmd string, args ...interface{}) (interface{}, error) {
	begin := time.Now()
	res, err := RedigoDoWithTimeout(c, readTimeout, cmd, args...)
	coreRedigoDoWithTimeout(begin, c, readTimeout, cmd, args, res, err)
	return res, err
}
func init() {
	dbs.init()
	tingyun3.Register(reflect.ValueOf(WrapRedigoDialContext).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapredigoDial).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapRedigoConnClose).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapRedigoDoWithTimeout).Pointer())
}
