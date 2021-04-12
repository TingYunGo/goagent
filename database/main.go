// +build linux
// +build amd64

package database

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
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
	if tystring.CaseCMP(tystring.SubString(name, 0, 5), "mysql") == 0 {
		return tingyun3.ComponentMysql
	} else if tystring.CaseCMP(tystring.SubString(name, 0, 7), "postgre") == 0 {
		return tingyun3.ComponentPostgreSQL
	} else if tystring.CaseCMP(tystring.SubString(name, 0, 9), "sqlserver") == 0 {
		return tingyun3.ComponentMSSQL
	} else if tystring.CaseCMP(tystring.SubString(name, 0, 6), "sqlite") == 0 {
		return tingyun3.ComponentSQLite
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
	// action := tingyun3.GetAction()
	// fmt.Println("action:", action.GetName())
	// fmt.Println("sql.Open(", driverName, ", ", dataSourceName, ")")
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
	return DBClose(db)
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
	// fmt.Println("DSN:", info.dsn, " , host:", info.host, ", dbname: ", info.dbname)
	component := action.CreateSQLComponent(getTingyunDBType(info.vender), info.host, info.dbname, query, callerName)
	component.FixBegin(begin)
	//PrepareContext失败,
	if stmt == nil || e != nil {
		// fmt.Println("prepare error: ", e)
		component.SetError(e, callerName, 3)
		component.End(1)
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

	// fmt.Println(callerName, "create db component ", info.vender, ":", info.dsn, ":", query)
	// fmt.Println("Host:", host, ",db:", dbname)
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
		component.SetError(e, callerName, 3)
		component.End(1)
		return
	}
	component.End(1)

	// fmt.Println(callerName, "create db component ", info.vender, ":", info.dsn, ":", query)
	// fmt.Println("Host:", host, ",db:", dbname)
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
	//QueryContext失败,
	if r == nil && e != nil {
		component.SetError(e, callerName, 3)
		component.End(1)
		return
	}
	if dbctx == nil {
		dbctx = (&databaseContext{}).init()
		tingyun3.LocalSet(1, dbctx)
	}
	dbctx.records[r] = component

	// fmt.Println(callerName, "create db component ", info.vender, ":", info.dsn, ":", query)
	// fmt.Println("Host:", host, ",db:", dbname)
}

//go:noinline
func DBPrepareContext(db *sql.DB, ctx context.Context, query string) (*sql.Stmt, error) {
	fmt.Println(query)
	return nil, nil
}

//go:noinline
func WrapDBPrepareContext(db *sql.DB, ctx context.Context, query string) (*sql.Stmt, error) {
	begin := time.Now()
	// fmt.Println("On DB prepare: ", query)
	stmt, e := DBPrepareContext(db, ctx, query)
	coreWrapPrepareContext(begin, db, query, stmt, e)
	return stmt, e
}

//go:noinline
func DBExecContext(db *sql.DB, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	fmt.Println(query)
	return nil, nil
}

//go:noinline
func WrapDBExecContext(db *sql.DB, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {

	begin := time.Now()
	r, e := DBExecContext(db, ctx, query, args...)
	coreWrapExecContext(begin, db, query, r, e)
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
	for callerName[0:13] == "database/sql." {
		skip++
		callerName = tingyun3.GetCallerName(skip)
	}
	return
}

//go:noinline
func WrapDBQueryContext(db *sql.DB, ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	begin := time.Now()
	r, e := DBQueryContext(db, ctx, query, args...)
	coreWrapQueryContext(begin, db, query, r, e)
	return r, e
}

//go:noinline
func ConnExecContext(c *sql.Conn, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	fmt.Println(query)
	return nil, nil
}

//go:noinline
func WrapConnExecContext(c *sql.Conn, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	begin := time.Now()
	r, e := ConnExecContext(c, ctx, query, args...)
	coreWrapExecContext(begin, getdb_byconn(c), query, r, e)
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
	begin := time.Now()
	r, e := ConnQueryContext(c, ctx, query, args...)
	coreWrapQueryContext(begin, getdb_byconn(c), query, r, e)
	return r, e
}

//go:noinline
func ConnPrepareContext(c *sql.Conn, ctx context.Context, query string) (*sql.Stmt, error) {
	fmt.Println(query)
	return nil, nil
}

//go:noinline
func WrapConnPrepareContext(c *sql.Conn, ctx context.Context, query string) (*sql.Stmt, error) {
	begin := time.Now()
	stmt, e := ConnPrepareContext(c, ctx, query)
	coreWrapPrepareContext(begin, getdb_byconn(c), query, stmt, e)
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
	begin := time.Now()
	stmt, e := TxPrepareContext(tx, ctx, query)
	coreWrapPrepareContext(begin, getdbByTx(tx), query, stmt, e)
	return stmt, e
}

//go:noinline
func TxExecContext(tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	fmt.Println(query)
	return nil, nil
}

//go:noinline
func WrapTxExecContext(tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	begin := time.Now()
	r, e := TxExecContext(tx, ctx, query, args...)
	coreWrapExecContext(begin, getdbByTx(tx), query, r, e)
	return r, e
}

//go:noinline
func TxQueryContext(tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	fmt.Println(query)
	return nil, nil
}

//go:noinline
func WrapTxQueryContext(tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	begin := time.Now()
	r, e := TxQueryContext(tx, ctx, query, args...)
	coreWrapQueryContext(begin, getdbByTx(tx), query, r, e)
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
	// fmt.Println("Stmt QueryContext")
	if c := tingyun3.LocalGet(1); c != nil {
		dbctx := c.(*databaseContext)
		component, found := dbctx.stmts[s]
		if !found {
			return r, e
		}
		delete(dbctx.stmts, s)
		if r == nil && e != nil {
			callerName := getCallName(2)
			component.SetError(e, callerName, 2)
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

	// fmt.Println("Stmt Close")
	if c := tingyun3.LocalGet(1); c != nil {
		dbctx := c.(*databaseContext)
		if c, found := dbctx.stmts[s]; found {
			c.End(1)
			delete(dbctx.stmts, s)
			// fmt.Println("Component finished")
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

	// fmt.Println("Stmt Exec")
	if c := tingyun3.LocalGet(1); c != nil {
		dbctx := c.(*databaseContext)
		if c, found := dbctx.stmts[s]; found {
			c.End(1)
			// fmt.Println("Component finished on exec")
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
	// fmt.Println("Rows Close")
	if c := tingyun3.LocalGet(1); c != nil {
		dbctx := c.(*databaseContext)
		if c, found := dbctx.records[rs]; found {
			c.End(1)
			delete(dbctx.records, rs)
			// fmt.Println("Component finished")
		}
		if dbctx.empty() {
			tingyun3.LocalDelete(1)
		}
	}
	return err
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
}
