// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64

package database

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"github.com/TingYunGo/goagent"
	"github.com/TingYunGo/goagent/libs/tystring"
)

type dbinstanceSet struct {
	lock  sync.RWMutex
	items map[*sql.DB]databaseInfo
}

func (d *dbinstanceSet) init() {
	d.items = map[*sql.DB]databaseInfo{}
}
func (d *dbinstanceSet) Get(db *sql.DB) *databaseInfo {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if i, found := d.items[db]; found {
		return &i
	}
	return nil
}
func (d *dbinstanceSet) Set(db *sql.DB, info databaseInfo) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.items == nil {
		d.init()
	}
	d.items[db] = info
}
func (d *dbinstanceSet) Delete(db *sql.DB) {
	d.lock.Lock()
	defer d.lock.Unlock()
	delete(d.items, db)
}

var dbs dbinstanceSet

type databaseContext struct {
	records map[*sql.Rows]*tingyun3.Component
	stmts   map[*sql.Stmt]*tingyun3.Component
}

func (ctx *databaseContext) init() *databaseContext {
	ctx.records = map[*sql.Rows]*tingyun3.Component{}
	ctx.stmts = map[*sql.Stmt]*tingyun3.Component{}
	return ctx
}
func (ctx *databaseContext) empty() bool {
	return len(ctx.records) == 0 && len(ctx.stmts) == 0
}
func (ctx *databaseContext) clear() {
	if ctx.records != nil {
		for s := range ctx.records {
			delete(ctx.records, s)
		}
		ctx.records = nil
	}
	if ctx.stmts != nil {
		for s := range ctx.stmts {
			delete(ctx.stmts, s)
		}
		ctx.stmts = nil
	}
}
func getTingyunDBType(name string) uint8 {
	if matchVendor(name, "mysql") {
		return tingyun3.ComponentMysql
	} else if matchVendor(name, "postgre") {
		return tingyun3.ComponentPostgreSQL
	} else if matchVendor(name, "sqlserver") {
		return tingyun3.ComponentMSSQL
	} else if matchVendor(name, "sqlite") {
		return tingyun3.ComponentSQLite
	} else if matchVendor(name, "oci") || matchVendor(name, "godror") {
		return tingyun3.ComponentOracle
	}
	return tingyun3.ComponentDefaultDB
}

//go:noinline
func DBOpen(driverName, dataSourceName string) (*sql.DB, error) {
	fmt.Println(driverName, dataSourceName)
	return nil, nil
}

//go:noinline
func WrapDBOpen(driverName, dataSourceName string) (*sql.DB, error) {
	db, err := DBOpen(driverName, dataSourceName)
	if db != nil {
		info := databaseInfo{}
		info.init(driverName, dataSourceName)
		dbs.Set(db, info)
	}
	return db, err
}

//go:noinline
func DBClose(db *sql.DB) error {
	fmt.Println(db)
	return nil
}

//go:noinline
func WrapDBClose(db *sql.DB) error {
	dbs.Delete(db)
	return db.Close()
}

func coreWrapPrepareContext(begin time.Time, db *sql.DB, query string, stmt *sql.Stmt, e error) {
	action := tingyun3.GetAction()
	if action == nil {
		return
	}
	info := dbs.Get(db)
	if info == nil {
		return
	}
	var dbctx *databaseContext = nil
	c := tingyun3.LocalGet(1)
	if c != nil {
		dbctx = c.(*databaseContext)
		if _, found := dbctx.stmts[stmt]; found {
			return
		}
	}
	callerName := getCallName(3)
	component := action.CreateSQLComponent(getTingyunDBType(info.vender), info.host, info.dbname, query, callerName)
	component.FixBegin(begin)
	if stmt == nil || e != nil {
		component.SetException(e, callerName, 3)
		component.End(2)
		return
	}
	if dbctx == nil {
		dbctx = (&databaseContext{}).init()
		tingyun3.LocalSet(1, dbctx)
	}
	action.OnEnd(func() {
		dbctx.clear()
		tingyun3.LocalDelete(1)
	})
	dbctx.stmts[stmt] = component
}
func coreWrapExecContext(begin time.Time, db *sql.DB, query string, r sql.Result, e error) {
	action := tingyun3.GetAction()
	if action == nil {
		return
	}
	info := dbs.Get(db)
	if info == nil {
		return
	}
	callerName := getCallName(3)
	component := action.CreateSQLComponent(getTingyunDBType(info.vender), info.host, info.dbname, query, callerName)
	component.FixBegin(begin)
	if r == nil && e != nil {
		component.SetException(e, callerName, 3)
		component.End(2)
		return
	}
	component.End(2)
}
func coreWrapQueryContext(begin time.Time, db *sql.DB, query string, r *sql.Rows, e error) {
	action := tingyun3.GetAction()
	if action == nil {
		return
	}
	info := dbs.Get(db)
	if info == nil {
		return
	}
	var dbctx *databaseContext = nil
	c := tingyun3.LocalGet(1)
	if c != nil {
		dbctx = c.(*databaseContext)
		if _, found := dbctx.records[r]; found {
			return
		}
	}
	callerName := getCallName(3)
	component := action.CreateSQLComponent(getTingyunDBType(info.vender), info.host, info.dbname, query, callerName)
	component.FixBegin(begin)
	if r == nil && e != nil {
		component.SetException(e, callerName, 3)
		component.End(2)
		return
	}
	component.End(2)
	if dbctx == nil {
		dbctx = (&databaseContext{}).init()
		tingyun3.LocalSet(1, dbctx)
	}
	dbctx.records[r] = component
}

//go:noinline
func DBPrepareContext(db *sql.DB, ctx context.Context, query string) (*sql.Stmt, error) {
	fmt.Println(query)
	return nil, nil
}

//go:noinline
func WrapDBPrepareContext(db *sql.DB, ctx context.Context, query string) (*sql.Stmt, error) {
	recursiveChecker := &recursiveCheck{rlsID: 2, success: false}
	begin, enter := recursiveChecker.enter()
	defer func() {
		recursiveChecker.leave()
		if exception := recover(); exception != nil {
			panic(exception)
		}
	}()
	stmt, e := DBPrepareContext(db, ctx, query)
	if enter {
		coreWrapPrepareContext(begin, db, query, stmt, e)
	}
	return stmt, e
}

//go:noinline
func DBExecContext(db *sql.DB, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	fmt.Println(query)
	return nil, nil
}

//go:noinline
func WrapDBExecContext(db *sql.DB, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {

	recursiveChecker := &recursiveCheck{rlsID: 2, success: false}

	begin, enter := recursiveChecker.enter()
	defer func() {
		recursiveChecker.leave()
		if exception := recover(); exception != nil {
			panic(exception)
		}
	}()

	r, e := DBExecContext(db, ctx, query, args...)
	if enter {
		coreWrapExecContext(begin, db, query, r, e)
	}
	return r, e
}

//go:noinline
func DBQueryContext(db *sql.DB, ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	fmt.Println(query)
	return nil, nil
}

//go:noinline
func getCallName(skip int) (callerName string) {
	skip++
	callerName = tingyun3.GetCallerName(skip)
	for isNativeMethod(callerName) {
		skip++
		callerName = tingyun3.GetCallerName(skip)
	}
	return
}

//go:noinline
func WrapDBQueryContext(db *sql.DB, ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	recursiveChecker := &recursiveCheck{rlsID: 2, success: false}

	begin, enter := recursiveChecker.enter()
	defer func() {
		recursiveChecker.leave()
		if exception := recover(); exception != nil {
			panic(exception)
		}
	}()
	r, e := DBQueryContext(db, ctx, query, args...)
	if enter {
		coreWrapQueryContext(begin, db, query, r, e)
	}
	return r, e
}

//go:noinline
func ConnExecContext(c *sql.Conn, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	fmt.Println(query)
	return nil, nil
}

//go:noinline
func WrapConnExecContext(c *sql.Conn, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	recursiveChecker := &recursiveCheck{rlsID: 2, success: false}

	begin, enter := recursiveChecker.enter()
	defer func() {
		recursiveChecker.leave()
		if exception := recover(); exception != nil {
			panic(exception)
		}
	}()
	r, e := ConnExecContext(c, ctx, query, args...)
	if enter {
		coreWrapExecContext(begin, getdb_byconn(c), query, r, e)
	}
	return r, e
}

//go:noinline
func ConnQueryContext(c *sql.Conn, ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	fmt.Println(query)
	return nil, nil
}

type connDB struct {
	db      *sql.DB
	closemu sync.RWMutex
	dc      unsafe.Pointer
	done    int32
}

func getdb_byconn(c *sql.Conn) *sql.DB {
	return (*connDB)(unsafe.Pointer(c)).db
}

//go:noinline
func WrapConnQueryContext(c *sql.Conn, ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	recursiveChecker := &recursiveCheck{rlsID: 2, success: false}

	begin, enter := recursiveChecker.enter()
	defer func() {
		recursiveChecker.leave()
		if exception := recover(); exception != nil {
			panic(exception)
		}
	}()
	r, e := ConnQueryContext(c, ctx, query, args...)
	if enter {
		coreWrapQueryContext(begin, getdb_byconn(c), query, r, e)
	}
	return r, e
}

//go:noinline
func ConnPrepareContext(c *sql.Conn, ctx context.Context, query string) (*sql.Stmt, error) {
	fmt.Println(query)
	return nil, nil
}

//go:noinline
func WrapConnPrepareContext(c *sql.Conn, ctx context.Context, query string) (*sql.Stmt, error) {
	recursiveChecker := &recursiveCheck{rlsID: 2, success: false}

	begin, enter := recursiveChecker.enter()
	defer func() {
		recursiveChecker.leave()
		if exception := recover(); exception != nil {
			panic(exception)
		}
	}()
	stmt, e := ConnPrepareContext(c, ctx, query)
	if enter {
		coreWrapPrepareContext(begin, getdb_byconn(c), query, stmt, e)
	}
	return stmt, e
}

type sqlTx struct {
	db *sql.DB
}

func getdbByTx(c *sql.Tx) *sql.DB {
	return (*sqlTx)(unsafe.Pointer(c)).db
}

//go:noinline
func TxPrepareContext(tx *sql.Tx, ctx context.Context, query string) (*sql.Stmt, error) {
	fmt.Println(query)
	return nil, nil
}

//go:noinline
func WrapTxPrepareContext(tx *sql.Tx, ctx context.Context, query string) (*sql.Stmt, error) {
	recursiveChecker := &recursiveCheck{rlsID: 2, success: false}

	begin, enter := recursiveChecker.enter()
	defer func() {
		recursiveChecker.leave()
		if exception := recover(); exception != nil {
			panic(exception)
		}
	}()
	stmt, e := TxPrepareContext(tx, ctx, query)
	if enter {
		coreWrapPrepareContext(begin, getdbByTx(tx), query, stmt, e)
	}
	return stmt, e
}

//go:noinline
func TxExecContext(tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	fmt.Println(query)
	return nil, nil
}

//go:noinline
func WrapTxExecContext(tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	recursiveChecker := &recursiveCheck{rlsID: 2, success: false}

	begin, enter := recursiveChecker.enter()
	defer func() {
		recursiveChecker.leave()
		if exception := recover(); exception != nil {
			panic(exception)
		}
	}()
	r, e := TxExecContext(tx, ctx, query, args...)
	if enter {
		coreWrapExecContext(begin, getdbByTx(tx), query, r, e)
	}
	return r, e
}

//go:noinline
func TxQueryContext(tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	fmt.Println(query)
	return nil, nil
}

//go:noinline
func WrapTxQueryContext(tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	recursiveChecker := &recursiveCheck{rlsID: 2, success: false}

	begin, enter := recursiveChecker.enter()
	defer func() {
		recursiveChecker.leave()
		if exception := recover(); exception != nil {
			panic(exception)
		}
	}()
	r, e := TxQueryContext(tx, ctx, query, args...)
	if enter {
		coreWrapQueryContext(begin, getdbByTx(tx), query, r, e)
	}
	return r, e
}

//go:noinline
func StmtQueryContext(s *sql.Stmt, ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	fmt.Println(s)
	return nil, nil
}

//go:noinline
func WrapStmtQueryContext(s *sql.Stmt, ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	r, e := StmtQueryContext(s, ctx, args...)
	if c := tingyun3.LocalGet(1); c != nil {
		dbctx := c.(*databaseContext)
		component, found := dbctx.stmts[s]
		if !found {
			return r, e
		}
		delete(dbctx.stmts, s)
		if r == nil && e != nil {
			callerName := getCallName(2)
			component.SetException(e, callerName, 2)
			if dbctx.empty() {
				dbctx.clear()
				tingyun3.LocalDelete(1)
			}
			return r, e
		}
		dbctx.records[r] = component
	}
	return r, e
}

//go:noinline
func StmtClose(s *sql.Stmt) error {
	fmt.Println(s)
	return nil
}

//go:noinline
func WrapStmtClose(s *sql.Stmt) error {
	err := StmtClose(s)

	if c := tingyun3.LocalGet(1); c != nil {
		dbctx := c.(*databaseContext)
		if c, found := dbctx.stmts[s]; found {
			c.End(1)
			delete(dbctx.stmts, s)
		}
		if dbctx.empty() {
			dbctx.clear()
			tingyun3.LocalDelete(1)
		}
	}
	return err
}

//go:noinline
func StmtExecContext(s *sql.Stmt, ctx context.Context, args ...interface{}) (sql.Result, error) {
	fmt.Println(s)
	return nil, nil
}

//go:noinline
func WrapStmtExecContext(s *sql.Stmt, ctx context.Context, args ...interface{}) (sql.Result, error) {
	r, e := StmtExecContext(s, ctx, args...)

	if c := tingyun3.LocalGet(1); c != nil {
		dbctx := c.(*databaseContext)
		if c, found := dbctx.stmts[s]; found {
			c.End(1)
		}
	}
	return r, e
}

//go:noinline
func RowsClose(rs *sql.Rows) error {
	fmt.Println(rs)
	return nil
}

//go:noinline
func WrapRowsClose(rs *sql.Rows) error {
	err := RowsClose(rs)
	if c := tingyun3.LocalGet(1); c != nil {
		dbctx := c.(*databaseContext)
		if c, found := dbctx.records[rs]; found {
			c.End(1)
			delete(dbctx.records, rs)
		}
		if dbctx.empty() {
			tingyun3.LocalDelete(1)
		}
	}
	return err
}

type driverConn struct {
	db        *sql.DB
	createdAt time.Time
}

//database/sql.(*DB).queryDC

//go:noinline
func DBqueryDC(db *sql.DB, ctx, txctx context.Context, dc *driverConn, releaseConn func(error), query string, args []interface{}) (*sql.Rows, error) {
	fmt.Println(query)
	return nil, nil
}

//go:noinline
func WrapDBqueryDC(db *sql.DB, ctx, txctx context.Context, dc *driverConn, releaseConn func(error), query string, args []interface{}) (*sql.Rows, error) {
	recursiveChecker := &recursiveCheck{rlsID: 2, success: false}

	begin, enter := recursiveChecker.enter()
	defer func() {
		recursiveChecker.leave()
		if exception := recover(); exception != nil {
			panic(exception)
		}
	}()
	r, e := DBqueryDC(db, ctx, txctx, dc, releaseConn, query, args)
	if enter {
		coreWrapQueryContext(begin, db, query, r, e)
	}
	return r, e
}

//go:noinline
func DBexecDC(db *sql.DB, ctx context.Context, dc *driverConn, release func(error), query string, args []interface{}) (res sql.Result, err error) {
	fmt.Println(query)
	return nil, nil
}

//go:noinline
func WrapDBexecDC(db *sql.DB, ctx context.Context, dc *driverConn, release func(error), query string, args []interface{}) (res sql.Result, err error) {
	recursiveChecker := &recursiveCheck{rlsID: 2, success: false}

	begin, enter := recursiveChecker.enter()
	defer func() {
		recursiveChecker.leave()
		if exception := recover(); exception != nil {
			panic(exception)
		}
	}()
	r, e := DBexecDC(db, ctx, dc, release, query, args)
	if enter {
		coreWrapExecContext(begin, db, query, r, e)
	}
	return r, e
}

type releaseConn func(error)
type stmtConnGrabber interface {
	grabConn(context.Context) (*driverConn, releaseConn, error)
	txCtx() context.Context
}

//go:noinline
func DBprepareDC(db *sql.DB, ctx context.Context, dc *driverConn, release func(error), cg stmtConnGrabber, query string) (*sql.Stmt, error) {
	fmt.Println(query)
	return nil, nil
}

//go:noinline
func WrapDBprepareDC(db *sql.DB, ctx context.Context, dc *driverConn, release func(error), cg stmtConnGrabber, query string) (*sql.Stmt, error) {
	recursiveChecker := &recursiveCheck{rlsID: 2, success: false}

	begin, enter := recursiveChecker.enter()
	defer func() {
		recursiveChecker.leave()
		if exception := recover(); exception != nil {
			panic(exception)
		}
	}()
	r, e := DBprepareDC(db, ctx, dc, release, cg, query)
	if enter {
		coreWrapPrepareContext(begin, db, query, r, e)
	}
	return r, e
}

func readConfigBoolean(name string, defaultValue bool) bool {
	v, exist := tingyun3.ConfigRead(name)
	if !exist {
		return defaultValue
	}
	if value, ok := v.(bool); ok {
		return value
	}
	if value, ok := v.(string); ok {
		return tystring.CaseCMP(value, "true") == 0
	}
	return defaultValue
}
func matchMethod(method, matcher string) bool {
	return tystring.SubString(method, 0, len(matcher)) == matcher
}
func isNativeMethod(method string) bool {

	if matchMethod(method, "database/sql.") {
		return true
	}
	if matchMethod(method, "git.codemonky.net/TingYunGo/goagent") {
		return true
	}
	if readConfigBoolean("GORM_ENABLED", false) {
		return matchMethod(method, "gorm.io/")
	}
	return false
}

type recursiveCheck struct {
	rlsID   int
	success bool
}

func (r *recursiveCheck) enter() (time.Time, bool) {
	if r.success {
		return time.Time{}, false
	}
	if data := tingyun3.LocalGet(r.rlsID); data != nil {
		return time.Time{}, false
	}
	r.success = true
	tingyun3.LocalSet(r.rlsID, 1)
	return time.Now(), true
}
func (r *recursiveCheck) leave() {
	if r.success {
		tingyun3.LocalDelete(r.rlsID)
		r.success = false
	}
}

func showStack() {
	fmt.Println("Routine:", tingyun3.GetGID())
	i := 1
	for i > 0 {
		if _, pc := tingyun3.GetCallerPC(i + 1); pc != 0 {
			pcinfo := runtime.FuncForPC(pc)
			file, line := pcinfo.FileLine(pc)
			fmt.Println("  ", i, pcinfo.Name(), file, ":", line)
		} else {
			break
		}
		i++
	}
}
func init() {
	dbs.init()
	tingyun3.Register(reflect.ValueOf(WrapDBOpen).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapDBClose).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapDBPrepareContext).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapDBQueryContext).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapDBExecContext).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapTxPrepareContext).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapTxQueryContext).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapTxExecContext).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapConnPrepareContext).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapConnQueryContext).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapConnExecContext).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapStmtQueryContext).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapStmtExecContext).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapStmtClose).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapRowsClose).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapDBqueryDC).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapDBexecDC).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapDBprepareDC).Pointer())
}
