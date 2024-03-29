// +build linux
// +build amd64 arm64
// +build cgo

package mongodb

import (
	"context"
	"reflect"
	"sync"
	"time"

	// "github.com/mongodb/mongo-go-driver/bson"
	// "github.com/mongodb/mongo-go-driver/mongo"
	// "github.com/mongodb/mongo-go-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/TingYunGo/goagent"
	"github.com/TingYunGo/goagent/libs/tystring"
)

type mongoInfo struct {
	hosts []string
}

const (
	StorageIndexMongo = tingyun3.StorageIndexMongo
)

var skipTokens = []string{
	"go.mongodb.org/mongo-driver/",
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

type mongoClientSet struct {
	lock  sync.RWMutex
	items map[*mongo.Client]mongoInfo
}

func (d *mongoClientSet) init() {
	d.items = map[*mongo.Client]mongoInfo{}
}
func (d *mongoClientSet) Get(client *mongo.Client) *mongoInfo {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if i, found := d.items[client]; found {
		return &i
	}
	return nil
}
func (d *mongoClientSet) Set(client *mongo.Client, info mongoInfo) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.items == nil {
		d.init()
	}
	d.items[client] = info
}
func (d *mongoClientSet) Delete(client *mongo.Client) {
	d.lock.Lock()
	defer d.lock.Unlock()
	delete(d.items, client)
}

var clients mongoClientSet

func getMongoHostName(coll *mongo.Collection) string {
	if clientInfo := clients.Get(coll.Database().Client()); clientInfo != nil {
		if len(clientInfo.hosts) > 0 {
			return clientInfo.hosts[0]
		}
	}
	return "Unknown"
}

//go:noinline
func mongodbNewClient(opts ...*options.ClientOptions) (*mongo.Client, error) {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmongodbNewClient(opts ...*options.ClientOptions) (*mongo.Client, error) {
	hosts := []string{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if len(opt.Hosts) > 0 {
			hosts = opt.Hosts
		}
	}
	clientInfo := mongoInfo{[]string{}}
	for _, host := range hosts {
		clientInfo.hosts = append(clientInfo.hosts, host)
	}
	c, e := mongodbNewClient(opts...)
	if c != nil {
		clients.Set(c, clientInfo)
	}
	return c, e
}

//go:noinline
func mongodbDisconnect(c *mongo.Client, ctx context.Context) error {
	trampoline.arg3 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmongodbDisconnect(c *mongo.Client, ctx context.Context) error {
	clients.Delete(c)
	return mongodbDisconnect(c, ctx)
}

func methodEnter(ctx context.Context) (interface{}, *tingyun3.Action, time.Time, bool) {
	if prehandle := tingyun3.LocalGet(StorageIndexMongo); prehandle != nil {
		return prehandle, nil, time.Time{}, false
	}
	action, _ := tingyun3.FindAction(ctx)
	isTask := false
	if action == nil {
		if action, _ = tingyun3.CreateTask(getCallName(2)); action != nil {
			isTask = true
		}
	}
	tingyun3.LocalSet(StorageIndexMongo, 1)
	return nil, action, time.Now(), isTask
}
func methodLeave(prehandle interface{}, e error, coll *mongo.Collection, action *tingyun3.Action, begin time.Time, invokeName string, isTask bool) {
	if prehandle != nil {
		return
	}
	defer func() {
		if action != nil && isTask {
			action.Finish()
		}
		if tingyun3.GetAction() == nil {
			tingyun3.LocalClear()
		}
	}()
	tingyun3.LocalDelete(StorageIndexMongo)
	if action == nil {
		return
	}
	callerName := getCallName(2)
	component := action.CreateMongoComponent(getMongoHostName(coll), coll.Database().Name(), coll.Name(), invokeName, callerName)
	if component == nil {
		return
	}
	component.FixBegin(begin)
	if e != nil {
		component.SetException(e, "mongo.Collection", 2)
	}
	component.End(2)
}

//go:noinline
func mongodbBulkWrite(coll *mongo.Collection, ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	trampoline.arg4 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmongodbBulkWrite(coll *mongo.Collection, ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	handle, action, begin, isTask := methodEnter(ctx)
	res, e := mongodbBulkWrite(coll, ctx, models, opts...)
	methodLeave(handle, e, coll, action, begin, "BulkWrite", isTask)
	return res, e
}

//go:noinline
func mongodbInsertOne(coll *mongo.Collection, ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	trampoline.arg5 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmongodbInsertOne(coll *mongo.Collection, ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	handle, action, begin, isTask := methodEnter(ctx)
	r, e := mongodbInsertOne(coll, ctx, document, opts...)
	methodLeave(handle, e, coll, action, begin, "InsertOne", isTask)
	return r, e
}

//go:noinline
func mongodbInsertMany(coll *mongo.Collection, ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	trampoline.arg6 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmongodbInsertMany(coll *mongo.Collection, ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	handle, action, begin, isTask := methodEnter(ctx)
	r, e := mongodbInsertMany(coll, ctx, documents, opts...)
	methodLeave(handle, e, coll, action, begin, "InsertMany", isTask)
	return r, e
}

//go:noinline
func mongodbDeleteOne(coll *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	trampoline.arg7 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmongodbDeleteOne(coll *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	handle, action, begin, isTask := methodEnter(ctx)
	r, e := mongodbDeleteOne(coll, ctx, filter, opts...)
	methodLeave(handle, e, coll, action, begin, "DeleteOne", isTask)
	return r, e
}

//go:noinline
func mongodbDeleteMany(coll *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	trampoline.arg8 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmongodbDeleteMany(coll *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	handle, action, begin, isTask := methodEnter(ctx)
	r, e := mongodbDeleteMany(coll, ctx, filter, opts...)
	methodLeave(handle, e, coll, action, begin, "DeleteMany", isTask)
	return r, e
}

//go:noinline
func mongodbUpdateByID(coll *mongo.Collection, ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	trampoline.arg9 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmongodbUpdateByID(coll *mongo.Collection, ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	handle, action, begin, isTask := methodEnter(ctx)
	r, e := mongodbUpdateByID(coll, ctx, id, update, opts...)
	methodLeave(handle, e, coll, action, begin, "UpdateByID", isTask)
	return r, e
}

//go:noinline
func mongodbUpdateOne(coll *mongo.Collection, ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	trampoline.arg10 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmongodbUpdateOne(coll *mongo.Collection, ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	handle, action, begin, isTask := methodEnter(ctx)
	r, e := mongodbUpdateOne(coll, ctx, filter, update, opts...)
	methodLeave(handle, e, coll, action, begin, "UpdateOne", isTask)
	return r, e
}

//go:noinline
func mongodbUpdateMany(coll *mongo.Collection, ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	trampoline.arg11 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmongodbUpdateMany(coll *mongo.Collection, ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	handle, action, begin, isTask := methodEnter(ctx)
	r, e := mongodbUpdateMany(coll, ctx, filter, update, opts...)
	methodLeave(handle, e, coll, action, begin, "UpdateMany", isTask)
	return r, e
}

//go:noinline
func mongodbReplaceOne(coll *mongo.Collection, ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	trampoline.arg12 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmongodbReplaceOne(coll *mongo.Collection, ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	handle, action, begin, isTask := methodEnter(ctx)
	r, e := mongodbReplaceOne(coll, ctx, filter, replacement, opts...)
	methodLeave(handle, e, coll, action, begin, "ReplaceOne", isTask)
	return r, e
}

//go:noinline
func mongodbAggregate(coll *mongo.Collection, ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	trampoline.arg13 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmongodbAggregate(coll *mongo.Collection, ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	handle, action, begin, isTask := methodEnter(ctx)
	r, e := mongodbAggregate(coll, ctx, pipeline, opts...)
	methodLeave(handle, e, coll, action, begin, "Aggregate", isTask)
	return r, e
}

//go:noinline
func mongodbCountDocuments(coll *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	trampoline.arg14 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return 0, nil
}

//go:noinline
func WrapmongodbCountDocuments(coll *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	handle, action, begin, isTask := methodEnter(ctx)
	r, e := mongodbCountDocuments(coll, ctx, filter, opts...)
	methodLeave(handle, e, coll, action, begin, "CountDocuments", isTask)
	return r, e
}

//go:noinline
func mongodbEstimatedDocumentCount(coll *mongo.Collection, ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	trampoline.arg15 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return 0, nil
}

//go:noinline
func WrapmongodbEstimatedDocumentCount(coll *mongo.Collection, ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	handle, action, begin, isTask := methodEnter(ctx)
	r, e := mongodbEstimatedDocumentCount(coll, ctx, opts...)
	methodLeave(handle, e, coll, action, begin, "EstimatedDocumentCount", isTask)
	return r, e
}

//go:noinline
func mongodbDistinct(coll *mongo.Collection, ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	trampoline.arg16 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmongodbDistinct(coll *mongo.Collection, ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	handle, action, begin, isTask := methodEnter(ctx)
	r, e := mongodbDistinct(coll, ctx, fieldName, filter, opts...)
	methodLeave(handle, e, coll, action, begin, "Distinct", isTask)
	return r, e
}

//go:noinline
func mongodbFind(coll *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	trampoline.arg17 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmongodbFind(coll *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	handle, action, begin, isTask := methodEnter(ctx)
	r, e := mongodbFind(coll, ctx, filter, opts...)
	methodLeave(handle, e, coll, action, begin, "Find", isTask)
	return r, e
}

//go:noinline
func mongodbFindOne(coll *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	trampoline.arg18 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmongodbFindOne(coll *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	handle, action, begin, isTask := methodEnter(ctx)
	r := mongodbFindOne(coll, ctx, filter, opts...)
	methodLeave(handle, r.Err(), coll, action, begin, "FindOne", isTask)
	return r
}

//go:noinline
func mongodbFindOneAndDelete(coll *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	trampoline.arg19 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmongodbFindOneAndDelete(coll *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	handle, action, begin, isTask := methodEnter(ctx)
	r := mongodbFindOneAndDelete(coll, ctx, filter, opts...)
	methodLeave(handle, r.Err(), coll, action, begin, "FindOneAndDelete", isTask)
	return r
}

//go:noinline
func mongodbFindOneAndReplace(coll *mongo.Collection, ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	trampoline.arg20 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmongodbFindOneAndReplace(coll *mongo.Collection, ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	handle, action, begin, isTask := methodEnter(ctx)
	r := mongodbFindOneAndReplace(coll, ctx, filter, replacement, opts...)
	methodLeave(handle, r.Err(), coll, action, begin, "FindOneAndReplace", isTask)
	return r
}

//go:noinline
func mongodbFindOneAndUpdate(coll *mongo.Collection, ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	trampoline.arg1 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmongodbFindOneAndUpdate(coll *mongo.Collection, ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	handle, action, begin, isTask := methodEnter(ctx)
	r := mongodbFindOneAndUpdate(coll, ctx, filter, update, opts...)
	methodLeave(handle, r.Err(), coll, action, begin, "FindOneAndUpdate", isTask)
	return r
}

//go:noinline
func mongodbWatch(coll *mongo.Collection, ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapmongodbWatch(coll *mongo.Collection, ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	handle, action, begin, isTask := methodEnter(ctx)
	r, e := mongodbWatch(coll, ctx, pipeline, opts...)
	methodLeave(handle, e, coll, action, begin, "Watch", isTask)
	return r, e
}

//go:noinline
func mongodbDrop(coll *mongo.Collection, ctx context.Context) error {
	trampoline.arg3 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapmongodbDrop(coll *mongo.Collection, ctx context.Context) error {
	handle, action, begin, isTask := methodEnter(ctx)
	e := mongodbDrop(coll, ctx)
	methodLeave(handle, e, coll, action, begin, "Drop", isTask)
	return e
}

func init() {
	clients.init()
	tingyun3.Register(reflect.ValueOf(WrapmongodbNewClient).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbDisconnect).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbBulkWrite).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbInsertOne).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbInsertMany).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbDeleteOne).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbDeleteMany).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbUpdateMany).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbUpdateByID).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbUpdateOne).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbReplaceOne).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbCountDocuments).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbEstimatedDocumentCount).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbDistinct).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbFind).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbFindOne).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbFindOneAndDelete).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbFindOneAndReplace).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbFindOneAndUpdate).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbWatch).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapmongodbDrop).Pointer())
	tingyun3.Register(reflect.ValueOf(initTrampoline).Pointer())
}
