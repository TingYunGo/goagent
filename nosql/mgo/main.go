// +build linux
// +build amd64 arm64
// +build cgo

package mgo

import (
	"reflect"
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/TingYunGo/goagent"
	"github.com/TingYunGo/goagent/libs/tystring"
)

const (
	StorageIndexMgo = tingyun3.StorageIndexMgo
)

var skipTokens = []string{
	"gopkg.in/mgo%2ev2/",
	"github.com/TingYunGo/goagent",
}

func isNative(methodName string) bool {
	for _, t := range skipTokens {
		if matchMethod(methodName, t) {
			return true
		}
	}
	return false
}

func matchMethod(method, matcher string) bool {
	return tystring.SubString(method, 0, len(matcher)) == matcher
}

//go:noinline
func getCallName(skip int) (callerName string) {
	callerTmp := [8]uintptr{}
	callerName = tingyun3.FindCallerName(skip+1, callerTmp[:], isNative)
	return
}

type mgoHostInfo struct {
	hosts []string
}

type mgoSessionSet struct {
	lock  sync.RWMutex
	items map[*mgo.Session]mgoHostInfo
}

func (d *mgoSessionSet) init() {
	d.items = map[*mgo.Session]mgoHostInfo{}
}
func (d *mgoSessionSet) Get(session *mgo.Session) *mgoHostInfo {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if i, found := d.items[session]; found {
		return &i
	}
	return nil
}
func (d *mgoSessionSet) Set(session *mgo.Session, info mgoHostInfo) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.items == nil {
		d.init()
	}
	d.items[session] = info
}
func (d *mgoSessionSet) Delete(session *mgo.Session) {
	d.lock.Lock()
	defer d.lock.Unlock()
	delete(d.items, session)
}

var sessions mgoSessionSet

func getMgoHostName(session *mgo.Session) string {
	hostInfo := sessions.Get(session)
	if hostInfo != nil {
		if len(hostInfo.hosts) > 0 {
			return hostInfo.hosts[0]
		}
	}
	return "Unknown"
}

//---------------------------mgo Methods-------------------------------//

//go:noinline
func mgoDialWithInfo(info *mgo.DialInfo) (*mgo.Session, error) {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmgoDialWithInfo(info *mgo.DialInfo) (*mgo.Session, error) {
	session, err := mgoDialWithInfo(info)
	if session != nil {
		hostInfo := mgoHostInfo{[]string{}}
		for _, host := range info.Addrs {
			hostInfo.hosts = append(hostInfo.hosts, host)
		}
		sessions.Set(session, hostInfo)
	}
	return session, err
}

//go:noinline
func mgocopySession(session *mgo.Session, keepCreds bool) *mgo.Session {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgocopySession(session *mgo.Session, keepCreds bool) *mgo.Session {
	r := mgocopySession(session, keepCreds)
	if r != nil && r != session {

		if pHostInfo := sessions.Get(session); pHostInfo != nil {
			sessions.Set(r, *pHostInfo)
		}
	}
	return r
}

//---------------------------Session Methods-------------------------------//

//go:noinline
func mgoSessionClose(session *mgo.Session) {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
}

//go:noinline
func WrapmgoSessionClose(session *mgo.Session) {
	mgoSessionClose(session)
	sessions.Delete(session)
}

//go:noinline
func mgoSessionFindRef(s *mgo.Session, ref *mgo.DBRef) *mgo.Query {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoSessionFindRef(s *mgo.Session, ref *mgo.DBRef) *mgo.Query {

	componentInfo := beginOperation(s, ref.Database, ref.Collection, "FindRef", true)
	if componentInfo != nil {
		if componentInfo.handleType == 0 {
			componentInfo.handleType = 2
		}
	}
	q := mgoSessionFindRef(s, ref)
	return q
}

//---------------------------Database Methods-------------------------------//

//go:noinline
func mgoDatabaseFindRef(d *mgo.Database, ref *mgo.DBRef) *mgo.Query {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoDatabaseFindRef(db *mgo.Database, ref *mgo.DBRef) *mgo.Query {
	database := ref.Database
	if database == "" {
		database = db.Name
	}
	componentInfo := beginOperation(db.Session, database, ref.Collection, "FindRef", true)
	if componentInfo != nil {
		if componentInfo.handleType == 0 {
			componentInfo.handleType = 2
		}
	}
	q := mgoDatabaseFindRef(db, ref)
	return q
}

//---------------------------Collection Methods-------------------------------//

//go:noinline
func mgoCollectionCount(c *mgo.Collection) (int, error) {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return 0, nil
}

//go:noinline
func WrapmgoCollectionCount(c *mgo.Collection) (int, error) {
	componentInfo := enterOperation(c, "Count", false)
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType == 1 {
			handled = true
		} else {
			componentInfo.handleType = 1
		}
	}
	count, err := mgoCollectionCount(c)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return count, err
}

//go:noinline
func mgoCollectionCreate(c *mgo.Collection, info *mgo.CollectionInfo) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoCollectionCreate(c *mgo.Collection, info *mgo.CollectionInfo) error {
	componentInfo := enterOperation(c, "Create", false)
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType == 1 {
			handled = true
		} else {
			componentInfo.handleType = 1
		}
	}
	err := mgoCollectionCreate(c, info)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return err
}

//go:noinline
func mgoCollectionDropCollection(c *mgo.Collection) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoCollectionDropCollection(c *mgo.Collection) error {
	componentInfo := enterOperation(c, "DropCollection", false)
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType == 1 {
			handled = true
		} else {
			componentInfo.handleType = 1
		}
	}
	err := mgoCollectionDropCollection(c)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return err
}

//go:noinline
func mgoCollectionDropIndex(c *mgo.Collection, key ...string) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoCollectionDropIndex(c *mgo.Collection, key ...string) error {
	componentInfo := enterOperation(c, "DropIndex", false)
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType == 1 {
			handled = true
		} else {
			componentInfo.handleType = 1
		}
	}
	err := mgoCollectionDropIndex(c, key...)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return err
}

//go:noinline
func mgoCollectionDropIndexName(c *mgo.Collection, name string) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoCollectionDropIndexName(c *mgo.Collection, name string) error {
	componentInfo := enterOperation(c, "DropIndexName", false)
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType == 1 {
			handled = true
		} else {
			componentInfo.handleType = 1
		}
	}
	err := mgoCollectionDropIndexName(c, name)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return err
}

//go:noinline
func mgoCollectionEnsureIndex(c *mgo.Collection, index mgo.Index) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoCollectionEnsureIndex(c *mgo.Collection, index mgo.Index) error {
	componentInfo := enterOperation(c, "EnsureIndex", false)
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType == 1 {
			handled = true
		} else {
			componentInfo.handleType = 1
		}
	}
	err := mgoCollectionEnsureIndex(c, index)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return err
}

//go:noinline
func mgoCollectionEnsureIndexKey(c *mgo.Collection, key ...string) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoCollectionEnsureIndexKey(c *mgo.Collection, key ...string) error {
	componentInfo := enterOperation(c, "EnsureIndexKey", false)
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType == 1 {
			handled = true
		} else {
			componentInfo.handleType = 1
		}
	}
	err := mgoCollectionEnsureIndexKey(c, key...)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return err
}

//go:noinline
func mgoCollectionFind(c *mgo.Collection, query interface{}) *mgo.Query {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoCollectionFind(c *mgo.Collection, query interface{}) *mgo.Query {

	componentInfo := enterOperation(c, "Find", true)
	if componentInfo != nil {
		if componentInfo.handleType == 0 {
			componentInfo.handleType = 2
		}
	}
	q := mgoCollectionFind(c, query)
	return q
}

//go:noinline
func mgoCollectionFindId(c *mgo.Collection, query interface{}) *mgo.Query {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoCollectionFindId(c *mgo.Collection, query interface{}) *mgo.Query {

	componentInfo := enterOperation(c, "FindId", true)
	if componentInfo != nil {
		if componentInfo.handleType == 0 {
			componentInfo.handleType = 2
		}
	}
	q := mgoCollectionFindId(c, query)
	return q
}

//go:noinline
func mgoCollectionIndexes(c *mgo.Collection) ([]mgo.Index, error) {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmgoCollectionIndexes(c *mgo.Collection) ([]mgo.Index, error) {
	componentInfo := enterOperation(c, "Indexes", false)
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType == 1 {
			handled = true
		} else {
			componentInfo.handleType = 1
		}
	}
	indexes, err := mgoCollectionIndexes(c)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return indexes, err
}

//go:noinline
func mgoCollectionInsert(c *mgo.Collection, docs ...interface{}) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoCollectionInsert(c *mgo.Collection, docs ...interface{}) error {
	componentInfo := enterOperation(c, "Insert", false)
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType == 1 {
			handled = true
		} else {
			componentInfo.handleType = 1
		}
	}
	err := mgoCollectionInsert(c, docs...)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return err
}

//go:noinline
func mgoCollectionPipe(c *mgo.Collection, pipeline interface{}) *mgo.Pipe {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoCollectionPipe(c *mgo.Collection, pipeline interface{}) *mgo.Pipe {

	componentInfo := enterOperation(c, "Pipe", true)
	if componentInfo != nil {
		if componentInfo.handleType == 0 {
			componentInfo.handleType = 2
		}
	}
	q := mgoCollectionPipe(c, pipeline)
	return q
}

//go:noinline
func mgoCollectionRemove(c *mgo.Collection, selector interface{}) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoCollectionRemove(c *mgo.Collection, selector interface{}) error {
	componentInfo := enterOperation(c, "Remove", false)
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType == 1 {
			handled = true
		} else {
			componentInfo.handleType = 1
		}
	}
	err := mgoCollectionRemove(c, selector)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return err
}

//go:noinline
func mgoCollectionRemoveAll(c *mgo.Collection, selector interface{}) (*mgo.ChangeInfo, error) {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmgoCollectionRemoveAll(c *mgo.Collection, selector interface{}) (*mgo.ChangeInfo, error) {
	componentInfo := enterOperation(c, "RemoveAll", false)
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType == 1 {
			handled = true
		} else {
			componentInfo.handleType = 1
		}
	}
	info, err := mgoCollectionRemoveAll(c, selector)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return info, err
}

//go:noinline
func mgoCollectionRemoveId(c *mgo.Collection, id interface{}) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoCollectionRemoveId(c *mgo.Collection, id interface{}) error {
	componentInfo := enterOperation(c, "RemoveId", false)
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType == 1 {
			handled = true
		} else {
			componentInfo.handleType = 1
		}
	}
	err := mgoCollectionRemoveId(c, id)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return err
}

//go:noinline
func mgoCollectionRepair(c *mgo.Collection) *mgo.Iter {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoCollectionRepair(c *mgo.Collection) *mgo.Iter {
	componentInfo := enterOperation(c, "Repair", false)
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType == 1 {
			handled = true
		} else {
			componentInfo.handleType = 1
		}
	}
	iter := mgoCollectionRepair(c)

	if !handled {
		leaveOperation(componentInfo, nil)
	}
	return iter
}

//go:noinline
func mgoCollectionUpdate(c *mgo.Collection, selector interface{}, update interface{}) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoCollectionUpdate(c *mgo.Collection, selector interface{}, update interface{}) error {
	componentInfo := enterOperation(c, "Update", false)
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType == 1 {
			handled = true
		} else {
			componentInfo.handleType = 1
		}
	}
	err := mgoCollectionUpdate(c, selector, update)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return err
}

//go:noinline
func mgoCollectionUpdateAll(c *mgo.Collection, selector interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmgoCollectionUpdateAll(c *mgo.Collection, selector interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	componentInfo := enterOperation(c, "UpdateAll", false)
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType == 1 {
			handled = true
		} else {
			componentInfo.handleType = 1
		}
	}
	info, err := mgoCollectionUpdateAll(c, selector, update)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return info, err
}

//go:noinline
func mgoCollectionUpdateId(c *mgo.Collection, id interface{}, update interface{}) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoCollectionUpdateId(c *mgo.Collection, id interface{}, update interface{}) error {
	componentInfo := enterOperation(c, "UpdateId", false)
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType == 1 {
			handled = true
		} else {
			componentInfo.handleType = 1
		}
	}
	err := mgoCollectionUpdateId(c, id, update)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return err
}

//go:noinline
func mgoCollectionUpsert(c *mgo.Collection, selector interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmgoCollectionUpsert(c *mgo.Collection, selector interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	componentInfo := enterOperation(c, "Upsert", false)
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType == 1 {
			handled = true
		} else {
			componentInfo.handleType = 1
		}
	}
	info, err := mgoCollectionUpsert(c, selector, update)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return info, err
}

//go:noinline
func mgoCollectionUpsertId(c *mgo.Collection, id interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmgoCollectionUpsertId(c *mgo.Collection, id interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	componentInfo := enterOperation(c, "UpsertId", false)
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType == 1 {
			handled = true
		} else {
			componentInfo.handleType = 1
		}
	}
	info, err := mgoCollectionUpsertId(c, id, update)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return info, err
}

//---------------------------Query Methods-------------------------------//

//go:noinline
func mgoQueryAll(q *mgo.Query, result interface{}) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoQueryAll(q *mgo.Query, result interface{}) error {

	componentInfo := getComponentInfo()
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType != 1 {
			componentInfo.operation += ".All"
			componentInfo.handleType = 1
		} else {
			handled = true
		}
	}
	r := mgoQueryAll(q, result)

	if !handled {
		leaveOperation(componentInfo, r)
	}
	return r
}

//go:noinline
func mgoQueryApply(q *mgo.Query, change mgo.Change, result interface{}) (*mgo.ChangeInfo, error) {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmgoQueryApply(q *mgo.Query, change mgo.Change, result interface{}) (*mgo.ChangeInfo, error) {

	componentInfo := getComponentInfo()
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType != 1 {
			componentInfo.operation += ".Apply"
			componentInfo.handleType = 1
		} else {
			handled = true
		}
	}
	r, e := mgoQueryApply(q, change, result)

	if !handled {
		leaveOperation(componentInfo, e)
	}
	return r, e
}

//go:noinline
func mgoQueryCount(q *mgo.Query) (int, error) {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return 0, nil
}

//go:noinline
func WrapmgoQueryCount(q *mgo.Query) (int, error) {

	componentInfo := getComponentInfo()
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType != 1 {
			componentInfo.operation += ".Count"
			componentInfo.handleType = 1
		} else {
			handled = true
		}
	}
	r, err := mgoQueryCount(q)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return r, err
}

//go:noinline
func mgoQueryDistinct(q *mgo.Query, key string, result interface{}) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoQueryDistinct(q *mgo.Query, key string, result interface{}) error {

	componentInfo := getComponentInfo()
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType != 1 {
			componentInfo.operation += ".Distinct"
			componentInfo.handleType = 1
		} else {
			handled = true
		}
	}
	err := mgoQueryDistinct(q, key, result)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return err
}

//go:noinline
func mgoQueryExplain(q *mgo.Query, result interface{}) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoQueryExplain(q *mgo.Query, result interface{}) error {

	componentInfo := getComponentInfo()
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType != 1 {
			componentInfo.operation += ".Explain"
			componentInfo.handleType = 1
		} else {
			handled = true
		}
	}
	err := mgoQueryExplain(q, result)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return err
}

//go:noinline
func mgoQueryFor(q *mgo.Query, result interface{}, f func() error) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoQueryFor(q *mgo.Query, result interface{}, f func() error) error {

	componentInfo := getComponentInfo()
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType != 1 {
			componentInfo.operation += ".For"
			componentInfo.handleType = 1
		} else {
			handled = true
		}
	}
	err := mgoQueryFor(q, result, f)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return err
}

//go:noinline
func mgoQueryIter(q *mgo.Query) *mgo.Iter {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoQueryIter(q *mgo.Query) *mgo.Iter {

	componentInfo := getComponentInfo()
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType != 1 {
			componentInfo.operation += ".Iter"
			componentInfo.handleType = 1
		} else {
			handled = true
		}
	}
	r := mgoQueryIter(q)

	if !handled {
		leaveOperation(componentInfo, nil)
	}
	return r
}

//go:noinline
func mgoQueryMapReduce(q *mgo.Query, job *mgo.MapReduce, result interface{}) (*mgo.MapReduceInfo, error) {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmgoQueryMapReduce(q *mgo.Query, job *mgo.MapReduce, result interface{}) (*mgo.MapReduceInfo, error) {

	componentInfo := getComponentInfo()
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType != 1 {
			componentInfo.operation += ".MapReduce"
			componentInfo.handleType = 1
		} else {
			handled = true
		}
	}
	r, err := mgoQueryMapReduce(q, job, result)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return r, err
}

//go:noinline
func mgoQueryOne(q *mgo.Query, result interface{}) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoQueryOne(q *mgo.Query, result interface{}) error {

	componentInfo := getComponentInfo()
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType != 1 {
			componentInfo.operation += ".One"
			componentInfo.handleType = 1
		} else {
			handled = true
		}
	}
	err := mgoQueryOne(q, result)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return err
}

//go:noinline
func mgoQueryTail(q *mgo.Query, timeout time.Duration) *mgo.Iter {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoQueryTail(q *mgo.Query, timeout time.Duration) *mgo.Iter {

	componentInfo := getComponentInfo()
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType != 1 {
			componentInfo.operation += ".Tail"
			componentInfo.handleType = 1
		} else {
			handled = true
		}
	}
	r := mgoQueryTail(q, timeout)

	if !handled {
		leaveOperation(componentInfo, nil)
	}
	return r
}

//---------------------------Pipe Methods-------------------------------//

//go:noinline
func mgoPipeAll(q *mgo.Pipe, result interface{}) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoPipeAll(q *mgo.Pipe, result interface{}) error {

	componentInfo := getComponentInfo()
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType != 1 {
			componentInfo.operation += ".All"
			componentInfo.handleType = 1
		} else {
			handled = true
		}
	}
	r := mgoPipeAll(q, result)

	if !handled {
		leaveOperation(componentInfo, r)
	}
	return r
}

//go:noinline
func mgoPipeExplain(q *mgo.Pipe, result interface{}) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoPipeExplain(q *mgo.Pipe, result interface{}) error {

	componentInfo := getComponentInfo()
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType != 1 {
			componentInfo.operation += ".Explain"
			componentInfo.handleType = 1
		} else {
			handled = true
		}
	}
	err := mgoPipeExplain(q, result)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return err
}

//go:noinline
func mgoPipeIter(q *mgo.Pipe) *mgo.Iter {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoPipeIter(q *mgo.Pipe) *mgo.Iter {

	componentInfo := getComponentInfo()
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType != 1 {
			componentInfo.operation += ".Iter"
			componentInfo.handleType = 1
		} else {
			handled = true
		}
	}
	r := mgoPipeIter(q)

	if !handled {
		leaveOperation(componentInfo, nil)
	}
	return r
}

//go:noinline
func mgoPipeOne(q *mgo.Pipe, result interface{}) error {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmgoPipeOne(q *mgo.Pipe, result interface{}) error {

	componentInfo := getComponentInfo()
	handled := false
	if componentInfo != nil {
		if componentInfo.handleType != 1 {
			componentInfo.operation += ".One"
			componentInfo.handleType = 1
		} else {
			handled = true
		}
	}
	err := mgoPipeOne(q, result)

	if !handled {
		leaveOperation(componentInfo, err)
	}
	return err
}

func init() {
	sessions.init()
	tingyun3.Register(reflect.ValueOf(WrapmgoDialWithInfo).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgocopySession).Pointer())

	tingyun3.Register(reflect.ValueOf(WrapmgoSessionClose).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoSessionFindRef).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoDatabaseFindRef).Pointer())

	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionCount).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionCreate).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionDropCollection).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionDropIndex).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionDropIndexName).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionEnsureIndex).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionEnsureIndexKey).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionFind).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionFindId).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionIndexes).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionInsert).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionPipe).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionRemove).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionRemoveAll).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionRemoveId).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionRepair).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionUpdate).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionUpdateAll).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionUpdateId).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionUpsert).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoCollectionUpsertId).Pointer())

	tingyun3.Register(reflect.ValueOf(WrapmgoQueryAll).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoQueryApply).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoQueryCount).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoQueryDistinct).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoQueryExplain).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoQueryFor).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoQueryIter).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoQueryMapReduce).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoQueryOne).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoQueryTail).Pointer())

	tingyun3.Register(reflect.ValueOf(WrapmgoPipeAll).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoPipeExplain).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoPipeIter).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmgoPipeOne).Pointer())

	tingyun3.Register(reflect.ValueOf(initTrampoline).Pointer())
}

func leaveOperation(c *componentContext, e error) {
	if c == nil {
		return
	}
	action := tingyun3.GetAction()
	if action != nil {
		defer func() {
			if action.IsTask() {
				action.Finish()
				tingyun3.LocalClear()
			}
		}()
		component := action.CreateMongoComponent(c.host, c.database, c.collection, c.operation, c.callerName)
		if component != nil {
			component.FixBegin(c.beginTime)
			if e != nil {
				component.SetException(e, "mgo.Collection", 2)
			}
			component.End(2)
		}
	}
	tingyun3.LocalDelete(StorageIndexMgo)
}
func getComponentInfo() *componentContext {
	if prehandle := tingyun3.LocalGet(StorageIndexMgo); prehandle != nil {
		if info, ok := prehandle.(*componentContext); ok {
			return info
		}
	}
	return nil
}
func enterOperation(c *mgo.Collection, operationName string, isPreinfo bool) *componentContext {
	return beginOperation(c.Database.Session, c.Database.Name, c.Name, operationName, isPreinfo)
}

func beginOperation(session *mgo.Session, database, collection, operationName string, isPreinfo bool) *componentContext {
	if info := getComponentInfo(); info != nil {
		return info
	}
	action := tingyun3.GetAction()
	callerName := ""
	if action == nil && !isPreinfo {
		callerName = getCallName(2)
		if action, _ = tingyun3.CreateTask(callerName); action != nil {
			tingyun3.SetAction(action)
		}
		if action == nil {
			return nil
		}
	}
	if callerName == "" {
		callerName = getCallName(2)
	}
	componentInfo := &componentContext{
		host:       getMgoHostName(session),
		database:   database,
		collection: collection,
		operation:  operationName,
		callerName: callerName,
		beginTime:  time.Now(),
		handleType: 0,
		hasAction:  action != nil,
	}

	tingyun3.LocalSet(StorageIndexMgo, componentInfo)
	return componentInfo
}

type componentContext struct {
	host       string
	database   string
	collection string
	operation  string
	callerName string
	beginTime  time.Time
	handleType int //1:InHandler; 2: Prepair Info
	hasAction  bool
}
